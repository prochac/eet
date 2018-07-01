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

	"golang.org/x/crypto/pkcs12"
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

func NewSigner(pkcsPath string, password string) (s *Signer, err error) {
	pemBytes, err := ioutil.ReadFile(pkcsPath)
	if err != nil {
		return nil, err
	}

	rawKey, cert, err := pkcs12.Decode(pemBytes, password)
	if err != nil {
		return nil, err
	}
	key, ok := rawKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("can't load private key")
	}

	signer := Signer{
		cert: cert,
		key:  key,
	}

	return &signer, nil
}
