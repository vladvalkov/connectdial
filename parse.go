package connectdial

import (
	"fmt"
	"net/url"
	"strings"
)

func parseProxy(proxy string) (*url.URL, error) {
	if proxy == "" {
		return nil, nil
	}
	proxyURL, err := url.Parse(proxy)
	if err != nil || (proxyURL.Scheme != "http" && proxyURL.Scheme != "https") {
		proxyURL, err = url.Parse("http://" + proxy)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %q: %v", ErrInvalidConnectionString, proxy, err)
	}

	if !strings.ContainsRune(proxyURL.Host, ':') {
		switch proxyURL.Scheme {
		case "http":
			proxyURL.Host = proxyURL.Host + ":80"
		case "https":
			proxyURL.Host = proxyURL.Host + ":443"
		default:
			return nil, fmt.Errorf("%w: %q: port is not specified", ErrInvalidConnectionString, proxy)
		}
	}

	return proxyURL, nil
}
