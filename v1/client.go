// Package irkit
package irkit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	defaultEndpoint = "https://api.getirkit.com"
)

// InternetClient is client for IRKit Internet HTTP API.
type InternetClient struct {
	URL    *url.URL
	client *http.Client
}

type LocalClient struct {
}

// Message represents IRKit signal
type Message struct {
	// Format is format of signal. "raw" only.
	Format string `json:"format"`

	// Freq is IRKit sub-carrier frequency. 38 or 40 only. [kHz]
	Freq int `json:"freq"`

	// Data is IRkit signal consists of ON/OFF of sub carrier frequency.
	// IRKit measures On to Off, Off to On interval using a 2MHz counter.
	// data value is an array of those intervals
	Data []int `json:"data"`
}

type RequestOption struct {
	// Body is key-value that will be added request body
	// as 'key=value' with '&'
	Body map[string]string
}

func DefaultInternetClient() *InternetClient {
	client, err := newInternetClient(defaultEndpoint)
	if err != nil {
		// Should not reach here
		panic(err)
	}

	return client
}

func newInternetClient(rawURL string) (*InternetClient, error) {
	if len(rawURL) == 0 {
		return nil, fmt.Errorf("missing url")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("faile to parse URL: %s", err)
	}

	client := &InternetClient{
		URL: parsedURL,
	}

	if err := client.init(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *InternetClient) init() error {
	c.client = http.DefaultClient
	return nil
}

func (c *InternetClient) newRequest(method, spath string, opt *RequestOption) (*http.Request, error) {
	if len(method) == 0 {
		return nil, fmt.Errorf("missing method")
	}

	if len(spath) == 0 {
		return nil, fmt.Errorf("missing spath")
	}

	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	kv := make([]string, 0, len(opt.Body))
	for k, v := range opt.Body {
		kv = append(kv, fmt.Sprintf("%s=%s", k, v))
	}
	r := strings.NewReader(strings.Join(kv, "&"))

	req, err := http.NewRequest(method, u.String(), r)
	if err != nil {
		return nil, err
	}

	// Set common headers
	req.Header.Set("User-Agent", "go-irkit")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

// GetKeys gets deviceid and clientkey.
func (c *InternetClient) GetKeys(ctx context.Context, token string) (deviceid, clientkey string, err error) {
	if ctx == nil {
		return "", "", fmt.Errorf("nil context")
	}

	if len(token) == 0 {
		return "", "", fmt.Errorf("missing token")
	}

	opt := &RequestOption{
		Body: map[string]string{
			"clienttoken": token,
		},
	}
	req, err := c.newRequest("POST", "/1/keys", opt)
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

	out := struct {
		Clientkey string `json:"clientkey"`
		Deviceid  string `json:"deviceid"`
	}{}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&out); err != nil {
		return "", "", err
	}

	return out.Deviceid, out.Clientkey, nil
}

// SendMessages sends IR signal through IRKit device identified by deviceid.
func (c *InternetClient) SendMessages(ctx context.Context, clientkey, deviceid string, msg *Message) error {
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

	opt := &RequestOption{
		Body: map[string]string{
			"clientkey": clientkey,
			"deviceid":  deviceid,
			"message":   string(buf),
		},
	}

	req, err := c.newRequest("POST", "/1/messages", opt)
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

// GetDevices gets devicekey and deviceid
func (c *InternetClient) GetDevices(ctx context.Context, clientkey string) (devicekey, deviceid string, err error) {
	if ctx == nil {
		return "", "", fmt.Errorf("nil context")
	}

	if len(clientkey) == 0 {
		return "", "", fmt.Errorf("missing clientkey")
	}

	opt := &RequestOption{
		Body: map[string]string{
			"clientkey": clientkey,
		},
	}

	req, err := c.newRequest("POST", "/1/devices", opt)
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

	out := struct {
		Devicekey string `json:"devicekey"`
		Deviceid  string `json:"deviceid"`
	}{}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&out); err != nil {
		return "", "", err
	}

	return out.Devicekey, out.Deviceid, nil
}

func (m *Message) validate() error {
	if m.Format != "raw" {
		return fmt.Errorf("format must be raw")
	}

	if m.Freq != 38 && m.Freq != 40 {
		return fmt.Errorf("freq must 38 or 40")
	}

	if len(m.Data) == 0 {
		return fmt.Errorf("empty data")
	}

	return nil
}
