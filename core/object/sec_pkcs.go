package object

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
	p12 "software.sslmate.com/src/go-pkcs12"
)

// PKCS returns the PKCS#12 format bytes of the private key and certificate
// chain stored in this secure keystore
func (t *sec) PKCS() ([]byte, error) {
	if !t.HasKey("private_key") {
		return nil, errors.Errorf("private_key does not exist")
	}
	privateKeyBytes, err := t.DecodeKey("private_key")
	if err != nil {
		return nil, err
	}
	block, rest := pem.Decode(privateKeyBytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.Errorf("failed to decode PEM block containing private key")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	if !t.HasKey("certificate_chain") {
		return nil, errors.Errorf("certificate_chain does not exist")
	}
	certificateChainBytes, err := t.DecodeKey("certificate_chain")
	if err != nil {
		return nil, err
	}
	l := make([]*x509.Certificate, 0)
	for {
		block, certificateChainBytes = pem.Decode(certificateChainBytes)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" {
			return nil, errors.Errorf("failed to decode PEM block containing certificate. actual type %s", block.Type)
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		l = append(l, cert)
		if rest == nil {
			break
		}
	}
	if len(l) < 1 {
		return nil, errors.Errorf("certificate_chain has no valid certificate")
	}
	return p12.Encode(rand.Reader, privateKey, l[0], l[1:], "foo")
}
