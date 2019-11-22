package eet

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

const XmlHeader = `<?xml version="1.0" encoding="UTF-8"?>`

type SOAPEnvelopeRequest struct {
	XMLName   xml.Name   `xml:"soap:Envelope"`
	XmlnsSoap string     `xml:"xmlns:soap,attr"`
	Header    SOAPHeader `xml:"SOAP-ENV:Header"`
	Body      SOAPBody   `xml:"soap:Body"`
}

func NewSOAPEnvelopeRequest(content interface{}, signer *Signer) (SOAPEnvelopeRequest, error) {
	bodyId := "id-" + strings.ToUpper(hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes()))
	certId := "X509-" + strings.ToUpper(hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes()))
	sigId := "SIG-" + strings.ToUpper(hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes()))
	keyId := "KI-" + strings.ToUpper(hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes()))
	secTokenId := "STR-" + strings.ToUpper(hex.EncodeToString(uuid.Must(uuid.NewV4()).Bytes()))

	envelope := SOAPEnvelopeRequest{
		XmlnsSoap: NsSoapUrl,
		Header: SOAPHeader{
			XmlnsSoapEnv: NsSoapEnvUrl,
			Security: WsseSecurity{
				XmlnsWsse:      NsWsseUrl,
				XmlnsWsu:       NsWsuUrl,
				MustUnderstand: "1",
				BinarySecurityToken: WsseBinarySecurityToken{
					EncodingType: EncodingBase64Url,
					ValueType:    ValueX509Url,
					WsuId:        certId,
				},
				Signature: DsSignature{
					XmlnsDs: NsDsUrl,
					Id:      sigId,
					SignedInfo: DsSignedInfo{
						XmlnsDs:   NsDsUrl,
						XmlnsSoap: NsSoapUrl,
						CanonicalizationMethod: DsCanonicalizationMethod{
							Algorithm: NsEcUrl,
							InclusiveNamespaces: EcInclusiveNamespaces{
								XmlnsEc:    NsEcUrl,
								PrefixList: NsSoap,
							},
						},
						SignatureMethod: DsSignatureMethod{
							Algorithm: AlgorithmSHA256,
						},
						Reference: DsReference{
							URI: "#" + bodyId,
							Transforms: DsTransforms{
								Transform: DsTransform{
									Algorithm: NsEcUrl,
									InclusiveNamespaces: EcInclusiveNamespaces{
										XmlnsEc: NsEcUrl,
									},
								},
							},
							DigestMethod: DsDigestMethod{
								Algorithm: AlgorithmDigestSHA256,
							},
							DigestValue: DsDigestValue{},
						},
					},
					SignatureValue: DsSignatureValue{},
					KeyInfo: DsKeyInfo{
						Id: keyId,
						SecurityTokenReference: WsseSecurityTokenReference{
							XmlnsWsse: NsWsseUrl,
							XmlnsWsu:  NsWsuUrl,
							WsuId:     secTokenId,
							Reference: WsseReference{
								URI:       "#" + certId,
								ValueType: ValueX509Url,
							},
						},
					},
				},
			},
		},
		Body: SOAPBody{
			XmlnsSoap: NsSoapUrl,
			XmlnsWsu:  NsWsuUrl,
			WsuId:     bodyId,
			Content:   content,
		},
	}

	// BinarySecurityToken
	envelope.Header.Security.BinarySecurityToken.Value = signer.Base64Cert()

	// DigestValue
	body, err := xml.Marshal(envelope.Body)
	if err != nil {
		return SOAPEnvelopeRequest{}, errors.Wrap(err, "Failed to xml.Marshal Body")
	}
	bodySum := sha256.Sum256(body)
	envelope.Header.Security.Signature.SignedInfo.Reference.DigestValue.Value = base64.StdEncoding.EncodeToString(bodySum[:])

	// SignatureValue
	signedInfo, err := xml.Marshal(envelope.Header.Security.Signature.SignedInfo)
	if err != nil {
		return SOAPEnvelopeRequest{}, errors.Wrap(err, "Failed to xml.Marshal Header.Security.Signature.SignedInfo")
	}
	signedSignedInfo, err := signer.Sign(signedInfo)
	if err != nil {
		return SOAPEnvelopeRequest{}, errors.Wrap(err, "Failed to Sign xml.Marshaled Header.Security.Signature.SignedInfo")
	}
	envelope.Header.Security.Signature.SignatureValue.Value = base64.StdEncoding.EncodeToString(signedSignedInfo)

	return envelope, nil
}

func (req SOAPEnvelopeRequest) Marshal() ([]byte, error) {
	b, err := xml.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to xml.Marshal SOAPEnvelopeRequest")
	}

	return []byte(XmlHeader + string(b)), nil
}

type SOAPHeader struct {
	XMLName      xml.Name     `xml:"SOAP-ENV:Header"`
	XmlnsSoapEnv string       `xml:"xmlns:SOAP-ENV,attr"`
	Security     WsseSecurity `xml:"wsse:Security"`
}

