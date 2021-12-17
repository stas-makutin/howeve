package services

// log operation codes
const (
	SvcAddFromConfig = "C"
)

// log operation status codes
const (
	SvcOcSuccess                  = "0"
	SvcOcCfgUnknownProtocol       = "P"
	SvcOcCfgUnknownTransport      = "T"
	SvcOcCfgProtocolNotSupported  = "X"
	SvcOcCfgTransportNotSupported = "x"
	SvcOcCfgUnknownParameter      = "N"
	SvcOcCfgNoRequiredParameter   = "R"
	SvcOcCfgInvalidParameterValue = "V"
	SvcOcCfgAliasExists           = "A"
	SvcOcCfgCreateError           = "C"
)
