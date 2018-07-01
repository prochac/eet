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

const ValidMessage string = `<?xml version="1.0" encoding="UTF-8"?><soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><SOAP-ENV:Header xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"><wsse:Security xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" soap:mustUnderstand="1"><wsse:BinarySecurityToken EncodingType="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary" ValueType="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3" wsu:Id="X509-A79845F15C5549CA0514761283545351">MIIEmDCCA4CgAwIBAgIEdHOXJzANBgkqhkiG9w0BAQsFADB3MRIwEAYKCZImiZPyLGQBGRYCQ1oxQzBBBgNVBAoMOsSMZXNrw6EgUmVwdWJsaWthIOKAkyBHZW5lcsOhbG7DrSBmaW5hbsSNbsOtIMWZZWRpdGVsc3R2w60xHDAaBgNVBAMTE0VFVCBDQSAxIFBsYXlncm91bmQwHhcNMTYwOTMwMDkwMzU5WhcNMTkwOTMwMDkwMzU5WjBDMRIwEAYKCZImiZPyLGQBGRYCQ1oxEzARBgNVBAMTCkNaMDAwMDAwMTkxGDAWBgNVBA0TD3ByYXZuaWNrYSBvc29iYTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJnNUPW8rAlLi2KAwu12W1vqLj02mWIifq/Jp0/tUjf9B8RpkDAD3GOqDdVuHSfxej92WiEouDy7X8uXzIDdZu4pXA3t3KntxM8rAlu2U6SqtF3kTR+AJCdwfkM53U3z4/qoyKqdQ8lGuMxJKs7X5uIjcY/UDSXMK9OTmXRhndjYcX1oILr5F2ONf1Z0kWyl/S9wI0cl0gQ1F91mzqgnlH80u2inMmmBp42ndR4TGS1nvjer5D73bkLg07TdeqnUg609WwjUJN96OKZMsKXzBMzt09NbhQcABWnAWbRTSVhsAdDO8vfmWx2C+gXUlkIvtO+9fbj81GS1xdNoAkpARUcCAwEAAaOCAV4wggFaMAkGA1UdEwQCMAAwHQYDVR0OBBYEFL/0b0Iw6FY33UT8iJEy1V7nZVR6MB8GA1UdIwQYMBaAFHwwdqzM1ofR7Mkf4nAILONf3gwHMA4GA1UdDwEB/wQEAwIGwDBjBgNVHSAEXDBaMFgGCmCGSAFlAwIBMAEwSjBIBggrBgEFBQcCAjA8DDpUZW50byBjZXJ0aWZpa8OhdCBieWwgdnlkw6FuIHBvdXplIHBybyB0ZXN0b3ZhY8OtIMO6xI1lbHkuMIGXBgNVHR8EgY8wgYwwgYmggYaggYOGKWh0dHA6Ly9jcmwuY2ExLXBnLmVldC5jei9lZXRjYTFwZy9hbGwuY3JshipodHRwOi8vY3JsMi5jYTEtcGcuZWV0LmN6L2VldGNhMXBnL2FsbC5jcmyGKmh0dHA6Ly9jcmwzLmNhMS1wZy5lZXQuY3ovZWV0Y2ExcGcvYWxsLmNybDANBgkqhkiG9w0BAQsFAAOCAQEAvXdWsU+Ibd1VysKnjoy6RCYVcI9+oRUSSTvQQDJLFjwn5Sm6Hebhci8ERGwAzd2R6uqPdzl1KCjmHOitypZ66e+/e9wj3BaDqgBKRZYvxZykaVUdtQgG0819JZmiXTbGgOCKiUPIXO80cnP7U1ZPkVNV7WZwh0I2k/fg1VLTI5HA/x4BeD77wiEOExa7eqePJET0jpTVK3LxSW59LLIJROh4/kfKQbTvDL5Ypw8WagAMVCPvWnGJIcUru+ApLU4pZD9bdHSa1Ib4LpFhtWrkHYM/XqKbj2bNKKjTo5T3sU0Bf2QD3QzkmcjlNVG0V+qAgimwTdPueU/mtExw+7z1/A==</wsse:BinarySecurityToken><ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="SIG-A79845F15C5549CA0514761283545705"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"><ec:InclusiveNamespaces xmlns:ec="http://www.w3.org/2001/10/xml-exc-c14n#" PrefixList="soap"/></ds:CanonicalizationMethod><ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"/><ds:Reference URI="#id-A79845F15C5549CA0514761283545594"><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"><ec:InclusiveNamespaces xmlns:ec="http://www.w3.org/2001/10/xml-exc-c14n#" PrefixList=""/></ds:Transform></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha256"/><ds:DigestValue>M8/dBI/LLuwxP8ZoeRKVVpzIhQhyUDI0l6bglhNWKC0=</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue>DlFiF51sdtc0zeqgCsuSY6EU5emX7Hka6Ox3gviR4dpqyrwj6O8cm4oWkUTl+erINf9CpOG2y1z5y83+DubuGmiOPsACeEZjwF5TCme/uU1tzXs+LsLS8WXZvuMUejFOrUFEUKyzRvHJW5lskV/DhiMsTwJ/MEtGtjRVuWnAEcHxq/3ALsc3HfSi/qAqAiGlz3OpIKf3Hai6iiD37eJwHNd9QdJy4D32DJwY4Gw21ku7TO3FiLdtBT+Xdr3At8sWbbihFwErKrutT/FAVQKffUzCjHdEm9//W/shHufxYdnEh0m8cYp5Z5xK2bfnfhaKPB0rUZrrcw8T7LlvZKGogg==</ds:SignatureValue><ds:KeyInfo Id="KI-A79845F15C5549CA0514761283545482"><wsse:SecurityTokenReference xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" wsu:Id="STR-A79845F15C5549CA0514761283545513"><wsse:Reference URI="#X509-A79845F15C5549CA0514761283545351" ValueType="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3"/></wsse:SecurityTokenReference></ds:KeyInfo></ds:Signature></wsse:Security></SOAP-ENV:Header><soap:Body xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" wsu:Id="id-A79845F15C5549CA0514761283545594"><Trzba xmlns="http://fs.mfcr.cz/eet/schema/v3"><Hlavicka dat_odesl="2016-08-19T19:06:37+02:00" prvni_zaslani="false" uuid_zpravy="2da635a5-d712-459d-9674-c12f335c39f7"/><Data celk_trzba="34113.00" cerp_zuct="679.00" cest_sluz="5460.00" dan1="-172.39" dan2="-530.73" dan3="975.65" dat_trzby="2016-08-05T00:30:12+02:00" dic_popl="CZ00000019" dic_poverujiciho="CZ683555118" id_pokl="/5546/RO24" id_provoz="273" porad_cis="0/6460/ZQ42" pouzit_zboz1="784.00" pouzit_zboz2="967.00" pouzit_zboz3="189.00" rezim="0" urceno_cerp_zuct="324.00" zakl_dan1="-820.92" zakl_dan2="-3538.20" zakl_dan3="9756.46" zakl_nepodl_dph="3036.00"/><KontrolniKody><pkp cipher="RSA2048" digest="SHA256" encoding="base64">a0asEiJhFCBlVtptSspKvEZhcrvnzF7SQ55C4DhnStnSu1b37GUI2+Dlme9P94UCPZ1oCUPJdsYOBZ3IX6aEgEe0FJKXYX0kXraYCJKIo3g64wRchE7iblIOBCK1uHh8qqHA66Isnhb6hqBOOdlt2aWO/0jCzlfeQr0axpPF1mohMnP3h3ICaxZh0dnMdju5OmMrq+91PL5T9KkR7bfGHqAoWJ0kmxY/mZumtRfGil2/xf7I5pdVeYXPgDO/Tojzm6J95n68fPDOXTDrTzKYmqDjpg3kmWepLNQKFXRmkQrkBLToJWG1LDUDm3UTTmPWzq4c0XnGcXJDZglxfolGpA==</pkp><bkp digest="SHA1" encoding="base16">9356D566-A3E48838-FB403790-D201244E-95DCBD92</bkp></KontrolniKody></Trzba></soap:Body></soap:Envelope>`

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

func NewDispatcher(service Service, keyPath, certPath, password string) (eet *Dispatcher, err error) {
	d := Dispatcher{service: service}
	// TODO wait for pkcs12.DecodeAll func
	if d.signer, err = signing.NewSigner("", password); err != nil {
		d.signer = &signing.Signer{}

		if err = d.signer.LoadCertificate(certPath); err != nil {
			return nil, err
		}
		if err := d.signer.LoadPrivateKey(keyPath); err != nil {
			return nil, err
		}
	}

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
