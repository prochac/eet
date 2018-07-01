package trzba

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"regexp"
	"time"

	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/prochac/goEET/signing"
)

type String20 string

func NewString20(string20 string) (String20, error) {
	regex := `[0-9a-zA-Z\.,:;/#\-_ ]{1,20}`
	if ok, _ := regexp.MatchString(regex, string20); ok {
		return String20(string20), nil
	}
	return "", errors.New("invalid String20")
}

type String25 string

func NewString25(string25 string) (String25, error) {
	regex := `[0-9a-zA-Z\.,:;/#\-_ ]{1,25}`
	if ok, _ := regexp.MatchString(regex, string25); ok {
		return String25(string25), nil
	}
	return "", errors.New("invalid String25")
}

type DateTimeType string

func NewDateTimeType(dateTime time.Time) DateTimeType {
	return DateTimeType(dateTime.Format(time.RFC3339))
}

type CastkaType string

func NewCastkaType(castka float64) (CastkaType, error) {
	strCastka := fmt.Sprintf("%0.2f", castka)

	regex := `((0|-?[1-9]\d{0,7})\.\d\d|-0\.(0[1-9]|[1-9]\d))`
	if ok, _ := regexp.MatchString(regex, strCastka); ok {
		return CastkaType(strCastka), nil
	}
	return "0.00", errors.New("invalid CastkaType")
}

type IdProvozType int32

func NewIdProvozType(idProvozType int) (IdProvozType, error) {
	if 0 < idProvozType && idProvozType < 1000000 {
		return IdProvozType(idProvozType), nil
	}
	return 0, errors.New("invalid IdProvozType")
}

func (i IdProvozType) String() string {
	return fmt.Sprintf("%d", i)
}

type RezimType int

const (
	BeznyRezim RezimType = iota
	ZjednodusenyRezim
)

type UUIDType string

func NewUUIDType(uuid string) (UUIDType, error) {
	regex := `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}`
	if ok, _ := regexp.MatchString(regex, uuid); ok {
		return UUIDType(uuid), nil
	}
	return "", errors.New("invalid UUIDType")

}

type CZDICType string

func NewCZDICType(dic string) (CZDICType, error) {
	regex := `CZ[0-9]{8,10}`
	if ok, _ := regexp.MatchString(regex, dic); ok {
		return CZDICType(dic), nil
	}
	return "", errors.New("invalid CZDICType")

}

type PkpDigestType string

const (
	PkpDigestTypeSHA256 PkpDigestType = "SHA256"
)

type PkpCipherType string

const (
	PkpCipherTypeRSA2048 PkpCipherType = "RSA2048"
)

type PkpEncodingType string

const (
	PkpEncodingTypeBase64 PkpEncodingType = "base64"
)

type BkpDigestType string

const (
	BkpDigestTypeSHA1 BkpDigestType = "SHA1"
)

type BkpEncodingType string

const (
	BkpEncodingTypeBase16 BkpEncodingType = "base16"
)

type Trzba struct {
	XMLName       xml.Name      `xml:"http://fs.mfcr.cz/eet/schema/v3 Trzba"`
	Hlavicka      Hlavicka      `xml:"Hlavicka"`
	Data          Data          `xml:"Data"`
	KontrolniKody KontrolniKody `xml:"KontrolniKody"`
}

type Hlavicka struct {
	XMLName      xml.Name     `xml:"Hlavicka"`
	DatOdesl     DateTimeType `xml:"dat_odesl,attr"`
	Overeni      bool         `xml:"overeni,attr,omitempty"`
	PrvniZaslani bool         `xml:"prvni_zaslani,attr"`
	UuidZpravy   UUIDType     `xml:"uuid_zpravy,attr"`
}

