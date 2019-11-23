package eet

import (
	"bytes"
	"crypto/x509"
	"encoding/xml"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Service string

const (
	PlaygroundService Service = "https://pg.eet.cz:443/eet/services/EETServiceSOAP/v3"
	ProductionService Service = "https://prod.eet.cz:443/eet/services/EETServiceSOAP/v3"
)

type Regime int

const (
	RegularRegime Regime = iota
	SimplifiedRegime
)

type Dispatcher struct {
	service     Service
	signer      *Signer
	certificate *x509.Certificate
	testing     bool
}

func NewDispatcher(service Service, certPath, password string) (*Dispatcher, error) {
	signer, err := NewSigner(certPath, password)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create signer")
	}

	d := Dispatcher{
		service: service,
		signer:  signer,
	}

	return &d, nil
}

func (d *Dispatcher) SendPayment(receipt Receipt) (*Response, error) {
	trzba, err := receipt.Trzba(d.signer)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert Receipt to Trzba")
	}

	envelope, err := NewSOAPEnvelopeRequest(trzba, d.signer)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create SOAPEnvelopeRequest")
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	if err := xml.NewEncoder(&buf).Encode(envelope); err != nil {
		return nil, errors.Wrap(err, "Failed to marshal SOAPEnvelopeRequest")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(string(d.service), "application/xml", &buf)
	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send payment")
	}

	var resEnvelope SOAPEnvelopeResponse
	if err := xml.NewDecoder(resp.Body).Decode(&resEnvelope); err != nil {
		return nil, errors.Wrap(err, "Failed to xml.Unmarshal SOAPEnvelopeResponse")
	}

	odpoved := resEnvelope.Body.Odpoved
	if odpoved.Chyba != nil {
		return nil, odpoved.Chyba
	}

	response := Response{
		DatPrij: odpoved.Hlavicka.DatPrij,
		Fik:     odpoved.Potvrzeni.Fik,
		Bkp:     odpoved.Hlavicka.Bkp,
		odpoved: odpoved,
	}

	return &response, nil
}
