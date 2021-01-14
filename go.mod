module github.com/containerssh/auditlogintegration

go 1.14

require (
	github.com/containerssh/auditlog v0.9.7
	github.com/containerssh/geoip v0.9.4
	github.com/containerssh/log v0.9.9
	github.com/containerssh/sshserver v0.9.16
	github.com/stretchr/testify v1.6.1
	golang.org/x/sys v0.0.0-20210113181707-4bcb84eeeb78 // indirect
	golang.org/x/text v0.3.4 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace (
	// Fixes CVE-2020-9283
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	// Fixes CVE-2020-14040
	golang.org/x/text v0.3.0 => golang.org/x/text v0.3.3
	golang.org/x/text v0.3.1 => golang.org/x/text v0.3.3
	golang.org/x/text v0.3.2 => golang.org/x/text v0.3.3
)
