// Package irkit is golang client for IRKit.
//
// IRKit is IRKit is a Wi-Fi enabled Open Source Infrared Remote Controller device.
// See more on offial documentation http://getirkit.com/en/
//
// See example usage on `v1/_example` directory.
package irkit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	defaultEndpoint = "https://api.getirkit.com"
)

// InternetClient is client for IRKit Internet HTTP API.
// https://api.getirkit.com
type InternetClient struct {
	URL    *url.URL
	client *http.Client
}

type localClient struct {
	// TODO
}

// requestOption stores option used for create
// internet client http.Request.
type requestOption struct {
	// params is key-value that will be added request
	// URL query
	params map[string]string

	// body is key-value that will be added request body
	// as 'key=value' with '&'
	body map[string]string
}

type getKeysResponse struct {
	Clientkey string `json:"clientkey"`
	Deviceid  string `json:"deviceid"`
}

type getDevicesResponse struct {
	Devicekey string `json:"devicekey"`
	Deviceid  string `json:"deviceid"`
}

// DefaultInternetClient creates InternetClient with default API endpoint.
func DefaultInternetClient() *InternetClient {
	client, err := newInternetClient(defaultEndpoint)
	if err != nil {
		// Should not reach here
		panic(err)
	}

	return client
}

// newInternetClient creates InternetClient with the given url.
// If any, return error.
func newInternetClient(rawURL string) (*InternetClient, error) {
	if len(rawURL) == 0 {
		return nil, fmt.Errorf("missing url")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("faile to parse URL: %s", err)
	}

	return &InternetClient{
		URL:    parsedURL,
		client: http.DefaultClient,
	}, nil
}

// newRequest creates request for IntenetClient.
// It returns request with the given context.
func (c *InternetClient) newRequest(ctx context.Context, method,
	spath string, opt *requestOption) (*http.Request, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context")
	}
	if len(method) == 0 {
		return nil, fmt.Errorf("missing method")
	}

	if len(spath) == 0 {
		return nil, fmt.Errorf("missing spath")
	}

	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	// Create request body reader
	var r io.Reader
	if len(opt.body) != 0 {
		kv := make([]string, 0, len(opt.body))
		for k, v := range opt.body {
			kv = append(kv, fmt.Sprintf("%s=%s", k, v))
		}
		r = strings.NewReader(strings.Join(kv, "&"))
	}

	req, err := http.NewRequest(method, u.String(), r)
	if err != nil {
		return nil, err
	}

	// Add query params
	if len(opt.params) != 0 {
		values := req.URL.Query()
		for k, v := range opt.params {
			values.Add(k, v)
		}

		req.URL.RawQuery = values.Encode()
	}

	// Set common headers
	req.Header.Set("User-Agent", "go-irkit")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Set context
	req = req.WithContext(ctx)

	return req, nil
}

// GetKeys gets deviceid and clientkey.
func (c *InternetClient) GetKeys(ctx context.Context,
	token string) (deviceid, clientkey string, err error) {
	if ctx == nil {
		return "", "", fmt.Errorf("nil context")
	}

	if len(token) == 0 {
		return "", "", fmt.Errorf("missing token")
	}

	opt := &requestOption{
		body: map[string]string{
			"clienttoken": token,
		},
	}
	req, err := c.newRequest(ctx, "POST", "/1/keys", opt)
	if err != nil {
		return "", "", err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("invalid status code: %s", res.Status)
	}

	var out getKeysResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&out); err != nil {
		return "", "", err
	}

	return out.Deviceid, out.Clientkey, nil
}

// SendMessages sends IR signal through IRKit device identified by deviceid.
func (c *InternetClient) SendMessages(ctx context.Context,
	clientkey, deviceid string, msg *Message) error {
	if ctx == nil {
		return fmt.Errorf("nil context")
	}

	if len(clientkey) == 0 {
		return fmt.Errorf("missing clientkey")
	}

	if len(deviceid) == 0 {
		return fmt.Errorf("missing deviceid")
	}

	if err := msg.validate(); err != nil {
		return fmt.Errorf("invalid message: %s", err)
	}

	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	opt := &requestOption{
		body: map[string]string{
			"clientkey": clientkey,
			"deviceid":  deviceid,
			"message":   string(buf),
		},
	}

	req, err := c.newRequest(ctx, "POST", "/1/messages", opt)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %s", res.Status)
	}

	return nil
}

// GetMessages gets latest received IR signal. This request is a
// long pooling request.
//
// If you provide clear=true, it it clears IR signal buffer on server
// on internet. When IRKit device receives an IR signal, device sends
// it over to our server on Internet, and server passes it over as response.
//
// Server will respond with an empty response after timeout.
func (c *InternetClient) GetMessages(ctx context.Context,
	clientkey string, clear bool) (*SignalInfo, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context")
	}

	if len(clientkey) == 0 {
		return nil, fmt.Errorf("missing clientkey")
	}

	clearStr := "0"
	if clear {
		clearStr = "1"
	}

	opt := &requestOption{
		params: map[string]string{
			"clear":     clearStr,
			"clientkey": clientkey,
		},
	}

	req, err := c.newRequest(ctx, "GET", "/1/messages", opt)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %s", res.Status)
	}

	var out SignalInfo
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}

// GetDevices gets devicekey and deviceid
func (c *InternetClient) GetDevices(ctx context.Context,
	clientkey string) (devicekey, deviceid string, err error) {
	if ctx == nil {
		return "", "", fmt.Errorf("nil context")
	}

	if len(clientkey) == 0 {
		return "", "", fmt.Errorf("missing clientkey")
	}

	opt := &requestOption{
		body: map[string]string{
			"clientkey": clientkey,
		},
	}

	req, err := c.newRequest(ctx, "POST", "/1/devices", opt)
	if err != nil {
		return "", "", err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("invalid status code: %s", res.Status)
	}

	var out getDevicesResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&out); err != nil {
		return "", "", err
	}

	return out.Devicekey, out.Deviceid, nil
}
