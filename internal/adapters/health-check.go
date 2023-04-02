package adapters

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Status string

// Option is the health-container options type
type Option func(*HealthCheckAdapter) error

// Status aliases
const (
	StatusOK                 Status = "OK"
	StatusPartiallyAvailable Status = "Partially Available"
	StatusUnavailable        Status = "Unavailable"
	StatusTimeout            Status = "Timeout during health check"
)

type (
	// HealthCheckFunc is the func which executes the check.
	HealthCheckFunc func(ctx context.Context) error
	// HealthConfig carries the parameters to run the check.
	HealthConfig struct {
		// Name is the name of the resource to be checked.
		Name string
		// Timeout is the timeout defined for every check.
		Timeout time.Duration
		// SkipOnErr if set to true, it will retrieve StatusOK providing the error message from the failed resource.
		SkipOnErr bool
		// Check is the func which executes the check.
		Check HealthCheckFunc
	}
	// HealthCheckInfo represents the health check response.
	HealthCheckInfo struct {
		// Status is the check status.
		Status Status `json:"status"`
		// Timestamp is the time in which the check occurred.
		Timestamp time.Time `json:"timestamp"`
		// Failures holds the failed checks along with their messages.
		Failures map[string]string `json:"failures,omitempty"`
		// System holds information of the go process.
		*SystemInfo `json:"system,omitempty"`
		// Component holds information on the component for which checks are made
		ComponentInfo `json:"component"`
	}
	// SystemInfo runtime variables about the go process.
	SystemInfo struct {
		// Version is the go version.
		Version string `json:"version"`
		// GoroutinesCount is the number of the current goroutines.
		GoroutinesCount int `json:"goroutines_count"`
		// TotalAllocBytes is the total bytes allocated.
		TotalAllocBytes int `json:"total_alloc_bytes"`
		// HeapObjectsCount is the number of objects in the go heap.
		HeapObjectsCount int `json:"heap_objects_count"`
		// TotalAllocBytes is the bytes allocated and not yet freed.
		AllocBytes int `json:"alloc_bytes"`
	}
	// ComponentInfo descriptive values about the component for which checks are made
	ComponentInfo struct {
		// Name is the name of the component.
		Name string `json:"name"`
		// Version is the component version.
		Version string `json:"version"`
	}
	// HealthCheckAdapter is the health-checks container
	HealthCheckAdapter struct {
		mu                  sync.Mutex
		checks              map[string]HealthConfig
		maxConcurrent       int
		instrumentationName string
		component           ComponentInfo
		systemInfoEnabled   bool
	}
)

// New instantiates and build new health check container
func NewHealthCheckAdapter() (*HealthCheckAdapter, error) {
	h := &HealthCheckAdapter{
		checks:        make(map[string]HealthConfig),
		maxConcurrent: runtime.NumCPU(),
	}

	return h, nil
}

// Register registers a check config to be performed.
func (h *HealthCheckAdapter) Register(c HealthConfig) error {
	if c.Timeout == 0 {
		c.Timeout = time.Second * 2
	}

	if c.Name == "" {
		return errors.New("[HealthCheckAdapter] Health check must have a name to be registered")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.checks[c.Name]; ok {
		return fmt.Errorf("[HealthCheckAdapter] Health check %q is already registered", c.Name)
	}

	h.checks[c.Name] = c

	return nil
}

// ---- Public methods

// Measure runs all the registered health checks and returns summary status
func (h *HealthCheckAdapter) Measure(ctx context.Context) HealthCheckInfo {
	h.mu.Lock()
	defer h.mu.Unlock()

	status := StatusOK
	failures := make(map[string]string)

	limiterCh := make(chan bool, h.maxConcurrent)
	defer close(limiterCh)

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)
	for _, c := range h.checks {
		limiterCh <- true
		wg.Add(1)

		go func(c HealthConfig) {
			defer func() {
				<-limiterCh
				wg.Done()
			}()

			resCh := make(chan error)

			go func() {
				resCh <- c.Check(ctx)
				defer close(resCh)
			}()

			timeout := time.NewTimer(c.Timeout)

			select {
			case <-timeout.C:
				mu.Lock()
				defer mu.Unlock()

				failures[c.Name] = string(StatusTimeout)
				status = getAvailability(status, c.SkipOnErr)
			case res := <-resCh:
				if !timeout.Stop() {
					<-timeout.C
				}

				mu.Lock()
				defer mu.Unlock()

				if res != nil {
					failures[c.Name] = res.Error()
					status = getAvailability(status, c.SkipOnErr)
				}
			}
		}(c)
	}

	wg.Wait()

	var systemMetrics *SystemInfo
	if h.systemInfoEnabled {
		systemMetrics = newSystemMetrics()
	}

	return newCheck(h.component, status, systemMetrics, failures)
}

// ---- Private methods

func newCheck(c ComponentInfo, s Status, system *SystemInfo, failures map[string]string) HealthCheckInfo {
	return HealthCheckInfo{
		Status:        s,
		Timestamp:     time.Now(),
		Failures:      failures,
		SystemInfo:    system,
		ComponentInfo: c,
	}
}

func newSystemMetrics() *SystemInfo {
	s := runtime.MemStats{}
	runtime.ReadMemStats(&s)

	return &SystemInfo{
		Version:          runtime.Version(),
		GoroutinesCount:  runtime.NumGoroutine(),
		TotalAllocBytes:  int(s.TotalAlloc),
		HeapObjectsCount: int(s.HeapObjects),
		AllocBytes:       int(s.Alloc),
	}
}

func getAvailability(s Status, skipOnErr bool) Status {
	if skipOnErr && s != StatusUnavailable {
		return StatusPartiallyAvailable
	}

	return StatusUnavailable
}

// ---- Static methods

// SetChecks adds checks to newly instantiated health-container
func SetChecks(checks ...HealthConfig) Option {
	return func(h *HealthCheckAdapter) error {
		for _, c := range checks {
			if err := h.Register(c); err != nil {
				return fmt.Errorf("could not register check %q: %w", c.Name, err)
			}
		}
		return nil
	}
}

// SetComponent sets the component description of the component to which this check refer
func SetComponent(component ComponentInfo) Option {
	return func(h *HealthCheckAdapter) error {
		h.component = component
		return nil
	}
}

// SetMaxConcurrent sets max number of concurrently running checks.
// Set to 1 if want to run all checks sequentially.
func SetMaxConcurrent(n int) Option {
	return func(h *HealthCheckAdapter) error {
		h.maxConcurrent = n
		return nil
	}
}

// SetSystemInfo enables the option to return system information about the go process.
func SetSystemInfo() Option {
	return func(h *HealthCheckAdapter) error {
		h.systemInfoEnabled = true
		return nil
	}
}
