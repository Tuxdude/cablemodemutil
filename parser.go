package cablemodemutil

import (
	"fmt"
	"log"
	"strings"
)

type actionResponseBody map[string]interface{}

// ParseRawStatus parses the raw status returned by the cable modem into the structured cable modem status.
func ParseRawStatus(status CableModemRawStatus) (*CableModemStatus, error) {
	err := validateSubResponses(status)
	if err != nil {
		return nil, fmt.Errorf("invalid status response. reason: %w", err)
	}

	result := CableModemStatus{}
	err = populateDeviceInfo(status, &result.Info)
	if err != nil {
		return nil, err
	}
	err = populateDeviceSettings(status, &result.Settings)
	if err != nil {
		return nil, err
	}
	err = populateAuthSettings(status, &result.Auth)
	if err != nil {
		return nil, err
	}
	err = populateSoftwareStatus(status, &result.Software)
	if err != nil {
		return nil, err
	}
	err = populateStartupStatus(status, &result.Startup)
	if err != nil {
		return nil, err
	}
	err = populateConnectionStatus(status, &result.Connection)
	if err != nil {
		return nil, err
	}
	result.Logs, err = populateLogEntries(status)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func actionResp(resp interface{}) actionResponseBody {
	body, ok := resp.(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("Cannot convert given type to actionResponseBody, data:%v", resp))
	}
	return actionResponseBody(body)
}

