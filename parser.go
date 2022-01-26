package cablemodemutil

import (
	"fmt"
	"log"
)

type actionResponseBody map[string]interface{}

// Parses the raw status returned by the cable modem into the structured cable modem status.
func ParseRawStatus(status CableModemRawStatus) (*CableModemStatus, error) {
	err := validateSubResponses(status)
	if err != nil {
		return nil, fmt.Errorf("invalid status response. reason: %w", err)
	}

	result := CableModemStatus{}
	err = populateDeviceInfo(status, &result.DeviceInfo)
	if err != nil {
		return nil, err
	}
	err = populateDeviceSettings(status, &result.DeviceSettings)
	if err != nil {
		return nil, err
	}
	err = populateAuthSettings(status, &result.AuthSettings)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func actionResp(resp interface{}) actionResponseBody {
	return actionResponseBody(resp.(map[string]interface{}))
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
	val, ok := status[key]
	if !ok {
		return fmt.Errorf("unable to find the response key %q in status response. response: %s", key, prettyPrintJSON(status))
	}

	unpacked := actionResp(val)
	key = actionResultKey(cmd)
	result, ok := unpacked[key].(string)
	if !ok {
		return fmt.Errorf("unable to find the result key %q in status response. response: %s", key, prettyPrintJSON(status))
	}

	if result != "OK" {
		return fmt.Errorf("result in unpacked resposne is %q, expected \"OK\".\nunpacked response: %v", result, prettyPrintJSON(unpacked))
	}
	return nil
}

// Compare the values for the specified keys and emits a warning message if they differ.
func warnIfMismatch(status CableModemRawStatus, desc string, expectedKey string, expectedSubKey string, compareAgainst map[string]string) {
	expected := actionResp(status[expectedKey])[expectedSubKey]

	for key, subKey := range compareAgainst {
		actual := actionResp(status[key])[subKey]
		if expected != actual {
			log.Printf("Warning: %s information mismatch between %q[%q]=%q and %q[%q]=%q", desc, expectedKey, expectedSubKey, expected, key, subKey, actual)
		}
	}
}

// Populates cable modem device information.
func populateDeviceInfo(status CableModemRawStatus, result *CableModemDeviceInfo) error {
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

	warnIfMismatch(status, "Serial Number", "GetArrisRegisterInfoResponse", "SerialNumber", map[string]string{"GetCustomerStatusSoftwareResponse": "StatusSoftwareSerialNum"})
	warnIfMismatch(status, "MAC Address", "GetArrisRegisterInfoResponse", "MacAddress", map[string]string{"GetCustomerStatusSoftwareResponse": "StatusSoftwareMac"})
	return nil
}

// Populates cable modem device settings.
func populateDeviceSettings(status CableModemRawStatus, result *CableModemDeviceSettings) error {
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
func populateAuthSettings(status CableModemRawStatus, result *CableModemAuthSettings) error {
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