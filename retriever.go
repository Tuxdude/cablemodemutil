package cablemodemutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

const (
	urlFormat           = "%s://%s/HNAP1/"
	tokenExpiryDuration = 10 * time.Minute
	loginAction         = "Login"
	queryAction         = "GetMultipleHNAPs"
)

var (
	statusSubCommands = []string{
		// MAC address, serial number and model info.
		"GetArrisRegisterInfo",
		// Software version info along with device MAC addr / serial.
		"GetCustomerStatusSoftware",
		// Short Device Status and signal info.
		"GetArrisDeviceStatus",
		// Current Time, uptime and short connection status.
		"GetCustomerStatusConnectionInfo",
		// Detailed connection status.
		"GetCustomerStatusStartupSequence",
		// Downstream channel info.
		"GetCustomerStatusDownstreamChannelInfo",
		// Upstream channel info.
		"GetCustomerStatusUpstreamChannelInfo",
		// Downstream/Upstream Frequency summary and configurable settings.
		"GetArrisConfigurationInfo",
		// Event log.
		"GetCustomerStatusLog",
		// User login/password information (Not so useful).
		"GetCustomerStatusSecAccount",
		// Ask me later and never ask (Not so useful).
		"GetArrisRegisterStatus",
		// Just contains 'XXX' (Not useful).
		"GetCustomerStatusXXX",
		// Just contains 'XXX' (Not useful).
		"GetArrisXXX",
		// Just contains 'XXX' (Not useful).
		"GetCustomerStatusLogXXX",
	}
)

// Retriever is used to retrieve the current status of the Cable Modem.
type Retriever struct {
	client        *httpClient
	username      string
	clearPassword string
}

// RetrieverInput is used to specify the input for building a Retriever.
type RetrieverInput struct {
	// The host name or IP address of the cable modem device.
	Host string
	// The protocol used to connect to the cable modem, either "http" or "https".
	Protocol string
	// If true skips verifying the cable modem's SSL certificate, false otherwise.
	SkipVerifyCert bool
	// User name for authenticating with the cable modem.
	Username string
	// Password for authenticating with the cable modem.
	ClearPassword string
}

// The token object containing the state of the authenticated session with the cable modem.
type token struct {
	// The UID of the session provided by the cable modem during authentication.
	uid string
	// The private key of the session after authentication, generated based on public key, challenge from the cable modem and the supplied password.
	privateKey string
	// The expiry timestamp of the credentials stored in this session.
	expiry time.Time
}

// Initial response from the cable modem to allow the client to initiate authentication.
type loginResponse struct {
	// The UID of the session provided by the cable modem during authentication.
	uid string
	// The public key provided by the cable modem during authentication.
	publicKey string
	// The challenge message provided by the cable modem during authentication.
	challenge string
}

// actionRequest represents the payload of the request containing the SOAP action command.
type actionRequest map[string]string

// actionResponse represents the payload of the response to a SOAP action command request.
type actionResponse map[string]interface{}

// soapRequest represents the full SOAP request body.
type soapRequest map[string]actionRequest

// soapResponse represents the full SOAP response body to a SOAP request sent.
type soapResponse map[string]actionResponse

// Returns a token that has been reset to the initial state.
func resetToken() *token {
	return &token{
		privateKey: "withoutLoginKey",
	}
}

