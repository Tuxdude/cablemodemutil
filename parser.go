package cablemodemutil

import (
	"fmt"
)

type actionResponseBody map[string]interface{}

// Parses the raw status returned by the cable modem into the structured cable modem status.
func ParseRawStatus(status CableModemRawStatus) (*CableModemStatus, error) {
	err := validateSubResponses(status)
	if err != nil {
		return nil, fmt.Errorf("invalid status response. reason: %w", err)
	}

	result := CableModemStatus{}
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
