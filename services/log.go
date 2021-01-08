package services

// log operation codes
const (
	SvcOpStart  = "S"
	SvcOpFinish = "F"
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
)
