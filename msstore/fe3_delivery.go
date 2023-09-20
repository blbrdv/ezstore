package msstore

import (
	"fmt"
	"strings"

	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
	"github.com/go-resty/resty/v2"
)

const (
	clientUrl        = "https://fe3.delivery.mp.microsoft.com/ClientWebService/client.asmx"
	clientSecuredUrl = "https://fe3.delivery.mp.microsoft.com/ClientWebService/client.asmx/secured"
	wuidInfoUrl      = "https://displaycatalog.mp.microsoft.com/v7.0/products/"
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
		SelectElement("//EncryptedData").
		InnerText(), nil
}

func getWUID(id string, market string, lang string) (string, error) {
	resp, err := fe3Client().
		R().
		Get(fmt.Sprintf("%s%s?market=%s&languages=%s-%s,%s,neutral", wuidInfoUrl, id, market, lang, market, lang))

	if err != nil {
		return "", err
	}

	json, err := jsonquery.Parse(strings.NewReader(resp.String()))

	if err != nil {
		return "", err
	}

	aboba := jsonquery.
		// FindOne(json, "//ProductListing/Product/DisplaySkuAvailabilities/[0]/Sku/Properties/FulfillmentData/WuCategoryId").
		FindOne(json, "//WuCategoryId").
		Value()

	return fmt.Sprintf("%v", aboba), nil
}
