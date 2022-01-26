package cablemodemutil

// CableModemRawStatus contains the raw status retrieved from the cable modem without any parsing.
type CableModemRawStatus map[string]interface{}

// CableModemDeviceInfo contains Cable Modem Device information.
type CableModemDeviceInfo struct {
	// Cable Modem model.
	Model string
	// Cable Modem serial number.
	SerialNumber string
	// Cable Modem MAC address.
	MACAddress string
}

// CableModemStatus contains Cable Modem Status.
type CableModemStatus struct {
	// Device related information.
	DeviceInfo CableModemDeviceInfo
}