type WsseSecurity struct {
	XMLName             xml.Name                `xml:"wsse:Security"`
	XmlnsWsse           string                  `xml:"xmlns:wsse,attr"`
	XmlnsWsu            string                  `xml:"xmlns:wsu,attr"`
	MustUnderstand      string                  `xml:"soap:mustUnderstand,attr"`
	BinarySecurityToken WsseBinarySecurityToken `xml:"wsse:BinarySecurityToken"`
	Signature           DsSignature             `xml:"ds:Signature"`
}

type DsSignature struct {
	XMLName        xml.Name         `xml:"ds:Signature"`
	XmlnsDs        string           `xml:"xmlns:ds,attr"`
	Id             string           `xml:"Id,attr"`
	SignedInfo     DsSignedInfo     `xml:"ds:SignedInfo"`
	SignatureValue DsSignatureValue `xml:"ds:SignatureValue"`
	KeyInfo        DsKeyInfo        `xml:"ds:KeyInfo"`
}

type DsSignedInfo struct {
	XMLName                xml.Name                 `xml:"ds:SignedInfo"`
	XmlnsDs                string                   `xml:"xmlns:ds,attr"`
	XmlnsSoap              string                   `xml:"xmlns:soap,attr"`
	CanonicalizationMethod DsCanonicalizationMethod `xml:"ds:CanonicalizationMethod"`
	SignatureMethod        DsSignatureMethod        `xml:"ds:SignatureMethod"`
	Reference              DsReference              `xml:"ds:Reference"`
}

type DsCanonicalizationMethod struct {
	XMLName             xml.Name              `xml:"ds:CanonicalizationMethod"`
	Algorithm           string                `xml:"Algorithm,attr"`
	InclusiveNamespaces EcInclusiveNamespaces `xml:"ec:InclusiveNamespaces"`
}

type EcInclusiveNamespaces struct {
	XMLName    xml.Name `xml:"ec:InclusiveNamespaces"`
	XmlnsEc    string   `xml:"xmlns:ec,attr"`
	PrefixList string   `xml:"PrefixList,attr"`
}

type DsSignatureMethod struct {
	XMLName   xml.Name `xml:"ds:SignatureMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

type DsReference struct {
	XMLName      xml.Name       `xml:"ds:Reference"`
	URI          string         `xml:"URI,attr"`
	Transforms   DsTransforms   `xml:"ds:Transforms"`
	DigestMethod DsDigestMethod `xml:"ds:DigestMethod"`
	DigestValue  DsDigestValue  `xml:"ds:DigestValue"`
}

type DsTransforms struct {
	XMLName   xml.Name    `xml:"ds:Transforms"`
	Transform DsTransform `xml:"ds:Transform"`
}

type DsTransform struct {
	XMLName             xml.Name              `xml:"ds:Transform"`
	Algorithm           string                `xml:"Algorithm,attr"`
	InclusiveNamespaces EcInclusiveNamespaces `xml:"ec:InclusiveNamespaces"`
}

type DsDigestMethod struct {
	XMLName   xml.Name `xml:"ds:DigestMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
}

type DsDigestValue struct {
	XMLName xml.Name `xml:"ds:DigestValue"`
	Value   string   `xml:",chardata"`
}

type DsSignatureValue struct {
	XMLName xml.Name `xml:"ds:SignatureValue"`
	Value   string   `xml:",chardata"`
}

type DsKeyInfo struct {
	XMLName                xml.Name                   `xml:"ds:KeyInfo"`
	Id                     string                     `xml:"Id,attr"`
	SecurityTokenReference WsseSecurityTokenReference `xml:"wsse:SecurityTokenReference"`
}

type WsseSecurityTokenReference struct {
	XMLName   xml.Name      `xml:"wsse:SecurityTokenReference"`
	XmlnsWsse string        `xml:"xmlns:wsse,attr"`
	XmlnsWsu  string        `xml:"xmlns:wsu,attr"`
	WsuId     string        `xml:"wsu:Id,attr"`
	Reference WsseReference `xml:"wsse:Reference"`
}

type WsseReference struct {
	XMLName   xml.Name `xml:"wsse:Reference"`
	URI       string   `xml:"URI,attr"`
	ValueType string   `xml:"ValueType,attr"`
}

type WsseBinarySecurityToken struct {
	XMLName      xml.Name `xml:"wsse:BinarySecurityToken"`
	EncodingType string   `xml:"EncodingType,attr"`
	ValueType    string   `xml:"ValueType,attr"`
	WsuId        string   `xml:"wsu:Id,attr"`
	Value        string   `xml:",chardata"`
}

type SOAPBody struct {
	XMLName   xml.Name    `xml:"soap:Body"`
	XmlnsSoap string      `xml:"xmlns:soap,attr"`
	XmlnsWsu  string      `xml:"xmlns:wsu,attr"`
	WsuId     string      `xml:"wsu:Id,attr"`
	Content   interface{} `xml:",omitempty"`
}
