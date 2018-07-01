package wsse

const (
	NsSoapUrl    = "http://schemas.xmlsoap.org/soap/envelope/"
	NsSoapEnvUrl = "http://schemas.xmlsoap.org/soap/envelope/"
	NsWsseUrl    = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	NsWsuUrl     = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
	NsDsUrl      = "http://www.w3.org/2000/09/xmldsig#"
	NsEcUrl      = "http://www.w3.org/2001/10/xml-exc-c14n#"
	NsEetUrl     = "http://fs.mfcr.cz/eet/schema/v3"

	NsSoap    = "soap"
	NsSoapEnv = "SOAP-ENV"
	NsWsse    = "wsse"
	NsWsu     = "wsu"
	NsDs      = "ds"
	NsEc      = "ec"

	EncodingBase64Url     = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary"
	ValueX509Url          = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509v3"
	AlgorithmC14N         = "http://www.w3.org/2001/10/xml-exc-c14n#"
	AlgorithmSHA256       = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"
	AlgorithmDigestSHA256 = "http://www.w3.org/2001/04/xmlenc#sha256"
)
