module github.com/containerssh/auditlogintegration

go 1.14

require (
	github.com/containerssh/auditlog v0.9.7
	github.com/containerssh/geoip v0.9.4
	github.com/containerssh/log v0.9.9
	github.com/containerssh/sshserver v0.9.16
	github.com/stretchr/testify v1.6.1
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
