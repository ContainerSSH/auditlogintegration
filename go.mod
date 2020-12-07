module github.com/containerssh/auditlogintegration

go 1.14

require (
	github.com/containerssh/auditlog v0.9.3
	github.com/containerssh/configuration v0.0.0-20201117205727-9e4ab1927a99
	github.com/containerssh/geoip v0.9.3
	github.com/containerssh/log v0.9.2
	github.com/containerssh/sshserver v0.9.8
)

replace github.com/containerssh/configuration v0.0.0-20201117205727-9e4ab1927a99 => ../configuration