type Data struct {
	XMLName         xml.Name     `xml:"Data"`
	CelkTrzba       CastkaType   `xml:"celk_trzba,attr"`
	CerpZuct        CastkaType   `xml:"cerp_zuct,attr,omitempty"`
	CestSluz        CastkaType   `xml:"cest_sluz,attr,omitempty"`
	Dan1            CastkaType   `xml:"dan1,attr,omitempty"`
	Dan2            CastkaType   `xml:"dan2,attr,omitempty"`
	Dan3            CastkaType   `xml:"dan3,attr,omitempty"`
	DatTrzby        DateTimeType `xml:"dat_trzby,attr"`
	DicPopl         CZDICType    `xml:"dic_popl,attr"`
	DicPoverujiciho CZDICType    `xml:"dic_poverujiciho,attr,omitempty"`
	IdPokl          String20     `xml:"id_pokl,attr"`
	IdProvoz        IdProvozType `xml:"id_provoz,attr"`
	PoradCis        String25     `xml:"porad_cis,attr"`
	PouzitZboz1     CastkaType   `xml:"pouzit_zboz1,attr,omitempty"`
	PouzitZboz2     CastkaType   `xml:"pouzit_zboz2,attr,omitempty"`
	PouzitZboz3     CastkaType   `xml:"pouzit_zboz3,attr,omitempty"`
	Rezim           RezimType    `xml:"rezim,attr"`
	UrcenoCerpZuct  CastkaType   `xml:"urceno_cerp_zuct,attr,omitempty"`
	ZaklDan1        CastkaType   `xml:"zakl_dan1,attr,omitempty"`
	ZaklDan2        CastkaType   `xml:"zakl_dan2,attr,omitempty"`
	ZaklDan3        CastkaType   `xml:"zakl_dan3,attr,omitempty"`
	ZaklNepodlDph   CastkaType   `xml:"zakl_nepodl_dph,attr,omitempty"`
}

type KontrolniKody struct {
	XMLName xml.Name `xml:"KontrolniKody"`
	Pkp     Pkp      `xml:"pkp"`
	Bkp     Bkp      `xml:"bkp"`
}

type Pkp struct {
	XMLName  xml.Name        `xml:"pkp"`
	Cipher   PkpCipherType   `xml:"cipher,attr"`
	Digest   PkpDigestType   `xml:"digest,attr"`
	Encoding PkpEncodingType `xml:"encoding,attr"`
	Value    string          `xml:",chardata"`
}

func NewPkp(t Trzba, signer *signing.Signer) Pkp {
	pkpStr := fmt.Sprintf("%s|%s|%s|%s|%s|%s", t.Data.DicPopl, t.Data.IdProvoz, t.Data.IdPokl, t.Data.PoradCis, t.Data.DatTrzby, t.Data.CelkTrzba)
	pkp, _ := signer.Sign([]byte(pkpStr))
	return Pkp{
		Cipher:   PkpCipherTypeRSA2048,
		Digest:   PkpDigestTypeSHA256,
		Encoding: PkpEncodingTypeBase64,
		Value:    base64.StdEncoding.EncodeToString(pkp),
	}
}

type Bkp struct {
	XMLName  xml.Name        `xml:"bkp"`
	Digest   BkpDigestType   `xml:"digest,attr"`
	Encoding BkpEncodingType `xml:"encoding,attr"`
	Value    string          `xml:",chardata"`
}

func NewBkp(pkp Pkp) Bkp {
	pkpValue, _ := base64.StdEncoding.DecodeString(pkp.Value)
	sumBkp := sha1.Sum(pkpValue)
	encodedBkp := hex.EncodeToString(sumBkp[:])
	finalBkp := fmt.Sprintf("%s-%s-%s-%s-%s", encodedBkp[0:8], encodedBkp[8:16], encodedBkp[16:24], encodedBkp[24:32], encodedBkp[32:40])
	return Bkp{
		Digest:   BkpDigestTypeSHA1,
		Encoding: BkpEncodingTypeBase16,
		Value:    finalBkp,
	}
}
