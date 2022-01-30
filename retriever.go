package cablemodemutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const (
	urlFormat           = "%s://%s/HNAP1/"
	tokenExpiryDuration = 10 * time.Minute
	loginAction         = "Login"
	queryAction         = "GetMultipleHNAPs"
)

// nolint:gochecknoglobals
var statusSubCommands = []string{
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

// Retriever is used to retrieve the current status of the Cable Modem.
type Retriever struct {
	client        *httpClient
	username      string
	clearPassword string
	debug         RetrieverDebug
	tok           *token
	tokMu         sync.Mutex
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
	// Debugging options.
	Debug RetrieverDebug
}

// RetrieverDebug is used to specify the debugging options of the Retriever.
type RetrieverDebug struct {
	// If set to true logs additional debug information except for the requests and responses sent/received to/from the cable modem, false otherwise.
	Debug bool
	// If set to true logs additional debug information about the requests sent to the cable modem, false otherwise.
	DebugReq bool
	// If set to true logs additional debug information about the resposnes received from the cable modem, false otherwise.
	DebugResp bool
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
func encodePayload(action string, payload actionRequest) (*bytes.Buffer, error) {
	req := soapRequest{
		action: payload,
	}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request as JSON, reason: %w", err)
	}
	return buf, nil
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

// NewStatusRetriever returns a new retriever object that can be used to query the Cable Modem status.
func NewStatusRetriever(input *RetrieverInput) *Retriever {
	url := fmt.Sprintf(urlFormat, input.Protocol, input.Host)
	r := Retriever{}
	r.client = newHTTPClient(url, input.SkipVerifyCert, &input.Debug)
	r.username = input.Username
	r.clearPassword = input.ClearPassword
	r.debug = input.Debug
	r.tok = resetToken()
	return &r
}

// Persist the token in the object.
func (r *Retriever) persistToken(tok *token) {
	r.tokMu.Lock()
	r.tok = &token{
		privateKey: tok.privateKey,
		uid:        tok.uid,
		expiry:     tok.expiry,
	}
	r.tokMu.Unlock()
	if r.debug.Debug {
		fmt.Println("Persisting a new token.")
		debugToken(tok)
	}
}

// Retrieves a copy of the persisted token.
func (r *Retriever) getToken() *token {
	r.tokMu.Lock()
	res := &token{
		privateKey: r.tok.privateKey,
		uid:        r.tok.uid,
		expiry:     r.tok.expiry,
	}
	r.tokMu.Unlock()
	return res
}

// Sends the SOAP request for the specified action containing the specified payload.
func (r *Retriever) sendReq(action string, payload actionRequest, tok *token) (actionResponse, error) {
	req, err := encodePayload(action, payload)
	if err != nil {
		return nil, err
	}
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
	hashedPassword, err := genHashedPassword(tok.privateKey, challenge)
	if err != nil {
		return fmt.Errorf("auth failed while generating hashed password, reason: %w", err)
	}

	payload := actionRequest{
		"LoginPassword": hashedPassword,
		"Captcha":       "",
		"PrivateLogin":  "LoginPassword",
		"Action":        "login",
		"Username":      r.username,
	}
	_, err = r.sendReq(loginAction, payload, tok)
	if err != nil {
		return fmt.Errorf("auth failed.\nreason: %w", err)
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

	privateKey, err := genPrivateKey(loginResp.publicKey, loginResp.challenge, r.clearPassword)
	if err != nil {
		return nil, fmt.Errorf("login failed while generating private key, reason: %w", err)
	}

	tok := &token{
		uid:        loginResp.uid,
		privateKey: privateKey,
		expiry:     expiry,
	}
	err = r.doAuth(loginResp.challenge, tok)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

// RawStatus retrieves the current detailed raw status from the cable modem.
func (r *Retriever) RawStatus() (CableModemRawStatus, error) {
	var err error
	loginAttempted := false
	payload := make(actionRequest)
	for _, cmd := range statusSubCommands {
		payload[cmd] = ""
	}

	tok := r.getToken()
	for {
		// If the token has expired, login to generate a fresh token.
		if time.Now().After(tok.expiry) {
			if r.debug.Debug {
				fmt.Println("Token expired, will attempt a new login.")
				debugToken(tok)
			}
			loginAttempted = true
			tok, err = r.login()
			if err != nil {
				return nil, err
			}
			r.persistToken(tok)
		}

		var st actionResponse
		// Compute the new token expiry time based on when we send the request.
		newExpiry := time.Now().Add(tokenExpiryDuration)
		// Fetch the current status.
		st, err = r.sendReq(queryAction, payload, tok)
		if err == nil {
			tok.expiry = newExpiry
			r.persistToken(tok)
			return CableModemRawStatus(st), nil
		}

		// If there is a failure in fetching the current status, and we
		// didn't generate a fresh token just now, generate a new token
		// and re-attempt fetching the current status.
		if loginAttempted {
			break
		}
		tok = resetToken()
	}
	return nil, err
}

// Status retrieves and parses the current detailed status from the cable modem.
func (r *Retriever) Status() (*CableModemStatus, error) {
	raw, err := r.RawStatus()
	if err != nil {
		return nil, err
	}
	return ParseRawStatus(raw)
}
