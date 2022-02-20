// Package cablemodemutil exposes functions to connect and retrieve the status from a cable modem.
package cablemodemutil

import "time"

// CableModemRawStatus contains the raw status retrieved from the cable modem without any parsing.
type CableModemRawStatus map[string]interface{}

// DeviceInfo contains Cable Modem Device information.
type DeviceInfo struct {
	// Cable Modem model.
	Model string
	// Cable Modem serial number.
	SerialNumber string
	// Cable Modem MAC address.
	MACAddress string
}

// DeviceSettings contains Cable Modem Device settings.
type DeviceSettings struct {
	// True if front panel LED lights are configured to be on, false otherwise.
	FrontPanelLightsOn bool
	// True if energy efficient ethernet setting is turned on, false otherwise.
	EnergyEfficientEthernetOn bool
	// True if ask me later setting has been opted into, false otherwise.
	AskMeLater bool
	// True if never ask setting has been opted into, false otherwise.
	NeverAsk bool
}

// AuthSettings contains Cable Modem Authentication settings.
type AuthSettings struct {
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

// SoftwareStatus contains Cable Modem Software status.
type SoftwareStatus struct {
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

// BootStatus contains Cable Modem Startup Boot status.
type BootStatus struct {
	// Boot status.
	Status bool
	// Operational.
	Operational bool
}

// ConfigFileStatus contains Cable Modem Startup Configuration file status.
type ConfigFileStatus struct {
	// Configuration file status.
	Status bool
	// Comments.
	Comment string
}

// ConnectivityStatus contains Cable Modem Startup Connectivity status.
type ConnectivityStatus struct {
	// Connectivity status.
	Status bool
	// Operational.
	Operational bool
}

// DownstreamStatus contains Cable Modem Startup Downstream Connection status.
type DownstreamStatus struct {
	// Frequency in Hz for the Downstream channel connection.
	FrequencyHZ uint32
	// Locked.
	Locked bool
}

// SecurityStatus contains Cable Modem Startup Security status.
type SecurityStatus struct {
	// Security status.
	Enabled bool
	// Comments.
	Comment string
}

// StartupStatus contains Cable Modem Startup Status.
type StartupStatus struct {
	// Boot status.
	Boot BootStatus
	// Configuration file status.
	ConfigFile ConfigFileStatus
	// Connectivity status.
	Connectivity ConnectivityStatus
	// Downstream connection status.
	Downstream DownstreamStatus
	// Security status.
	Security SecurityStatus
}

// DownstreamChannelInfo contains Cable Modem Downstream channel information.
type DownstreamChannelInfo struct {
	// Lock status.
	Locked bool
	// Modulation.
	Modulation string
	// Channel ID.
	ChannelID uint32
	// Frequency of the channel in Hz.
	FrequencyHZ uint32
	// Signal Power in dB mV.
	SignalPowerDBMV float32
	// Signal SNR/MER in dB.
	SignalSNRMERDB float32
	// Corrected errors.
	CorrectedErrors uint32
	// Uncorrected errors.
	UncorrectedErrors uint32
}

// UpstreamChannelInfo contains Cable Modem Upstream channel information.
type UpstreamChannelInfo struct {
	// Lock status.
	Locked bool
	// Modulation.
	Modulation string
	// Channel ID.
	ChannelID uint32
	// Width of the channel in Hz.
	WidthHZ uint32
	// Frequency of the channel in Hz.
	FrequencyHZ uint32
	// Signal Power in dB mV.
	SignalPowerDBMV float32
}

// DownstreamConnectionStatus contains Cable Modem Connection status
// pertaining to the downstream channels of the connection.
type DownstreamConnectionStatus struct {
	// Downstream plan for the connection.
	Plan string
	// Primary Downstream channel frequency for the connection.
	FrequencyHZ uint32
	// Primary Downstream channel signal power in dB mV.
	SignalPowerDBMV float32
	// Primary Downstream channel signal SNR in dB.
	SignalSNRDB float32
	// Downstream channel information.
	Channels []DownstreamChannelInfo
}

// UpstreamConnectionStatus contains Cable Modem Connection status
// pertaining to the upstream channels of the connection.
type UpstreamConnectionStatus struct {
	// Primary upstream channel ID.
	ChannelID uint32
	// Upstream channel information.
	Channels []UpstreamChannelInfo
}

// ConnectionStatus contains Cable Modem Connection status.
type ConnectionStatus struct {
	// Current system time on the device when the query was run.
	SystemTime time.Time
	// Duration for which the connection has been up.
	UpTime time.Duration
	// DOCSIS network access status.
	DOCSISNetworkAccessAllowed bool
	// Internet connection status.
	InternetConnected bool
	// Downstream connection status.
	Downstream DownstreamConnectionStatus
	// Upstream connection status.
	Upstream UpstreamConnectionStatus
}

// LogEntry contains Cable Modem Log entry.
type LogEntry struct {
	// Timestamp for this log entry.
	Timestamp time.Time
	// The log string in the entry.
	Log string
}

// CableModemStatus contains detailed status of the Cable Modem.
type CableModemStatus struct {
	// Device related information.
	Info DeviceInfo
	// General settings.
	Settings DeviceSettings
	// Auth settings.
	Auth AuthSettings
	// Software status.
	Software SoftwareStatus
	// Startup status.
	Startup StartupStatus
	// Connection status.
	Connection ConnectionStatus
	// Logs.
	Logs []LogEntry
}
