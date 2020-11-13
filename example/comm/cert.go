package comm

import (
	"crypto/tls"
	"crypto/x509"
)

// GetTLSConfig returns a tls config for quic.clemente.io
func GetTLSConfig(certFile string, keyFile string) *tls.Config {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	pool := x509.NewCertPool()
	if tmp, err := x509.ParseCertificate(cert.Certificate[0]); err != nil {
		panic(err)
	} else {
		pool.AddCert(tmp)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
}

// GetCertificate returns a certificate for quic.clemente.io
func GetCertificate(certFile string, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	return cert
}
