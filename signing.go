package goEET

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"io/ioutil"

	"github.com/pkg/errors"
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

//func (s *Signer) LoadCertificate(certPath string) error {
//	pemBytes, err := ioutil.ReadFile(certPath)
//	if err != nil {
//		return errors.Wrapf(err, "Failed to read file %s", certPath)
//	}
//
//	block, _ := pem.Decode(pemBytes)
//	if block == nil {
//		return errors.New("No certificate found: decoded block is empty")
//	}
//
//	var certificate *x509.Certificate
//	switch block.Type {
//	case "CERTIFICATE":
//		certificate, err = x509.ParseCertificate(block.Bytes)
//		if err != nil {
//			return errors.Wrap(err, "Failed to parse certificate")
//		}
//	default:
//		return fmt.Errorf("ssl: unsupported key type %q", block.Type)
//	}
//	s.cert = certificate
//
//	return nil
//}
//
//func (s *Signer) LoadPrivateKey(keyPath string) (err error) {
//	pemBytes, err := ioutil.ReadFile(keyPath)
//	if err != nil {
//		return errors.Wrapf(err, "Failed to read file %s", keyPath)
//	}
//
//	block, _ := pem.Decode(pemBytes)
//	if block == nil {
//		return errors.New("No private key found: decoded block is empty")
//	}
//
//	var privateKey *rsa.PrivateKey
//	switch block.Type {
//	case "RSA PRIVATE KEY":
//		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
//		if err != nil {
//			return errors.Wrap(err, "Failed to parse private key")
//		}
//	default:
//		return fmt.Errorf("ssl: unsupported key type %q", block.Type)
//	}
//	s.key = privateKey
//
//	return nil
//}

func NewSigner(certPath string, password string) (*Signer, error) {
	pfxData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read file %s", certPath)
	}

	privateKeys, certificates, err := pkcs12.DecodeAll(pfxData, password)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode private keys and certificates from pfxData")
	}

	var privateKey *rsa.PrivateKey
	for _, rawKey := range privateKeys {
		if key, ok := rawKey.(*rsa.PrivateKey); ok {
			privateKey = key
			break
		}
	}
	if privateKey == nil {
		return nil, errors.New("can't load private key")
	}

	var certificate *x509.Certificate
	for _, cert := range certificates {
		if !cert.IsCA {
			certificate = cert
			break
		}
	}
	if certificate == nil {
		return nil, errors.New("can't load certificate")
	}

	s := Signer{
		key:  privateKey,
		cert: certificate,
	}

	return &s, nil
}
