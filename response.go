package eet

import (
	"time"
)

type Response struct {
	odpoved Odpoved
	DatPrij time.Time
	Fik     string
	Bkp     string
}

func (r Response) Warnings() []string {
	warnings := make([]string, len(r.odpoved.Varovani))
	for i, v := range r.odpoved.Varovani {
		warnings[i] = v.Varovani
	}
	return warnings
}
