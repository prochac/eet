package goEET

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"crypto/x509"

	"github.com/prochac/goEET/signing"
	"github.com/prochac/goEET/wsse"
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
	signer      *signing.Signer
	certificate *x509.Certificate
	testing     bool
}

func NewDispatcher(service Service, certPath, password string) (eet *Dispatcher, err error) {
	d := Dispatcher{service: service}
	d.signer, err = signing.NewSigner(certPath, password)

	return &d, err
}

func (d *Dispatcher) SendPayment(receipt Receipt) (r *Response, err error) {
	trzba, err := receipt.toTrzba(d.signer)
	if err != nil {
		return nil, err
	}

	envelope := wsse.NewSOAPEnvelope(trzba, d.signer)

	b, err := xml.Marshal(envelope)
	if err != nil {
		panic(err)
	}
	xmlHeader := `<?xml version="1.0" encoding="UTF-8"?>`
	buf := bytes.NewBuffer([]byte(xmlHeader + string(b)))

	resp, err := http.Post(d.service.ToString(), "application/xml", buf)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return parseResponseFromWsseResponse(respBody)
}
