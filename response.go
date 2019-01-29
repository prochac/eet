package eet

import (
	"log"
	"time"
)

type Response struct {
	odpoved Odpoved
	DatPrij time.Time
	Fik     string
	Bkp     string
}

func (r Response) Warnings() (warnings []string) {
	for _, v := range r.odpoved.Varovani {
		log.Println(v)
		warnings = append(warnings, v.Varovani)
	}
	return warnings
}
