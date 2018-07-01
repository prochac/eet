package goEET

import (
	"time"

	"github.com/prochac/goEET/signing"
	"github.com/prochac/goEET/trzba"
)

type Receipt struct {
	UuidZpravy      string
	PrvniZaslani    bool
	DicPopl         string
	DicPoverujiciho string
	IdProvoz        int
	IdPokl          string
	PoradCis        string
	DatTrzby        time.Time
	CelkTrzba       float64
	ZaklNepodlDph   float64
	ZaklDan1        float64
	Dan1            float64
	ZaklDan2        float64
	Dan2            float64
	ZaklDan3        float64
	Dan3            float64
	CestSluz        float64
	PouzitZboz1     float64
	PouzitZboz2     float64
	PouzitZboz3     float64
	UrcenoCerpZuct  float64
	CerpZuct        float64
	Rezim           Regime
}

func (r Receipt) toTrzba(signer *signing.Signer) (t trzba.Trzba, err error) {
	// Hlavicka
	t.Hlavicka.DatOdesl = trzba.NewDateTimeType(time.Now())
	t.Hlavicka.UuidZpravy, err = trzba.NewUUIDType(r.UuidZpravy)
	if err != nil {
		return
	}
	t.Hlavicka.PrvniZaslani = r.PrvniZaslani
	// Data
	if t.Data.DicPopl, err = trzba.NewCZDICType(r.DicPopl); err != nil {
		return
	}
	if t.Data.DicPoverujiciho, err = trzba.NewCZDICType(r.DicPoverujiciho); len(r.DicPoverujiciho) != 0 && err != nil {
		return
	}
	if t.Data.IdProvoz, err = trzba.NewIdProvozType(r.IdProvoz); err != nil {
		return
	}
	if t.Data.IdPokl, err = trzba.NewString20(r.IdPokl); err != nil {
		return
	}
	if t.Data.PoradCis, err = trzba.NewString25(r.PoradCis); err != nil {
		return
	}
	t.Data.DatTrzby = trzba.NewDateTimeType(r.DatTrzby)
	if t.Data.CelkTrzba, err = trzba.NewCastkaType(r.CelkTrzba); err != nil {
		return
	}
	if t.Data.ZaklNepodlDph, err = trzba.NewCastkaType(r.ZaklNepodlDph); err != nil {
		return
	}
	if t.Data.ZaklDan1, err = trzba.NewCastkaType(r.ZaklDan1); err != nil {
		return
	}
	if t.Data.Dan1, err = trzba.NewCastkaType(r.Dan1); err != nil {
		return
	}
	if t.Data.ZaklDan2, err = trzba.NewCastkaType(r.ZaklDan2); err != nil {
		return
	}
	if t.Data.Dan2, err = trzba.NewCastkaType(r.Dan2); err != nil {
		return
	}
	if t.Data.ZaklDan3, err = trzba.NewCastkaType(r.ZaklDan3); err != nil {
		return
	}
	if t.Data.Dan3, err = trzba.NewCastkaType(r.Dan3); err != nil {
		return
	}
	if t.Data.CestSluz, err = trzba.NewCastkaType(r.CestSluz); err != nil {
		return
	}
	if t.Data.PouzitZboz1, err = trzba.NewCastkaType(r.PouzitZboz1); err != nil {
		return
	}
	if t.Data.PouzitZboz2, err = trzba.NewCastkaType(r.PouzitZboz2); err != nil {
		return
	}
	if t.Data.PouzitZboz3, err = trzba.NewCastkaType(r.PouzitZboz3); err != nil {
		return
	}
	if t.Data.UrcenoCerpZuct, err = trzba.NewCastkaType(r.UrcenoCerpZuct); err != nil {
		return
	}
	if t.Data.CerpZuct, err = trzba.NewCastkaType(r.CerpZuct); err != nil {
		return
	}
	if r.Rezim == RegularRegime {
		t.Data.Rezim = trzba.BeznyRezim
	} else if r.Rezim == SimplifiedRegime {
		t.Data.Rezim = trzba.ZjednodusenyRezim
	}
	// KontrolniKody
	t.KontrolniKody.Pkp = trzba.NewPkp(t, signer)
	t.KontrolniKody.Bkp = trzba.NewBkp(t.KontrolniKody.Pkp)

	return t, nil
}
