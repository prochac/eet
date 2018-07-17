package goEET

import (
	"encoding/xml"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/prochac/goEET/odpoved"
)

type Response struct {
	odpoved odpoved.Odpoved
	DatPrij time.Time
	Fik     string
	Bkp     string
}

func (r Response) Warnings() (warnings []string) {
	for _, v := range r.odpoved.Varovani {
		logs.Info(v)
		warnings = append(warnings, v.Varovani)
	}
	return warnings
}

func parseResponseFromWsseResponse(rawResponse []byte) (*Response, error) {
	resEnvelope := odpoved.SOAPEnvelope{}

	if err := xml.Unmarshal(rawResponse, &resEnvelope); err != nil {
		panic(err)
	}
	o := resEnvelope.Body.Odpoved

	if o.Chyba != nil {
		return nil, o.Chyba
	}

	resp := Response{
		DatPrij: o.Hlavicka.DatPrij,
		Fik:     o.Potvrzeni.Fik,
		Bkp:     o.Hlavicka.Bkp,
	}
	return &resp, nil
}
