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

// CableModemDeviceSettings contains Cable Modem Device settings.
type CableModemDeviceSettings struct {
	// True if front panel LED lights are configured to be on, false otherwise.
	FrontPanelLightsOn bool
	// True if energy efficient ethernet setting is turned on, false otherwise.
	EnergyEfficientEthernetOn bool
	// True if ask me later setting has been opted into, false otherwise.
	AskMeLater bool
	// True if never ask setting has been opted into, false otherwise.
	NeverAsk bool
}

// CableModemAuthSettings contains Cable Modem Authentication settings.
type CableModemAuthSettings struct {
	// Hash of the current login.
	CurrentLogin string
	// Hash of the admin username.
	CurrentNameAdmin string
	// Hash of the current user's username.
	CurrentNameUser string
	// Hash of the admin password.
	CurrentPasswordAdmin string
	// Hash of the current user's password.
	CurrentPasswordUser string
}

// CableModemSoftwareStatus contains Cable Modem Software status.
type CableModemSoftwareStatus struct {
	// Firmware version.
	FirmwareVersion string
	// True if certificate has been installed, false otherwise.
	CertificateInstalled bool
	// Customer version.
	CustomerVersion string
	// HD version.
	HDVersion string
	// DOCSIS specification version.
	DOCSISSpecVersion string
}

// CableModemStartupBootStatus contains Cable Modem Startup Boot status.
type CableModemStartupBootStatus struct {
	// Boot status.
	Status string
	// Comments.
	Comment string
}

// CableModemStartupConfigurationFileStatus contains Cable Modem Startup Configuration file status.
type CableModemStartupConfigurationFileStatus struct {
	// Configuration file status.
	Status string
	// Comments.
	Comment string
}

// CableModemStartupConnectivityStatus contains Cable Modem Startup Connectivity status.
type CableModemStartupConnectivityStatus struct {
	// Connectivity status.
	Status string
	// Comments.
	Comment string
}

// CableModemStartupDownstreamConnectionStatus contains Cable Modem Startup Downstream Connection status.
type CableModemStartupDownstreamConnectionStatus struct {
	// Frequency in Hz for the Downstream channel connection.
	FrequencyHZ uint32
	// Comments.
	Comment string
}

// CableModemStartupSecurityStatus contains Cable Modem Startup Security status.
type CableModemStartupSecurityStatus struct {
	// Security status.
	Status string
	// Comments.
	Comment string
}

// CableModemStartupStatus contains Cable Modem Startup Status.
type CableModemStartupStatus struct {
	// Boot status.
	Boot CableModemStartupBootStatus
	// Configuration file status.
	ConfigurationFile CableModemStartupConfigurationFileStatus
	// Connectivity status.
	Connectivity CableModemStartupConnectivityStatus
	// Downstream connection status.
	DownstreamConnection CableModemStartupDownstreamConnectionStatus
	// Security status.
	Security CableModemStartupSecurityStatus
}

// CableModemStatus contains Cable Modem Status.
type CableModemStatus struct {
	// Device related information.
	DeviceInfo CableModemDeviceInfo
	// General settings.
	DeviceSettings CableModemDeviceSettings
	// Auth settings.
	AuthSettings CableModemAuthSettings
	// Software status.
	SoftwareStatus CableModemSoftwareStatus
	// Startup status.
	StartupStatus CableModemStartupStatus
}
