package connectdial

import (
	"net/url"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tcs := []struct {
		name    string
		conn    string
		user    *url.Userinfo
		address string
		err     bool
	}{
		{name: "IP address with port",
			conn:    "127.0.0.1:8888",
			address: "127.0.0.1:8888"},
		{name: "IP address with user",
			conn:    "user@127.0.0.1:8888",
			user:    url.User("user"),
			address: "127.0.0.1:8888"},
		{name: "IP address with user and password",
			conn:    "user:pass@127.0.0.1:8888",
			user:    url.UserPassword("user", "pass"),
			address: "127.0.0.1:8888"},
		{name: "IP address with protocol and user",
			conn:    "http://user:pass@127.0.0.1:8888",
			user:    url.UserPassword("user", "pass"),
			address: "127.0.0.1:8888"},
		{name: "IP address with protocol and user",
			conn:    "https://user:pass@127.0.0.1:8888",
			user:    url.UserPassword("user", "pass"),
			address: "127.0.0.1:8888"},
		{name: "IP address without port",
			conn:    "http://127.0.0.1",
			address: "127.0.0.1:80"},
		{name: "IP address without port",
			conn:    "127.0.0.1",
			address: "127.0.0.1:80"},
		{name: "IP address without port",
			conn:    "https://127.0.0.1",
			address: "127.0.0.1:443"},
		{name: "domain name without a protocol",
			conn:    "domain.com",
			address: "domain.com:80"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			u, err := parseProxy(tc.conn)
			if g, e := err != nil, tc.err; g != e {
				if g && !e {
					t.Errorf("got unexpected error while parsing: %v", err)
				}
				if !g && e {
					t.Errorf("expected to get an error")
				}
				t.FailNow()
			}
			if g, e := u.User, tc.user; !reflect.DeepEqual(g, e) {
				t.Errorf(`got user "%v" expected "%v"`, g, e)
			}
			if g, e := u.Host, tc.address; g != e {
				t.Errorf(`got user "%v" expected "%v"`, g, e)
			}
		})
	}
}
