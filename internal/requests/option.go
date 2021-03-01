package requests

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type Requests struct {
	client          *http.Client
	Options         *Options
	CustomHeaders   map[string]string
}

// Options contains configuration options for the client
type Options struct {
	DefaultUserAgent string
	HTTPProxy        string
	SocksProxy       string
	Threads          int
	Timeout time.Duration
	CustomHeaders map[string]string
	FollowRedirects      bool
	FollowHostRedirects  bool
	Unsafe               bool
}

// DefaultOptions contains the default options
var DefaultOptions = Options{
	Threads:  25,
	Timeout:  30 * time.Second,
	Unsafe:   false,
	DefaultUserAgent:         "Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36",
}

func DefaultTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConnsPerHost: -1,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableKeepAlives: true,
	}
	return transport
}

