package goEET

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
)

type SOAPEnvelopeResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		XMLName xml.Name `xml:"Body"`
		Odpoved Odpoved
	}
}

func ParseSOAPEnvelopeResponse(r io.Reader) (resp SOAPEnvelopeResponse, err error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return SOAPEnvelopeResponse{}, errors.Wrap(err, "Failed to read all data from reader")
	}

	if err := xml.Unmarshal(b, &resp); err != nil {
		return SOAPEnvelopeResponse{}, errors.Wrap(err, "Failed to xml.Unmarshal SOAPEnvelopeResponse")
	}

	return resp, nil
}

type Odpoved struct {
	XMLName   xml.Name        `xml:"Odpoved"`
	Hlavicka  OdpovedHlavicka `xml:"Hlavicka"`
	Potvrzeni *Potvrzeni      `xml:"Potvrzeni"`
	Chyba     *Chyba          `xml:"Chyba"`
	Varovani  []Varovani      `xml:"Varovani"`
}

type OdpovedHlavicka struct {
	XMLName    xml.Name  `xml:"Hlavicka"`
	UuidZpravy string    `xml:"uuid_zpravy,attr"`
	DatPrij    time.Time `xml:"dat_prij,attr"`
	DatOdmit   time.Time `xml:"dat_odmit,attr"`
	Bkp        string    `xml:"bkp,attr"`
}

type Potvrzeni struct {
	XMLName xml.Name `xml:"Potvrzeni"`
	Fik     string   `xml:"fik,attr"`
	Test    bool     `xml:"test,attr"`
}

type Chyba struct {
	XMLName xml.Name `xml:"Chyba"`
	Kod     int      `xml:"kod,attr"`
	Test    bool     `xml:"test,attr"`
	Chyba   string   `xml:",chardata"`
}

func (ch Chyba) Error() string {
	return fmt.Sprintf("%d %s", ch.Kod, ch.Chyba)
}

type Varovani struct {
	XMLName  xml.Name `xml:"Varovani"`
	KodVarov int      `xml:"kod_varov,attr"`
	Varovani string   `xml:",chardata"`
}
