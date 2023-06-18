package msstore

import (
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/go-resty/resty/v2"
)

const (
	clientUrl        = "https://fe3.delivery.mp.microsoft.com/ClientWebService/client.asmx"
	clientSecuredUrl = "https://fe3.delivery.mp.microsoft.com/ClientWebService/client.asmx/secured"
)

func fe3Client() *resty.Client {
	return http().
		SetHeader("Content-Type", "application/soap+xml")
}

func getCookie() (string, error) {
	resp, err := fe3Client().
		R().
		SetBody(getCookiePayload).
		Post(clientUrl)

	if err != nil {
		return "", err
	}

	xml, err := xmlquery.Parse(strings.NewReader(resp.String()))

	if err != nil {
		return "", err
	}

	return xml.
		SelectElement("//s:Envelope/s:Body/GetCookieResponse/GetCookieResult/EncryptedData").
		InnerText(), nil
}