// Validates all the sub-responses within the status response were successful and have the expected payload.
func validateSubResponses(status CableModemRawStatus) error {
	for _, cmd := range statusSubCommands {
		err := validateSubResponse(status, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}

// Validates a specific command's sub-response within the status response.
func validateSubResponse(status CableModemRawStatus, cmd string) error {
	key := actionResponseKey(cmd)
	val, keyExists := status[key]
	if !keyExists {
		return fmt.Errorf(
			"unable to find the response key %q in status response. response: %s",
			key,
			prettyPrintJSON(status),
		)
	}

	unpacked := actionResp(val)
	key = actionResultKey(cmd)
	result, keyExists := unpacked[key].(string)
	if !keyExists {
		return fmt.Errorf(
			"unable to find the result key %q in status response. response: %s",
			key,
			prettyPrintJSON(status),
		)
	}

	if result != "OK" {
		return fmt.Errorf(
			"result in unpacked resposne is %q, expected \"OK\".\nunpacked response: %v",
			result,
			prettyPrintJSON(unpacked),
		)
	}
	return nil
}

// Compare the values for the specified keys and emits a warning message if they differ.
func warnIfMismatch(
	status CableModemRawStatus,
	desc string,
	expectedKey string,
	expectedSubKey string,
	compareAgainst map[string]string,
) {
	expected := actionResp(status[expectedKey])[expectedSubKey]

	for key, subKey := range compareAgainst {
		actual := actionResp(status[key])[subKey]
		if expected != actual {
			log.Printf(
				"Warning: %s information mismatch between %q[%q]=%q and %q[%q]=%q",
				desc,
				expectedKey,
				expectedSubKey,
				expected,
				key,
				subKey,
				actual,
			)
		}
	}
}

// Populates cable modem device information.
func populateDeviceInfo(status CableModemRawStatus, result *DeviceInfo) error {
	var err error
	data := actionResp(status["GetArrisRegisterInfoResponse"])

	result.Model, err = parseString(data, "ModelName", "Model Name")
	if err != nil {
		return err
	}
	result.SerialNumber, err = parseString(data, "SerialNumber", "Serial Number")
	if err != nil {
		return err
	}
	result.MACAddress, err = parseString(data, "MacAddress", "MAC Address")
	if err != nil {
		return err
	}

	warnIfMismatch(
		status,
		"Serial Number",
		"GetArrisRegisterInfoResponse",
		"SerialNumber",
		map[string]string{
			"GetCustomerStatusSoftwareResponse": "StatusSoftwareSerialNum",
		},
	)
	warnIfMismatch(
		status,
		"MAC Address",
		"GetArrisRegisterInfoResponse",
		"MacAddress",
		map[string]string{
			"GetCustomerStatusSoftwareResponse": "StatusSoftwareMac",
		},
	)
	return nil
}

// Populates cable modem device settings.
func populateDeviceSettings(status CableModemRawStatus, result *DeviceSettings) error {
	var err error
	conf := actionResp(status["GetArrisConfigurationInfoResponse"])
	reg := actionResp(status["GetArrisRegisterStatusResponse"])

	result.FrontPanelLightsOn, err = parseBool(conf, "LedStatus", "1", "LED Status")
	if err != nil {
		return err
	}
	result.EnergyEfficientEthernetOn, err = parseBool(conf, "ethSWEthEEE", "1", "Energy Efficient Ethernet")
	if err != nil {
		return err
	}
	result.AskMeLater, err = parseBool(reg, "AskMeLater", "1", "Ask Me Later")
	if err != nil {
		return err
	}
	result.NeverAsk, err = parseBool(reg, "NeverAsk", "1", "Never Ask")
	if err != nil {
		return err
	}

	return nil
}

// Populates cable modem auth settings.
func populateAuthSettings(status CableModemRawStatus, result *AuthSettings) error {
	var err error
	acc := actionResp(status["GetCustomerStatusSecAccountResponse"])

	result.CurrentLogin, err = parseString(acc, "CurrentLogin", "Current Login")
	if err != nil {
		return err
	}
	result.CurrentNameAdmin, err = parseString(acc, "CurrentNameAdmin", "Current Admin Username")
	if err != nil {
		return err
	}
	result.CurrentNameUser, err = parseString(acc, "CurrentNameUser", "Current Username")
	if err != nil {
		return err
	}
	result.CurrentPasswordAdmin, err = parseString(acc, "CurrentPwAdmin", "Current Admin Password")
	if err != nil {
		return err
	}
	result.CurrentPasswordUser, err = parseString(acc, "CurrentPwUser", "Current User Password")
	if err != nil {
		return err
	}

	return nil
}

// Populates cable modem software status.
func populateSoftwareStatus(status CableModemRawStatus, result *SoftwareStatus) error {
	var err error
	sw := actionResp(status["GetCustomerStatusSoftwareResponse"])
	result.FirmwareVersion, err = parseString(sw, "StatusSoftwareSfVer", "Firmware Version")
	if err != nil {
		return err
	}
	result.CertificateInstalled, err = parseBool(sw, "StatusSoftwareCertificate", "Installed", "Certificate Installed")
	if err != nil {
		return err
	}
	result.CustomerVersion, err = parseString(sw, "StatusSoftwareCustomerVer", "Customer Version")
	if err != nil {
		return err
	}
	result.HDVersion, err = parseString(sw, "StatusSoftwareHdVer", "HD Version")
	if err != nil {
		return err
	}
	result.DOCSISSpecVersion, err = parseString(sw, "StatusSoftwareSpecVer", "DOCSIS Spec Version")
	if err != nil {
		return err
	}

	return nil
}

// Populates cable modem startup status.
// nolint:funlen
func populateStartupStatus(status CableModemRawStatus, result *StartupStatus) error {
	var err error
	startup := actionResp(status["GetCustomerStatusStartupSequenceResponse"])

	result.Boot.Status, err = parseBool(startup, "CustomerConnBootStatus", "OK", "Boot Status")
	if err != nil {
		return err
	}
	result.Boot.Operational, err = parseBool(
		startup,
		"CustomerConnBootComment",
		"Operational",
		"Boot Comment",
	)
	if err != nil {
		return err
	}
	result.ConfigFile.Status, err = parseBool(
		startup,
		"CustomerConnConfigurationFileStatus",
		"OK",
		"Configuration File Status",
	)
	if err != nil {
		return err
	}
	result.ConfigFile.Comment, err = parseString(
		startup,
		"CustomerConnConfigurationFileComment",
		"Configuration File Comment",
	)
	if err != nil {
		return err
	}
	result.Connectivity.Status, err = parseBool(
		startup,
		"CustomerConnConnectivityStatus",
		"OK",
		"Connectivity Status",
	)
	if err != nil {
		return err
	}
	result.Connectivity.Operational, err = parseBool(
		startup,
		"CustomerConnConnectivityComment",
		"Operational",
		"Connectivity Comment",
	)
	if err != nil {
		return err
	}
	result.Downstream.FrequencyHZ, err = parseFreq(
		startup, "CustomerConnDSFreq", true, "Downstream Connection Frequency")
	if err != nil {
		return err
	}
	result.Downstream.Locked, err = parseBool(
		startup,
		"CustomerConnDSComment",
		"Locked",
		"Downstream Connection Comment",
	)
	if err != nil {
		return err
	}
	result.Security.Enabled, err = parseBool(
		startup,
		"CustomerConnSecurityStatus",
		"Enabled",
		"Security Status",
	)
	if err != nil {
		return err
	}
	result.Security.Comment, err = parseString(
		startup,
		"CustomerConnSecurityComment",
		"Security Comment",
	)
	if err != nil {
		return err
	}
	return nil
}

// populates cable modem connection status.
// nolint:funlen
func populateConnectionStatus(status CableModemRawStatus, result *ConnectionStatus) error {
	var err error
	conn := actionResp(status["GetCustomerStatusConnectionInfoResponse"])
	dev := actionResp(status["GetArrisDeviceStatusResponse"])
	config := actionResp(status["GetArrisConfigurationInfoResponse"])

	result.SystemTime, err = parseSystemTimestamp(conn, "CustomerCurSystemTime", "Current System Time")
	if err != nil {
		return err
	}
	result.UpTime, err = parseDuration(conn, "CustomerConnSystemUpTime", "System Up Time")
	if err != nil {
		return err
	}
	result.EstablishedAt = result.SystemTime.Add(-result.UpTime)
	result.DOCSISNetworkAccessAllowed, err = parseBool(
		conn,
		"CustomerConnNetworkAccess",
		"Allowed",
		"DOCSIS Network Access",
	)
	if err != nil {
		return err
	}
	result.InternetConnected, err = parseBool(
		dev,
		"InternetConnection",
		"Connected",
		"Internet Connection Status",
	)
	if err != nil {
		return err
	}
	result.DownstreamPlan, err = parseString(config, "DownstreamPlan", "Downstream Plan")
	if err != nil {
		return err
	}
	result.DownstreamFrequencyHZ, err = parseFreq(config, "DownstreamFrequency", false, "Downstream Frequency")
	if err != nil {
		return err
	}
	result.DownstreamSignalPowerDBMV, err = parseSignalPowerInt(
		dev,
		"DownstreamSignalPower",
		true,
		"Downstream Signal Power",
	)
	if err != nil {
		return err
	}
	result.DownstreamSignalSNRDB, err = parseSignalSNR(dev, "DownstreamSignalSnr", true, "Downstream Signal SNR")
	if err != nil {
		return err
	}
	result.UpstreamChannelID, err = parseChannelID(config, "UpstreamChannelId", "Upstream Channel ID")
	if err != nil {
		return err
	}
	result.DownstreamChannels, err = populateDownstreamChannels(status)
	if err != nil {
		return err
	}
	result.UpstreamChannels, err = populateUpstreamChannels(status)
	if err != nil {
		return err
	}

	// TODO: Verify downstream frequency is the same in all three places below..
	// GetArrisConfigurationInfoResponse.DownstreamFrequency (no HZ suffix in string)
	// GetArrisDeviceStatusResponse.DownstreamFrequency (has HZ suffix in string)
	// GetCustomerStatusStartupSequenceResponse.CustomerConnDSFreq (has HZ suffix in string)

	return nil
}

// Populates cable modem downstream channel information.
func populateDownstreamChannels(status CableModemRawStatus) ([]DownstreamChannelInfo, error) {
	var err error
	dsInfo := actionResp(status["GetCustomerStatusDownstreamChannelInfoResponse"])
	squashedRows, err := parseString(dsInfo, "CustomerConnDownstreamChannel", "Downstream Channel info")
	if err != nil {
		return nil, err
	}

	// Each row is delimited by a '|+|'
	rows := strings.Split(squashedRows, "|+|")
	result := make([]DownstreamChannelInfo, len(rows))
	for i, row := range rows {
		// Each column is delimited by a '^'
		cols := strings.Split(row, "^")
		// The columns are:
		// Row ID, Lock Status, Modulation, Channel ID, Frequency, Power, SNR, Corrected Err, Uncorrected Err, Blank
		if len(cols) != 10 {
			return nil, fmt.Errorf(
				"expected 10 columns in a downstream channel, actual %d row=%q",
				len(cols),
				row,
			)
		}

		result[i].LockStatus = cols[1]
		result[i].Modulation = cols[2]
		result[i].ChannelID, err = parseChannelIDStr(cols[3], "Downstream Channel ID")
		if err != nil {
			return nil, err
		}
		result[i].FrequencyHZ, err = parseFreqStr(cols[4], false, "Downstream Channel Frequency")
		if err != nil {
			return nil, err
		}
		result[i].SignalPowerDBMV, err = parseSignalPowerIntStr(cols[5], false, "Downstream Channel Signal Power")
		if err != nil {
			return nil, err
		}
		result[i].SignalSNRMERDB, err = parseSignalSNRStr(cols[6], false, "Downstream Channel Signal SNR/MER")
		if err != nil {
			return nil, err
		}
		result[i].CorrectedErrors, err = parseSignalErrorsStr(cols[7], "Downstream Channel Signal Corrected Errors")
		if err != nil {
			return nil, err
		}
		result[i].UncorrectedErrors, err = parseSignalErrorsStr(cols[8], "Downstream Channel Signal Uncorrected Errors")
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Populates cable modem upstream channel information.
func populateUpstreamChannels(status CableModemRawStatus) ([]UpstreamChannelInfo, error) {
	var err error
	usInfo := actionResp(status["GetCustomerStatusUpstreamChannelInfoResponse"])
	squashedRows, err := parseString(usInfo, "CustomerConnUpstreamChannel", "Upstream Channel info")
	if err != nil {
		return nil, err
	}

	// Each row is delimited by a '|+|'
	rows := strings.Split(squashedRows, "|+|")
	result := make([]UpstreamChannelInfo, len(rows))
	for i, row := range rows {
		// Each column is delimited by a '^'
		cols := strings.Split(row, "^")
		// The columns are:
		// Row ID, Lock Status, Modulation, Channel ID, Width, Frequency, Power, Blank
		if len(cols) != 8 {
			return nil, fmt.Errorf(
				"expected 8 columns in a upstream channel, actual %d row=%q",
				len(cols),
				row,
			)
		}

		result[i].LockStatus = cols[1]
		result[i].Modulation = cols[2]
		result[i].ChannelID, err = parseChannelIDStr(cols[3], "Upstream Channel ID")
		if err != nil {
			return nil, err
		}
		result[i].WidthHZ, err = parseFreqStr(cols[4], false, "Upstream Channel Width")
		if err != nil {
			return nil, err
		}
		result[i].FrequencyHZ, err = parseFreqStr(cols[5], false, "Upstream Channel Frequency")
		if err != nil {
			return nil, err
		}
		result[i].SignalPowerDBMV, err = parseSignalPowerFloatStr(cols[6], false, "Upstream Channel Signal Power")
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Populates cable modem log entries.
func populateLogEntries(status CableModemRawStatus) ([]LogEntry, error) {
	var err error
	usInfo := actionResp(status["GetCustomerStatusLogResponse"])
	squashedRows, err := parseString(usInfo, "CustomerStatusLogList", "Log list")
	if err != nil {
		return nil, err
	}

	// Each row is delimited by a '}-{'
	rows := strings.Split(squashedRows, "}-{")
	result := make([]LogEntry, len(rows))
	for i, row := range rows {
		// Each column is delimited by a '^'
		cols := strings.Split(row, "^")
		// The columns are:
		// 0, Time, Date, 3, Log
		if len(cols) != 5 {
			return nil, fmt.Errorf(
				"expected 5 columns in a log entry, actual %d row=%q",
				len(cols),
				row,
			)
		}

		result[i].Timestamp, err = parseLogTimestamp(cols[2], cols[1])
		if err != nil {
			return nil, err
		}
		result[i].Log = parseLogEntry(cols[4])
	}

	// TODO:
	// 1. Count number of success/failed login attempts.
	//		  Expose timestamps and IPs of successful and failed login attempts.
	// 2. Count other kinds of errors (three known categories so far)
	//        STARTED_UNICAST_MAINTENANCE_RANGING_NO_RESPONSE_RECEIVED
	//        RNG_RSP_CCAP_COMMAND_POWER_EXCEEDS_TOP_OF_DRW
	//        DYNAMIC_RANGE_WINDOW_VIOLATION
	// 3. Parse CMTS MAC from the latest log entry with the info (if available).

	return result, nil
}
