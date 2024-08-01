package utils

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	"weatherbot/config"
	"weatherbot/internal/logger"
	"weatherbot/internal/weather"
)

type Auth struct {
	User     string
	Password string
}

type HTTPProxyDialer struct {
	ProxyURL *url.URL
}

type RequestParams struct {
	Method      string
	Url         string
	QueryParams *map[string]string
	Headers     *map[string]string
	Body        interface{}
}

const httpClientTimeOut = 30 * time.Second
const Retries = 3
const RetryTimeout = 3 * time.Second

func NewRequest(params *RequestParams) (*http.Request, error) {
	u, err := url.Parse(params.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid Url: %v", err)
	}

	if params.Method == http.MethodGet && len(*params.QueryParams) > 0 {
		q := u.Query()
		for key, val := range *params.QueryParams {
			q.Add(key, val)
		}
		u.RawQuery = q.Encode()
	}

	var reqBody io.Reader
	switch body := params.Body.(type) {
	case nil:
		// do nothing
	case url.Values:
		// form body
		reqBody = strings.NewReader(body.Encode())
	case []byte:
		// json
		reqBody = bytes.NewReader(body)
	default:
		return nil, fmt.Errorf("unsupported body type: %T", body)
	}

	req, err := http.NewRequest(params.Method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	if params.Method == http.MethodPost {
		if _, ok := params.Body.(url.Values); ok {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	if params.Headers != nil {
		for key, val := range *params.Headers {
			req.Header.Set(key, val)
		}
	}

	return req, nil
}

func (d *HTTPProxyDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	dialer := &net.Dialer{}
	return dialer.DialContext(ctx, network, addr)
}

func GetDialer(proxyURL string, auth *proxy.Auth) (proxy.ContextDialer, error) {
	uri, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proxy URL: %v", err)
	}

	switch uri.Scheme {
	case "http", "https":
		return getHTTPDialer(uri, auth)
	case "socks5":
		return getSOCKSDialer(uri, auth)
	default:
		return nil, fmt.Errorf("unsupported proxy scheme: %s", uri.Scheme)
	}
}

func getHTTPDialer(uri *url.URL, auth *proxy.Auth) (proxy.ContextDialer, error) {
	return &HTTPProxyDialer{
		ProxyURL: uri,
	}, nil
}

func getSOCKSDialer(uri *url.URL, auth *proxy.Auth) (proxy.ContextDialer, error) {
	socksAuth := &proxy.Auth{
		User:     auth.User,
		Password: auth.Password,
	}
	dialer, err := proxy.SOCKS5("tcp", uri.Host, socksAuth, proxy.Direct)
	if err != nil {
		return nil, err
	}

	contextDialer, ok := dialer.(proxy.ContextDialer)
	if !ok {
		return nil, fmt.Errorf("SOCKS5 dialer does not support ContextDialer")
	}

	return contextDialer, nil
}

func getHttpClient() *http.Client {
	client := &http.Client{}
	if proxyURL := config.GetConfigValue("PROXY_URL"); proxyURL != "" {
		proxyURI, err := url.Parse(proxyURL)
		if err != nil {
			logger.Logger().Fatalf("Failed to parse proxy URL: %v", err)
		}
		// Create a SOCKS5 dialer
		auth := &proxy.Auth{
			User:     proxyURI.User.Username(),
			Password: getPassword(proxyURI),
		}
		dialer, err := GetDialer(proxyURL, auth)
		if err != nil {
			logger.Logger().Fatalf("Failed to create dialer: %v", err)
		}

		// Create transport that uses the SOCKS5 dialer
		transport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, addr)
			},
		}
		client.Transport = transport
		client.Timeout = httpClientTimeOut
	}

	return client
}

func getPassword(proxyURI *url.URL) string {
	password, _ := proxyURI.User.Password()
	return password
}

func DoRequestWithRetry(req *http.Request, maxRetires int, initialWait time.Duration) (*http.Response, error) {
	var response *http.Response
	var err error

	wait := initialWait
	client := getHttpClient()

	for attempt := 0; attempt < maxRetires; attempt++ {
		if req.Body != nil {
			if req.Body, err = req.GetBody(); err != nil {
				return nil, fmt.Errorf("failed to get request body: %w", err)
			}
		}

		response, err = client.Do(req)
		if err == nil && response.StatusCode == http.StatusOK {
			return response, nil
		}

		if response != nil {
			response.Body.Close()
		}

		time.Sleep(wait)
		wait *= 2
	}

	return nil, fmt.Errorf("after %d attempts, last error: %w", maxRetires, err)
}

func GetQueryParams(api weather.UrlParamsInterface, cityInfo *weather.CityInfo, additional *map[string]string) *map[string]string {
	queryParams := api.GetUrlParams(cityInfo)
	if additional != nil {
		for key, val := range *additional {
			(*queryParams)[key] = val
		}
	}
	return queryParams
}
