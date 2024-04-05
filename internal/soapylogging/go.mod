module soapy_logging

go 1.22.1

require (
	github.com/pothosware/go-soapy-sdr v0.7.4
	internal/settings v0.0.0-00010101000000-000000000000
)

replace internal/settings => ../settings
