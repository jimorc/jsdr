module github.com/jimorc/jsdr/internal/settings

go 1.22.1

require (
    internal/soapylogging v1.0.0
)

require github.com/pothosware/go-soapy-sdr v0.7.4

replace internal/soapylogging => ../soapylogging
