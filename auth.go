package connectdial

import (
	"encoding/base64"
	"net/url"
	"strings"
)

/*
ProxyAuth is the interface used for user authorisation.
We avoid String() method to distinguish the interface from common Stringer implementations.
*/
type ProxyAuth interface {
	Type() string
	Credentials() string
}

func BasicUserinfo(userinfo *url.Userinfo) ProxyAuth {
	if userinfo == nil {
		return nil
	}
	user := userinfo.Username()
	pass, set := userinfo.Password()
	return basicAuth{
		user:        user,
		password:    pass,
		passwordSet: set,
	}
}

func BasicUserPassword(user, password string) ProxyAuth {
	return basicAuth{
		user:        user,
		password:    password,
		passwordSet: true,
	}
}

type basicAuth struct {
	user        string
	password    string
	passwordSet bool
}

func (a basicAuth) User() string {
	return a.user
}
func (a basicAuth) Password() (string, bool) {
	return a.password, a.passwordSet
}

func (basicAuth) Type() string {
	return "Basic"
}
func (a basicAuth) Credentials() string {
	var credStr strings.Builder
	credStr.WriteString(a.User())
	if pass, exists := a.Password(); exists {
		credStr.WriteRune(':')
		credStr.WriteString(pass)
	}
	base64Cred := base64.StdEncoding.EncodeToString([]byte(credStr.String()))
	return base64Cred
}
