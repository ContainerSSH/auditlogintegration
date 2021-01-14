# Changelog

## 0.9.3: Bumping version

This release bumps the version to work aroung Go caching.

## 0.9.2: Updated to sshserver 0.9.16

This release adds support for the new interactions in [sshserver](https://github.com/containerssh/sshserver) 0.9.16. 

## 0.9.1: Tests, bump to latest sshserver and auditlog (December 11, 2020)

This release updates the integration to match the latest [sshserver](https://github.com/containerssh/sshserver), [auditlog](https://github.com/containerssh/auditlog), and [log](https://github.com/containerssh/log).

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