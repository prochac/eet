package signing

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"io/ioutil"

	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"

	"github.com/prochac/crypto/pkcs12"
)

type Signer struct {
	cert *x509.Certificate
	key  *rsa.PrivateKey
}

func (s *Signer) Base64Cert() string {
	rawCert := s.cert.Raw
	return base64.StdEncoding.EncodeToString(rawCert)
}

// Sign signs data with rsa-sha256
func (s *Signer) Sign(data []byte) ([]byte, error) {
	hashed := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, s.key, crypto.SHA256, hashed[:])
}

func (s *Signer) LoadCertificate(certPath string) error {
	pemBytes, err := ioutil.ReadFile(certPath)
	if err != nil {
		return err
	}
	b, _ := pem.Decode(pemBytes)
	if b == nil || b.Type != "CERTIFICATE" {
		return errors.New("ssl: load certificate failed")
	}
	certificate, err := x509.ParseCertificate(b.Bytes)
	if err != nil {
		return err
	}
	s.cert = certificate
	return nil
}

func (s *Signer) LoadPrivateKey(keyPath string) (err error) {
	pemBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return errors.New("ssl: no key found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		s.key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("ssl: unsupported key type %q", block.Type)
	}
	return nil
}

func NewSigner(certPath string, password string) (s *Signer, err error) {
	pfxData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	rawKeys, certs, err := pkcs12.DecodeAll(pfxData, password)
	if err != nil {
		return nil, err
	}

	signer := Signer{}

	for _, rawKey := range rawKeys {
		if key, ok := rawKey.(*rsa.PrivateKey); ok {
			signer.key = key
			break
		}
	}
	if signer.key == nil {
		return nil, errors.New("can't load private key")
	}

	for _, cert := range certs {
		if !cert.IsCA {
			signer.cert = cert
			break
		}
	}
	if signer.cert == nil {
		return nil, errors.New("can't load certificate")
	}

	return &signer, nil
}
