package source

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

var (
	proxyURL    string
	rwMux       sync.RWMutex
	statusCodes = map[string]int{
		http.MethodPost:   http.StatusCreated,
		http.MethodGet:    http.StatusOK,
		http.MethodPut:    http.StatusOK,
		http.MethodPatch:  http.StatusOK,
		http.MethodDelete: http.StatusNoContent,
	}
)

type Requester interface {
	Method() string
	RequestURL() string
	Body() io.Reader
}

func SetProxy(rawURL string) {
	_, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("set proxy failed: %s", err)
		return
	}
	rwMux.Lock()
	proxyURL = rawURL
	rwMux.Unlock()
}

func proxyTransport() http.RoundTripper {
	transport := &http.Transport{}
	rwMux.RLock()
	defer rwMux.RUnlock()
	rawURL, err := url.Parse(proxyURL)
	if err != nil {
		fmt.Printf("proxy url is wrong: %s", err)
		return transport
	}

	switch rawURL.Scheme {
	case "http", "https":
		transport.Proxy = http.ProxyURL(rawURL)
	case "socks5":
		dialer, err := proxy.SOCKS5("tcp", rawURL.Host, nil, proxy.Direct)
		if err != nil {
			fmt.Printf("can't connect to the socks proxy: %s", err)
			return transport
		}

		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	}
	return transport
}

func Do(r Requester, data interface{}) error {
	req, err := http.NewRequest(r.Method(), r.RequestURL(), r.Body())
	if err != nil {
		return err
	}

	client := http.Client{
		Transport: proxyTransport(),
	}

	var resp *http.Response
	for num := 10; num > 0; num-- {
		resp, err = client.Do(req)
		if err == nil || err != io.EOF {
			break
		}
		time.Sleep(time.Second * 3)
	}

	if err != nil {
		return fmt.Errorf("do request failed: %s", err)
	}

	expected := statusCodes[r.Method()]
	if resp.StatusCode != expected {
		return fmt.Errorf("response status code is not matched, expected: %d, but %d", expected, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %s", err)
	}

	err = json.Unmarshal(body[5:], data)
	if err != nil {
		return fmt.Errorf("unmarshal json body failed: %s", err)
	}
	return nil
}
