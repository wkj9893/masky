package dns

import (
	"context"
	"net"
)

func SetResolver(dns string) {
	if _, _, err := net.SplitHostPort(dns); err != nil {
		// Append the default DNS port
		dns = net.JoinHostPort(dns, "53")
	}
	dialer := net.Dialer{}
	net.DefaultResolver = &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, dns)
		},
	}
}
