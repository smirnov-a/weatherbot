package utils

import (
	"context"
	"fmt"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"time"
	"weatherbot/config"
	"weatherbot/internal/logger"
)

type Auth struct {
	User     string
	Password string
}

type HTTPProxyDialer struct {
	ProxyURL *url.URL
}

const timeOut = 30

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

func GetHttpClient() *http.Client {
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
		client.Timeout = timeOut * time.Second
	}

	return client
}

func getPassword(proxyURI *url.URL) string {
	password, _ := proxyURI.User.Password()
	return password
}
