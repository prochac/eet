package goEET

import (
	"time"

	"fmt"

	"testing"

	"github.com/satori/go.uuid"
)

func TestDispatcher_SendPayment(t *testing.T) {
	d, err := NewDispatcher(PlaygroundService, "cert_test/EET_CA1_Playground-CZ00000019.p12", "eet")
	if err != nil {
		t.Fatal(err)
		return
	}
	r := Receipt{
		UuidZpravy: uuid.Must(uuid.NewV4()).String(),
		DicPopl:    "CZ00000019",
		IdProvoz:   273,
		IdPokl:     "/5546/RO24",
		PoradCis:   "0/6460/ZQ42",
		DatTrzby:   time.Now(),
		CelkTrzba:  0,
		Rezim:      RegularRegime,
	}
	response, err := d.SendPayment(r)
	if err != nil {
		t.Fatal(err)
		return
	}
	if len(response.Fik) == 0 {
		t.Fatal("send receipt failed")

	}
	fmt.Println("Fik: ", response.Fik)
	fmt.Println("Bkp: ", response.Bkp)
}
