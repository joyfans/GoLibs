package socket

import (
    "fmt"
    "net"
    "time"
    "golang.org/x/net/proxy"
)

type Socks5Dialer struct {
    socks5 proxy.Dialer
    timeout time.Duration
}

func newSocks5Dialer(network string, proxyAddress string, proxyPort int, auth *proxy.Auth) *Socks5Dialer {
    s := &Socks5Dialer{
        timeout: 0,
    }

    var err error

    s.socks5, err = proxy.SOCKS5(network, fmt.Sprintf("%s:%d", proxyAddress, proxyPort), auth, s)
    RaiseSocketError(err)

    return s
}

func (self *Socks5Dialer) setTimeout(timeout time.Duration) {
    if timeout >= 0 {
        self.timeout = timeout
    }
}

func (self *Socks5Dialer) Connect(network, host string, port int, timeout time.Duration) (net.Conn, error) {
    self.setTimeout(timeout)
    return self.socks5.Dial(network, fmt.Sprintf("%s:%d", host, port))
}

func (self *Socks5Dialer) Dial(network, address string) (net.Conn, error) {
    if self.timeout <= 0 {
        return net.Dial(network, address)
    }

    return net.DialTimeout(network, address, self.timeout)
}
