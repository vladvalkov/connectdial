package connectdial

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/elazarl/goproxy"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDialer(t *testing.T) {

	useStdLib := func(proxyUrl *url.URL, proxyCert, handlerCert *x509.Certificate) *http.Transport {
		return &http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			TLSClientConfig: makeTLSForCertificates(proxyCert, handlerCert),
		}
	}
	useCustomDialer := func(t *testing.T, proxyUrl *url.URL, proxyCert, handlerCert *x509.Certificate) *http.Transport {
		c := Config{
			ConnectionString: proxyUrl.String(),
		}
		if proxyCert != nil {
			c.TLS = makeTLSForCertificates(proxyCert)
		}
		dialer, err := New(c)
		if err != nil {
			t.Errorf("could not create a dialer: %v", err)
		}
		return &http.Transport{
			DialContext:     dialer.DialContext,
			TLSClientConfig: makeTLSForCertificates(handlerCert),
		}
	}

	tcs := []struct {
		name                 string
		proxyTLS, handlerTLS bool
		stdLib               bool
	}{
		{name: "stdlib no auth or tls",
			proxyTLS: false, handlerTLS: false,
			stdLib: true},
		{name: "stdlib no auth with handler tls",
			proxyTLS: false, handlerTLS: true,
			stdLib: true},
		{name: "stdlib no auth with proxy tls",
			proxyTLS: true, handlerTLS: false,
			stdLib: true},
		{name: "stdlib no auth with proxy and handler tls",
			proxyTLS: true, handlerTLS: true,
			stdLib: true},
		{name: "connectdial no auth or tls",
			proxyTLS: false, handlerTLS: false},
		{name: "connectdial no auth with handler tls",
			proxyTLS: false, handlerTLS: true},
		{name: "connectdial no auth with proxy tls",
			proxyTLS: true, handlerTLS: false},
		{name: "connectdial no auth with proxy and handler tls",
			proxyTLS: true, handlerTLS: true},
		/*
			TODO: Add tests for proxies with authentication
		*/
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			helloWorldUrl, handlerCert := createHelloWorldServer(t, tc.handlerTLS)
			proxyUrl, proxyCert := createProxyServer(t, tc.proxyTLS)

			var transport *http.Transport
			if tc.stdLib {
				transport = useStdLib(proxyUrl, proxyCert, handlerCert)
			} else {
				transport = useCustomDialer(t, proxyUrl, proxyCert, handlerCert)
			}
			client := http.Client{
				Transport: transport,
			}
			if err := expectHelloWorld(client.Get(helloWorldUrl.String())); err != nil {
				t.Errorf("%v", err)
			}
		})
	}
}

func expectHelloWorld(resp *http.Response, err error) error {
	if err != nil {
		return fmt.Errorf("got error: %w", err)
	}
	if resp == nil && err == nil {
		return fmt.Errorf("both response and error are nil")
	}
	if g, e := resp.StatusCode, 200; g != e {
		return fmt.Errorf("got non-200 status code")
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body")
	}
	if g, e := string(b), "Hello World!"; g != e {
		return fmt.Errorf("response body is different from expected")
	}
	return nil
}

func createProxyServer(t *testing.T, useTLS bool) (*url.URL, *x509.Certificate) {
	proxy := goproxy.NewProxyHttpServer()

	srv := httptest.NewUnstartedServer(proxy)
	if useTLS {
		srv.StartTLS()
	} else {
		srv.Start()
	}
	t.Cleanup(srv.Close)

	u, _ := url.Parse(srv.URL)
	return u, srv.Certificate()
}
func createHelloWorldServer(t *testing.T, useTLS bool) (*url.URL, *x509.Certificate) {
	helloWorld := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Hello World!"))
	})
	srv := httptest.NewUnstartedServer(helloWorld)
	if useTLS {
		srv.StartTLS()
	} else {
		srv.Start()
	}
	t.Cleanup(srv.Close)

	u, _ := url.Parse(srv.URL)
	return u, srv.Certificate()
}

func makeTLSForCertificates(certs ...*x509.Certificate) *tls.Config {
	var certPool = x509.NewCertPool()
	for _, v := range certs {
		if v != nil {
			certPool.AddCert(v)
		}
	}

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}
	return tlsConfig
}
