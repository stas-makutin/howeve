package zwave

// Serial API command ID
const (
	ZW_VERSION = 0x15
)

// library type
const (
	ZW_LIB_CONTROLLER_STATIC = 0x01
	ZW_LIB_CONTROLLER        = 0x02
	ZW_LIB_SLAVE_ENHANCED    = 0x03
	ZW_LIB_SLAVE             = 0x04
	ZW_LIB_INSTALLER         = 0x05
	ZW_LIB_SLAVE_ROUTING     = 0x06
	ZW_LIB_CONTROLLER_BRIDGE = 0x07
	ZW_LIB_DUT               = 0x08
)
