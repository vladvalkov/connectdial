package connectdial

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrUnsupportedDialNetwork  = errors.New("unsupported dialing protocol was passed")
	ErrAuthenticationFailed    = errors.New("authentication failed for proxy")
	ErrConnectionFailed        = errors.New("could not connect to proxy")
	ErrInvalidConnectionString = errors.New("invalid proxy connection string")
)

func New(c Config) (*Dialer, error) {
	u, err := parseProxy(c.ConnectionString)
	if err != nil {
		return nil, err
	}

	var auth ProxyAuth
	if c.Auth != nil {
		auth = c.Auth
	} else {
		auth = BasicUserinfo(u.User)
	}

	return &Dialer{
		proxyHost: u.Host,
		tls:       c.TLS,
		auth:      auth,
		dialer:    &net.Dialer{},
	}, nil
}

type Config struct {
	ConnectionString string
	Auth             ProxyAuth
	TLS              *tls.Config
}

type Dialer struct {
	proxyHost string
	auth      ProxyAuth
	tls       *tls.Config
	dialer    *net.Dialer
}

func (d *Dialer) Dial(network string, addr string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, addr)
}

func (d *Dialer) DialContext(ctx context.Context, network string, addr string) (net.Conn, error) {
	if network != "tcp" {
		return nil, fmt.Errorf("connectdial: %w: only tcp is supported", ErrUnsupportedDialNetwork)
	}
	proxyConn, err := d.dialProxy(ctx)
	if err != nil {
		return nil, fmt.Errorf("connectdial: could not dial proxy: %w", err)
	}
	err = d.authProxy(proxyConn, addr)
	if err != nil {
		return nil, fmt.Errorf("connectdial: %w", err)
	}
	return proxyConn, nil
}

func (d *Dialer) authProxy(conn net.Conn, addr string) error {
	req := &http.Request{
		Method: "CONNECT",
		URL:    &url.URL{Opaque: addr},
		Host:   addr,
		Header: make(http.Header),
	}
	if d.auth != nil {
		req.Header.Set("Proxy-Authorization", authHeader(d.auth))
	}
	response, err := doRoundTrip(conn, req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrConnectionFailed, err)
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("%w: got status code %d", ErrAuthenticationFailed, response.StatusCode)
	}
	return nil
}

func (d *Dialer) dialProxy(ctx context.Context) (net.Conn, error) {
	if d.tls != nil {
		return tls.DialWithDialer(d.dialer, "tcp", d.proxyHost, d.tls)
	}
	return d.dialer.DialContext(ctx, "tcp", d.proxyHost)
}

func doRoundTrip(conn net.Conn, r *http.Request) (*http.Response, error) {
	err := r.Write(conn)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(conn)
	return http.ReadResponse(reader, r)
}

func authHeader(auth ProxyAuth) string {
	var str strings.Builder
	t, cr := auth.Type(), auth.Credentials()
	str.Grow(len(t) + 1 + len(cr))
	str.WriteString(t)
	str.WriteRune(' ')
	str.WriteString(cr)
	return str.String()
}
