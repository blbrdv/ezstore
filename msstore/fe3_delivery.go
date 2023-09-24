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
	msaToken         = "<Device>dAA9AEUAdwBBAHcAQQBzAE4AMwBCAEEAQQBVADEAYgB5AHMAZQBtAGIAZQBEAFYAQwArADMAZgBtADcAbwBXAHkASAA3AGIAbgBnAEcAWQBtAEEAQQBMAGoAbQBqAFYAVQB2AFEAYwA0AEsAVwBFAC8AYwBDAEwANQBYAGUANABnAHYAWABkAGkAegBHAGwAZABjADEAZAAvAFcAeQAvAHgASgBQAG4AVwBRAGUAYwBtAHYAbwBjAGkAZwA5AGoAZABwAE4AawBIAG0AYQBzAHAAVABKAEwARAArAFAAYwBBAFgAbQAvAFQAcAA3AEgAagBzAEYANAA0AEgAdABsAC8AMQBtAHUAcgAwAFMAdQBtAG8AMABZAGEAdgBqAFIANwArADQAcABoAC8AcwA4ADEANgBFAFkANQBNAFIAbQBnAFIAQwA2ADMAQwBSAEoAQQBVAHYAZgBzADQAaQB2AHgAYwB5AEwAbAA2AHoAOABlAHgAMABrAFgAOQBPAHcAYQB0ADEAdQBwAFMAOAAxAEgANgA4AEEASABzAEoAegBnAFQAQQBMAG8AbgBBADIAWQBBAEEAQQBpAGcANQBJADMAUQAvAFYASABLAHcANABBAEIAcQA5AFMAcQBhADEAQgA4AGsAVQAxAGEAbwBLAEEAdQA0AHYAbABWAG4AdwBWADMAUQB6AHMATgBtAEQAaQBqAGgANQBkAEcAcgBpADgAQQBlAEUARQBWAEcAbQBXAGgASQBCAE0AUAAyAEQAVwA0ADMAZABWAGkARABUAHoAVQB0AHQARQBMAEgAaABSAGYAcgBhAGIAWgBsAHQAQQBUAEUATABmAHMARQBGAFUAYQBRAFMASgB4ADUAeQBRADgAagBaAEUAZQAyAHgANABCADMAMQB2AEIAMgBqAC8AUgBLAGEAWQAvAHEAeQB0AHoANwBUAHYAdAB3AHQAagBzADYAUQBYAEIAZQA4AHMAZwBJAG8AOQBiADUAQQBCADcAOAAxAHMANgAvAGQAUwBFAHgATgBEAEQAYQBRAHoAQQBYAFAAWABCAFkAdQBYAFEARQBzAE8AegA4AHQAcgBpAGUATQBiAEIAZQBUAFkAOQBiAG8AQgBOAE8AaQBVADcATgBSAEYAOQAzAG8AVgArAFYAQQBiAGgAcAAwAHAAUgBQAFMAZQBmAEcARwBPAHEAdwBTAGcANwA3AHMAaAA5AEoASABNAHAARABNAFMAbgBrAHEAcgAyAGYARgBpAEMAUABrAHcAVgBvAHgANgBuAG4AeABGAEQAbwBXAC8AYQAxAHQAYQBaAHcAegB5AGwATABMADEAMgB3AHUAYgBtADUAdQBtAHAAcQB5AFcAYwBLAFIAagB5AGgAMgBKAFQARgBKAFcANQBnAFgARQBJADUAcAA4ADAARwB1ADIAbgB4AEwAUgBOAHcAaQB3AHIANwBXAE0AUgBBAFYASwBGAFcATQBlAFIAegBsADkAVQBxAGcALwBwAFgALwB2AGUATAB3AFMAawAyAFMAUwBIAGYAYQBLADYAagBhAG8AWQB1AG4AUgBHAHIAOABtAGIARQBvAEgAbABGADYASgBDAGEAYQBUAEIAWABCAGMAdgB1AGUAQwBKAG8AOQA4AGgAUgBBAHIARwB3ADQAKwBQAEgAZQBUAGIATgBTAEUAWABYAHoAdgBaADYAdQBXADUARQBBAGYAZABaAG0AUwA4ADgAVgBKAGMAWgBhAEYASwA3AHgAeABnADAAdwBvAG4ANwBoADAAeABDADYAWgBCADAAYwBZAGoATAByAC8ARwBlAE8AegA5AEcANABRAFUASAA5AEUAawB5ADAAZAB5AEYALwByAGUAVQAxAEkAeQBpAGEAcABwAGgATwBQADgAUwAyAHQANABCAHIAUABaAFgAVAB2AEMAMABQADcAegBPACsAZgBHAGsAeABWAG0AKwBVAGYAWgBiAFEANQA1AHMAdwBFAD0AJgBwAD0A</Device>"
)

func fe3Client() *resty.Client {
	return http().
		SetHeader("Content-Type", "application/soap+xml")
}

func GetCookie() (string, error) {
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

func GetWUID(id string, market string, lang string) (string, error) {
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
		FindOne(json, "//WuCategoryId").
		Value()

	return fmt.Sprintf("%v", aboba), nil
}

func GetProducts(cookie string, categoryIdentifier string) ([]string, error) {
	var list []string

	resp, err := fe3Client().
		R().
		SetBody(wuidRequest(msaToken, cookie, categoryIdentifier)).
		Post(clientUrl)

	if err != nil {
		return list, err
	}

	xml, err := xmlquery.Parse(strings.NewReader(resp.String()))

	if err != nil {
		return list, err
	}

	rawxml := strings.Replace(
		strings.Replace(
			strings.Replace(xml.OutputXML(true), "&lt;", "<", -1),
			"&gt;", ">", -1),
		"&#34;", "\"", -1)

	newxml, err := xmlquery.Parse(strings.NewReader(rawxml))

	if err != nil {
		return list, err
	}

	for _, element := range newxml.SelectElements("//SecuredFragment/../../UpdateIdentity") {
		num := element.SelectAttr("RevisionNumber")

		if num != "" {
			list = append(list, element.SelectAttr("UpdateID"))
		}
	}

	return list, nil
}
