package relay

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

// API facade for the Relay REST API.
type Relay struct {
	httpClient *http.Client
	baseURL    string

	sessionToken *string
	hostName     string

	onLibraryRequest func()
	onQueueRequest   func()
}

// Options for creating a relay facade.
type RelayOptions struct {
	// The base URL of the Relay
	BaseURL string

	// Optional HTTP client to use when connecting to the relay.
	// If not provided, a default HTTP client will be used.
	HttpClient *http.Client

	// Optional existing session token. If not provided, the client will create a new session.
	SessionToken *string

	// Optional hostname.
	//
	// If not provided, will default to "Beetroot NeuroFunk" if there is no session token set.
	// Otherwise, will default to the session name on the relay.
	HostName *string

	// Library request handler.
	OnLibraryRequest func()

	// Track queue request handler.
	OnQueueRequest func()
}

// Creates a new relay interface.
func NewRelay(options *RelayOptions) (*Relay, error) {
	if options == nil {
		return nil, fmt.Errorf("options must not be nil")
	}

	if options.OnLibraryRequest == nil {
		return nil, fmt.Errorf("options.OnLibraryRequest must not be nil")
	}

	if options.OnQueueRequest == nil {
		return nil, fmt.Errorf("options.OnQueueRequest must not be nil")
	}

	relay := &Relay{
		baseURL:      options.BaseURL,
		sessionToken: nil,
		hostName:     "Beetroot NeuroFunk",
		httpClient:   nil,
	}

	if options.HttpClient == nil {
		relay.httpClient = http.DefaultClient
	}

	if options.SessionToken != nil {
		err := validateSessionToken(options.HttpClient, *options.SessionToken)
		if err != nil {
			return nil, err
		}

		relay.sessionToken = options.SessionToken
	}

	if options.HostName != nil {
		relay.hostName = *options.HostName
	}

	relay.onLibraryRequest = options.OnLibraryRequest
	relay.onQueueRequest = options.OnLibraryRequest

	return relay, nil
}

// Changes the host's name on the session.
func (relay *Relay) ChangeName(name string) (err error) {
	// TODO: PATCH /player/session
	panic("not implemented")
}

// Deletes the session. If subsequent calls are made after this function, a new session will be made.
func (relay *Relay) DeleteSession() (err error) {
	// TODO: DELETE /player/session

	relay.sessionToken = nil
	panic("not implemented")
}

// Connects to the relay's player socket.
func (relay *Relay) Connect() (conn *Connection, err error) {
	token, err := relay.getActiveSessionToken()
	if err != nil {
		return nil, err
	}

	socket, response, err := websocket.DefaultDialer.Dial(relay.path("/player/connect"), http.Header{
		"User-Agent":      []string{"Beetroot/1 NeuroFunk/0.1"},
		"X-Session-Token": []string{token},
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusContinue {
		if response.StatusCode == http.StatusUnauthorized {
			respBody, err := io.ReadAll(response.Body)
			if err != nil {
				return nil, err
			}

			var respError Error
			err = json.Unmarshal(respBody, &respError)
			if err != nil {
				return nil, err
			}

			return nil, fmt.Errorf("%s: %s (ID: %s)", respError.Code, respError.Message, respError.ID)
		} else {
			return nil, fmt.Errorf("relay responded with unexpected \"%s\"", response.Status)
		}
	}

	return &Connection{
		Relay:  relay,
		Socket: socket,
	}, nil
}

// Returns the full path of an endpoint.
func (relay *Relay) path(path string) string {
	return relay.baseURL + path
}

// Returns the current session token, or if not creates a new one.
func (relay *Relay) getActiveSessionToken() (token string, err error) {
	if relay.sessionToken != nil {
		return *relay.sessionToken, nil
	}

	// TODO: POST /player/session
	panic("not implemented")
}

// Validates the session token and retrieves the session details.
func validateSessionToken(client *http.Client, token string) (err error) {
	// TODO: GET /player/session
	panic("not implemented")
}
