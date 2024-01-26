# connectdial
## Introduction

Dialer implementation to enforce tunnelling of HTTP requests.

This package provides a workaround for the limitation of Go's standard library regarding HTTP request tunneling. 

In the default Go proxy implementation, tunneling is limited to HTTPS requests, which meets most of the cases. However, there are situations where you need to use tunneling for unencrypted requests as well. This package allows to bypass the limitation and enforce tunnelling for all HTTP requests. 

Compatible with `golang.org/x/net/proxy`

## Usage
```go
config := connectdial.Config{
    ConnectionString: "127.0.0.1:8888",
}

dialer, err := connectdial.New(config)
if err != nil { /*...*/ }

transport := &http.Transport{DialContext: dialer.DialContext}
client := http.Client{Transport: transport}

resp, err := client.Get("https://example.com")
if err != nil { /*...*/ }

log.Println(resp.StatusCode) // 200
```
