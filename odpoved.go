package eet

import (
	"encoding/xml"
	"fmt"
	"time"
)

type SOAPEnvelopeResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		XMLName xml.Name `xml:"Body"`
		Odpoved Odpoved
	}
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