// Encodes the specified payload for the specified action into a byte buffer to send this as a request.
func encodePayload(action string, payload actionRequest) *bytes.Buffer {
	req := soapRequest{
		action: payload,
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(req)
	return buf
}

// Decodes the specified byte array response for the specified action into the response payload.
func decodePayload(action string, resp *[]byte) (actionResponse, error) {
	var payload soapResponse
	err := json.Unmarshal(*resp, &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response payload, reason:%w", err)
	}
	return unpackResponse(action, payload)
}

// Validates and unpacks the response for the specified SOAP action.
func unpackResponse(action string, resp soapResponse) (actionResponse, error) {
	if len(resp) != 1 {
		return nil, fmt.Errorf(
			"action: %s, invalid number of keys (%d) in response, expected 1.\nresponse: %v", action, len(resp), prettyPrintJSON(resp))
	}

	respKey := actionResponseKey(action)
	unpacked, ok := resp[respKey]
	if !ok {
		return nil, fmt.Errorf(
			"action: %s, unable to find the response key %q in response.\nresponse: %v", action, respKey, prettyPrintJSON(resp))
	}

	resultKey := actionResultKey(action)
	result, ok := unpacked[resultKey]
	if !ok {
		return nil, fmt.Errorf(
			"action: %s, unable to find the result key %q in unpacked response.\nunpacked response: %v", action, resultKey, prettyPrintJSON(unpacked))
	}
	if result != "OK" {
		return nil, fmt.Errorf(
			"action: %s, result in unpacked resposne is %q, expected \"OK\".\nunpacked response: %v", action, result, prettyPrintJSON(unpacked))
	}

	return unpacked, nil
}

// Builds a new Retriever object to query the Cable Modem.
func NewStatusRetriever(input *RetrieverInput) *Retriever {
	url := fmt.Sprintf(urlFormat, input.Protocol, input.Host)
	r := Retriever{}
	r.client = newHttpClient(url, input.SkipVerifyCert)
	r.username = input.Username
	r.clearPassword = input.ClearPassword
	return &r
}

// Sends the SOAP request for the specified action containing the specified payload.
func (r *Retriever) sendReq(action string, payload actionRequest, tok *token) (actionResponse, error) {
	req := encodePayload(action, payload)
	resp, err := r.client.sendPOST(action, req, tok)
	if err != nil {
		return nil, err
	}
	return decodePayload(action, resp)
}

// Retrieves the cookie, public key and challenge information from the cable modem that can be used for initiating an authentication request.
func (r *Retriever) getLoginResponse() (*loginResponse, error) {
	payload := actionRequest{
		"LoginPassword": "",
		"Captcha":       "",
		"PrivateLogin":  "LoginPassword",
		"Action":        "request",
		"Username":      r.username,
	}
	tok := resetToken()
	resp, err := r.sendReq(loginAction, payload, tok)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve login challenge\nreason: %w", err)
	}

	var ok bool
	result := loginResponse{}
	result.uid, ok = resp["Cookie"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"unable to find key 'Cookie' in login response\nresponse: %v", prettyPrintJSON(resp))
	}

	result.publicKey, ok = resp["PublicKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"unable to find key 'PublicKey' in login response\nresponse: %v", prettyPrintJSON(resp))
	}

	result.challenge, ok = resp["Challenge"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"unable to find key 'Challenge' in login response\nresponse: %v", prettyPrintJSON(resp))
	}

	return &result, nil
}

// Performs authentication with the cable modem and returns the response.
func (r *Retriever) doAuth(challenge string, tok *token) error {
	hashedPassword := genHashedPassword(tok.privateKey, challenge, r.clearPassword)
	payload := actionRequest{
		"LoginPassword": hashedPassword,
		"Captcha":       "",
		"PrivateLogin":  "LoginPassword",
		"Action":        "login",
		"Username":      r.username,
	}

	_, err := r.sendReq(loginAction, payload, tok)
	if err != nil {
		return fmt.Errorf("login failed.\nreason: %w", err)
	}
	return nil
}

// Login to the cable modem using the specified username and password.
func (r *Retriever) login() (*token, error) {
	loginResp, err := r.getLoginResponse()
	if err != nil {
		return nil, err
	}
	// Compute the expiry time as soon as we obtain the response.
	expiry := time.Now().Add(tokenExpiryDuration)

	privateKey := genPrivateKey(loginResp.publicKey, loginResp.challenge, r.clearPassword)
	tok := &token{
		uid:        loginResp.uid,
		privateKey: privateKey,
		expiry:     expiry,
	}
	err = r.doAuth(loginResp.challenge, tok)
	if err != nil {
		return nil, fmt.Errorf("login failed.\nreason: %w", err)
	}
	return tok, nil
}

// Retrieves the current detailed raw status from the cable modem.
func (r *Retriever) RawStatus() (CableModemRawStatus, error) {
	payload := make(actionRequest)
	for _, cmd := range statusSubCommands {
		payload[cmd] = ""
	}

	tok, err := r.login()
	if err != nil {
		return nil, err
	}

	// Fetch the current status.
	st, err := r.sendReq(queryAction, payload, tok)
	if err != nil {
		return nil, err
	}
	return CableModemRawStatus(st), nil
}
