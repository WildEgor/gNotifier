# ------------ STAGE: Base
FROM golang:1.20 as base
ARG APP_PATH=/app
COPY . $APP_PATH
WORKDIR $APP_PATH
RUN go mod download && mkdir -p dist

# ------------ STAGE: Developing
FROM base as dev
WORKDIR /app
RUN go install -mod=mod github.com/cosmtrek/air
RUN cp /app/.air-unix.toml /app/.air.toml
ENTRYPOINT ["air"]

# ------------ STAGE: Test app
# FROM base as test
# ENTRYPOINT make test

# ------------ STAGE: Build app
FROM base as builder
ARG APP_PATH=/app
WORKDIR $APP_PATH
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dist/app cmd/main.go

# ------------ STAGE: Execute app
FROM alpine:latest as production
RUN apk --no-cache add ca-certificates
ARG APP_PATH=/root/
WORKDIR $APP_PATH
COPY --from=builder /app/dist/app .
EXPOSE 8888
RUN ls  
CMD ["./app"]