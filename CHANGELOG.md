# Changelog

## 0.9.0: Initial Release (November 27, 2020)

In order to use this library you will need two things:

1. An [audit logger from the auditlog library](https://github.com/containerssh/auditlog).
2. A [handler from the sshserver library](https://github.com/containerssh/sshserver) to act as a backend to this library.

You can then create the audit logging handler like this:

```go
handler := auditlogintegration.New(
    backend,
    auditLogger,
)
```

You can then pass this handler to the SSH server as [described in the readme](https://github.com/containerssh/sshserver):

```go
server, err := sshserver.New(
    cfg,
    handler,
    logger,
)
```