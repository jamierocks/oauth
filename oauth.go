// A library for Go client applications that need to perform OAuth authorization against a server,
// typically GitHub.com.
package oauth

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cli/oauth/api"
	"github.com/cli/oauth/device"
)

type httpClient interface {
	PostForm(string, url.Values) (*http.Response, error)
}

type OAuthFlow struct {
	// The host to authorize the app with.
	Hostname string
	// OAuth scopes to request from the user.
	Scopes []string
	// OAuth application ID.
	ClientID string
	// OAuth application secret. Only applicable in web application flow.
	ClientSecret string
	// The localhost URI for web application flow callback, e.g. "http://127.0.0.1/callback".
	CallbackURI string

	// Render an HTML page to the user upon completion of web application flow.
	WriteSuccessHTML func(io.Writer)
	// Open a web browser at a URL. Defaults to opening the default system browser.
	BrowseURL func(string) error

	// The HTTP client to use for API POST requests. Defaults to http.DefaultClient.
	HTTPClient httpClient
	// The stream to listen to keyboard input on. Defaults to os.Stdin.
	Stdin io.Reader
	// The stream to print UI messages to. Defaults to os.Stdout.
	Stdout io.Writer
}

func deviceInitURL(host string) string {
	return fmt.Sprintf("https://%s/login/device/code", host)
}

func webappInitURL(host string) string {
	return fmt.Sprintf("https://%s/login/oauth/authorize", host)
}

func tokenURL(host string) string {
	return fmt.Sprintf("https://%s/login/oauth/access_token", host)
}

// DetectFlow tries to perform Device flow first and falls back to Web application flow.
func (oa *OAuthFlow) DetectFlow() (*api.AccessToken, error) {
	accessToken, err := oa.DeviceFlow()
	if errors.Is(err, device.ErrUnsupported) {
		return oa.WebAppFlow()
	}
	return accessToken, err
}
