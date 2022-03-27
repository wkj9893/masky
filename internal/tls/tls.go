package tls

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"math/big"
)

const defaultALPN = "masky"

var ClientTLSConfig = &tls.Config{
	InsecureSkipVerify: true,
	NextProtos:         []string{defaultALPN},
	ClientSessionCache: tls.NewLRUClientSessionCache(1),
	CipherSuites:       []uint16{tls.TLS_CHACHA20_POLY1305_SHA256},
}

func GenerateTLSConfig() (*tls.Config, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, pub, priv)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{cert},
			PrivateKey:  priv,
		}},
		NextProtos: []string{defaultALPN},
	}, nil
}
