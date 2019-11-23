package eet

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/pkcs12"
)

type Signer struct {
	cert *x509.Certificate
	key  *rsa.PrivateKey
}

func NewSigner(certPath string, password string) (*Signer, error) {
	pfxData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", certPath, err)
	}

	privateKey, certificate, _, err := decodeAll(pfxData, password)
	if err != nil {
		return nil, fmt.Errorf("decoding private key and certificates: %w", err)
	}

	s := Signer{
		key:  privateKey,
		cert: certificate,
	}

	return &s, nil
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

// decodeAll extracts all certificate and private keys from pfxData.
func decodeAll(pfxData []byte, password string) (*rsa.PrivateKey, *x509.Certificate, *x509.Certificate, error) {
	var (
		privateKey  *rsa.PrivateKey
		certificate *x509.Certificate
		caCert      *x509.Certificate
	)
	blocks, err := pkcs12.ToPEM(pfxData, password)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("converting binary pfx data to PEM: %w", err)
	}
	for _, block := range blocks {
		switch block.Type {
		case "PRIVATE KEY":
			if privateKey != nil {
				return nil, nil, nil, errors.New("only one private key expected")
			}

			privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("parsing private key: %w", err)
			}
		case "CERTIFICATE":
			certs, err := x509.ParseCertificates(block.Bytes)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("parsing cetrtificates: %w", err)
			}
			for _, cert := range certs {
				if cert.IsCA {
					if caCert != nil {
						return nil, nil, nil, errors.New("only one CA certificate expected")
					}
					caCert = cert
					continue
				}
				if certificate != nil {
					return nil, nil, nil, errors.New("only one certificate expected")
				}
				certificate = cert
			}
		}
	}

	return privateKey, certificate, caCert, nil
}
