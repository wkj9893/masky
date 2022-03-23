package tls

import (
	"crypto/rand"
	"crypto/rsa"
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

// Setup a bare-bones TLS config for the server
func GenerateTLSConfig() (*tls.Config, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{cert},
			PrivateKey:  key,
		}},
		NextProtos: []string{defaultALPN},
	}, nil
}
