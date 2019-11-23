package eet

import (
	"time"

	"github.com/pkg/errors"
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

func (r Receipt) Trzba(signer *Signer) (Trzba, error) {
	var t Trzba
	var err error
	// Hlavicka
	t.Hlavicka.DatOdesl = NewDateTimeType(time.Now())
	t.Hlavicka.UuidZpravy, err = NewUUIDType(r.UuidZpravy)
	if err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create UuidZpravy")
	}
	t.Hlavicka.PrvniZaslani = r.PrvniZaslani
	// Data
	if t.Data.DicPopl, err = NewCZDICType(r.DicPopl); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create DicPopl")
	}
	if t.Data.DicPoverujiciho, err = NewCZDICType(r.DicPoverujiciho); len(r.DicPoverujiciho) != 0 && err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create DicPoverujiciho")
	}
	if t.Data.IdProvoz, err = NewIdProvozType(r.IdProvoz); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create IdProvoz")
	}
	if t.Data.IdPokl, err = NewString20(r.IdPokl); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create IdPokl")
	}
	if t.Data.PoradCis, err = NewString25(r.PoradCis); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create PoradCis")
	}
	t.Data.DatTrzby = NewDateTimeType(r.DatTrzby)
	if t.Data.CelkTrzba, err = NewCastkaType(r.CelkTrzba); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create CelkTrzba")
	}
	if t.Data.ZaklNepodlDph, err = NewCastkaType(r.ZaklNepodlDph); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create ZaklNepodlDph")
	}
	if t.Data.ZaklDan1, err = NewCastkaType(r.ZaklDan1); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create ZaklDan1")
	}
	if t.Data.Dan1, err = NewCastkaType(r.Dan1); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create Dan1")
	}
	if t.Data.ZaklDan2, err = NewCastkaType(r.ZaklDan2); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create ZaklDan2")
	}
	if t.Data.Dan2, err = NewCastkaType(r.Dan2); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create Dan2")
	}
	if t.Data.ZaklDan3, err = NewCastkaType(r.ZaklDan3); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create ZaklDan3")
	}
	if t.Data.Dan3, err = NewCastkaType(r.Dan3); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create Dan3")
	}
	if t.Data.CestSluz, err = NewCastkaType(r.CestSluz); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create CestSluz")
	}
	if t.Data.PouzitZboz1, err = NewCastkaType(r.PouzitZboz1); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create PouzitZboz1")
	}
	if t.Data.PouzitZboz2, err = NewCastkaType(r.PouzitZboz2); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create PouzitZboz2")
	}
	if t.Data.PouzitZboz3, err = NewCastkaType(r.PouzitZboz3); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create PouzitZboz3")
	}
	if t.Data.UrcenoCerpZuct, err = NewCastkaType(r.UrcenoCerpZuct); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create UrcenoCerpZuct")
	}
	if t.Data.CerpZuct, err = NewCastkaType(r.CerpZuct); err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create CerpZuct")
	}
	if r.Rezim == SimplifiedRegime {
		t.Data.Rezim = ZjednodusenyRezim
	}
	// KontrolniKody
	pkp, err := NewPkp(t, signer)
	if err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create PKP")
	}
	t.KontrolniKody.Pkp = pkp

	bkp, err := NewBkp(t.KontrolniKody.Pkp)
	if err != nil {
		return Trzba{}, errors.Wrap(err, "Failed to create BKP")
	}
	t.KontrolniKody.Bkp = bkp

	return t, nil
}
