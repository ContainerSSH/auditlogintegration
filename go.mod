module github.com/containerssh/auditlogintegration

go 1.14

require (
	github.com/aws/aws-sdk-go v1.37.30 // indirect
	github.com/containerssh/auditlog v0.9.9
	github.com/containerssh/geoip v0.9.4
	github.com/containerssh/log v0.9.13
	github.com/containerssh/service v0.9.3 // indirect
	github.com/containerssh/sshserver v0.9.19
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/mattn/go-shellwords v1.0.11 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20210309074719-68d13333faf2 // indirect
)

replace (
	// Fixes CVE-2020-9283
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 => golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975 => golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 => golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9 => golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975
	// Fixes CVE-2020-14040
	golang.org/x/text v0.3.0 => golang.org/x/text v0.3.3
	golang.org/x/text v0.3.1 => golang.org/x/text v0.3.3
	golang.org/x/text v0.3.2 => golang.org/x/text v0.3.3
)
