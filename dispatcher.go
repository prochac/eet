package goEET

import (
	"bytes"
	"crypto/x509"
	"net/http"

	"github.com/pkg/errors"
)

type Service string

const (
	PlaygroundService Service = "https://pg.eet.cz:443/eet/services/EETServiceSOAP/v3"
	ProductionService Service = "https://prod.eet.cz:443/eet/services/EETServiceSOAP/v3"
)

func (s Service) ToString() string {
	return string(s)
}

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

func NewDispatcher(service Service, certPath, password string) (eet *Dispatcher, err error) {
	signer, err := NewSigner(certPath, password)
	if err != nil {
		return &Dispatcher{}, errors.Wrap(err, "Failed to create signer")
	}

	d := Dispatcher{
		service: service,
		signer:  signer,
	}

	return &d, nil
}

func (d *Dispatcher) SendPayment(receipt Receipt) (r *Response, err error) {
	trzba, err := receipt.Trzba(d.signer)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert Receipt to Trzba")
	}

	envelope, err := NewSOAPEnvelopeRequest(trzba, d.signer)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create SOAPEnvelopeRequest")
	}

	b, err := envelope.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal SOAPEnvelopeRequest")
	}

	buf := bytes.NewBuffer(b)
	resp, err := http.Post(d.service.ToString(), "application/xml", buf)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send payment")
	}
	defer resp.Body.Close()

	resEnvelope, err := ParseSOAPEnvelopeResponse(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed parse response")
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
