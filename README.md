# gNotifier
Go notifier microservice (sms, smtp or push notifications). 
Expose API via amqp and REST.

## Features:
- Send email;
- Send sms;
- Send push (FCM or APN);
- Sub/unsub push tokens to own unique ID.

## Developing:
Wire DI container:
```bash
go install github.com/google/wire/cmd/wire@latest
cd internal
wire
```

Hot reload mode:
```bash
go install github.com/cosmtrek/air@latest
air -v
air server --port 8080
air -d
```

## Usage

```bash
   docker-compose up --build notifier
```

## Docs:

- [Fiber](https://gofiber.io/)
- [Wire](https://github.com/google/wire)
- [Testify](https://github.com/stretchr/testify)
