package store

import (
	"fmt"
	"github.com/imroc/req/v3"
	net "net/url"
	"regexp"
	"strings"

	"github.com/antchfx/xmlquery"
)

const (
	clientURL        = "https://fe3.delivery.mp.microsoft.com/ClientWebService/client.asmx"
	clientSecuredURL = "https://fe3.delivery.mp.microsoft.com/ClientWebService/client.asmx/secured"
	msaToken         = "<Device>dAA9AEUAdwBBAHcAQQBzAE4AMwBCAEEAQQBVADEAYgB5AHMAZQBtAGIAZQBEAFYAQwArADMAZgBtADcAbwBXAHkASAA3AGIAbgBnAEcAWQBtAEEAQQBMAGoAbQBqAFYAVQB2AFEAYwA0AEsAVwBFAC8AYwBDAEwANQBYAGUANABnAHYAWABkAGkAegBHAGwAZABjADEAZAAvAFcAeQAvAHgASgBQAG4AVwBRAGUAYwBtAHYAbwBjAGkAZwA5AGoAZABwAE4AawBIAG0AYQBzAHAAVABKAEwARAArAFAAYwBBAFgAbQAvAFQAcAA3AEgAagBzAEYANAA0AEgAdABsAC8AMQBtAHUAcgAwAFMAdQBtAG8AMABZAGEAdgBqAFIANwArADQAcABoAC8AcwA4ADEANgBFAFkANQBNAFIAbQBnAFIAQwA2ADMAQwBSAEoAQQBVAHYAZgBzADQAaQB2AHgAYwB5AEwAbAA2AHoAOABlAHgAMABrAFgAOQBPAHcAYQB0ADEAdQBwAFMAOAAxAEgANgA4AEEASABzAEoAegBnAFQAQQBMAG8AbgBBADIAWQBBAEEAQQBpAGcANQBJADMAUQAvAFYASABLAHcANABBAEIAcQA5AFMAcQBhADEAQgA4AGsAVQAxAGEAbwBLAEEAdQA0AHYAbABWAG4AdwBWADMAUQB6AHMATgBtAEQAaQBqAGgANQBkAEcAcgBpADgAQQBlAEUARQBWAEcAbQBXAGgASQBCAE0AUAAyAEQAVwA0ADMAZABWAGkARABUAHoAVQB0AHQARQBMAEgAaABSAGYAcgBhAGIAWgBsAHQAQQBUAEUATABmAHMARQBGAFUAYQBRAFMASgB4ADUAeQBRADgAagBaAEUAZQAyAHgANABCADMAMQB2AEIAMgBqAC8AUgBLAGEAWQAvAHEAeQB0AHoANwBUAHYAdAB3AHQAagBzADYAUQBYAEIAZQA4AHMAZwBJAG8AOQBiADUAQQBCADcAOAAxAHMANgAvAGQAUwBFAHgATgBEAEQAYQBRAHoAQQBYAFAAWABCAFkAdQBYAFEARQBzAE8AegA4AHQAcgBpAGUATQBiAEIAZQBUAFkAOQBiAG8AQgBOAE8AaQBVADcATgBSAEYAOQAzAG8AVgArAFYAQQBiAGgAcAAwAHAAUgBQAFMAZQBmAEcARwBPAHEAdwBTAGcANwA3AHMAaAA5AEoASABNAHAARABNAFMAbgBrAHEAcgAyAGYARgBpAEMAUABrAHcAVgBvAHgANgBuAG4AeABGAEQAbwBXAC8AYQAxAHQAYQBaAHcAegB5AGwATABMADEAMgB3AHUAYgBtADUAdQBtAHAAcQB5AFcAYwBLAFIAagB5AGgAMgBKAFQARgBKAFcANQBnAFgARQBJADUAcAA4ADAARwB1ADIAbgB4AEwAUgBOAHcAaQB3AHIANwBXAE0AUgBBAFYASwBGAFcATQBlAFIAegBsADkAVQBxAGcALwBwAFgALwB2AGUATAB3AFMAawAyAFMAUwBIAGYAYQBLADYAagBhAG8AWQB1AG4AUgBHAHIAOABtAGIARQBvAEgAbABGADYASgBDAGEAYQBUAEIAWABCAGMAdgB1AGUAQwBKAG8AOQA4AGgAUgBBAHIARwB3ADQAKwBQAEgAZQBUAGIATgBTAEUAWABYAHoAdgBaADYAdQBXADUARQBBAGYAZABaAG0AUwA4ADgAVgBKAGMAWgBhAEYASwA3AHgAeABnADAAdwBvAG4ANwBoADAAeABDADYAWgBCADAAYwBZAGoATAByAC8ARwBlAE8AegA5AEcANABRAFUASAA5AEUAawB5ADAAZAB5AEYALwByAGUAVQAxAEkAeQBpAGEAcABwAGgATwBQADgAUwAyAHQANABCAHIAUABaAFgAVAB2AEMAMABQADcAegBPACsAZgBHAGsAeABWAG0AKwBVAGYAWgBiAFEANQA1AHMAdwBFAD0AJgBwAD0A</Device>"
)

type productInfo struct {
	UpdateID       string
	RevisionNumber string
}

func getXMLClient() *req.Client {
	return client.SetCommonHeader("Content-Type", "application/soap+xml")
}

var xmlClient = getXMLClient()

func getCookie() (string, error) {
	resp, err := xmlClient.R().SetBody(getCookiePayload).Post(clientURL)
	if err != nil {
		return "", err
	}
	if resp.IsErrorState() {
		return "", fmt.Errorf("server error: %s", resp.ErrorResult())
	}

	data, err := xmlquery.Parse(strings.NewReader(resp.String()))
	if err != nil {
		return "", err
	}

	return data.SelectElement("//EncryptedData").InnerText(), nil
}

func getProducts(cookie string, categoryIdentifier string) ([]productInfo, error) {
	var list []productInfo

	resp, err := xmlClient.R().SetBody(wuidRequest(msaToken, cookie, categoryIdentifier)).Post(clientURL)
	if err != nil {
		return list, err
	}
	if resp.IsErrorState() {
		return list, fmt.Errorf("server error: %s", resp.ErrorResult())
	}

	data, err := xmlquery.Parse(strings.NewReader(resp.String()))
	if err != nil {
		return list, err
	}

	undeformedXMLStr := strings.Replace(
		strings.Replace(
			strings.Replace(
				data.OutputXML(true),
				"&lt;",
				"<",
				-1,
			),
			"&gt;",
			">",
			-1,
		),
		"&#34;",
		"\"",
		-1,
	)

	data, err = xmlquery.Parse(strings.NewReader(undeformedXMLStr))
	if err != nil {
		return list, err
	}

	for _, element := range data.SelectElements("//SecuredFragment/../../UpdateIdentity") {
		revisionNumber := element.SelectAttr("RevisionNumber")
		updateID := element.SelectAttr("UpdateID")

		if revisionNumber != "" {
			list = append(list, productInfo{updateID, revisionNumber})
		}
	}

	return list, nil
}

func getURL(info productInfo) ([]string, error) {
	var result []string

	resp, err := xmlClient.R().
		SetBody(fe3FileURL(msaToken, info.UpdateID, info.RevisionNumber)).
		Post(clientSecuredURL)
	if err != nil {
		return result, err
	}
	if resp.IsErrorState() {
		return result, fmt.Errorf("server error: %s", resp.ErrorResult())
	}

	data, err := xmlquery.Parse(strings.NewReader(resp.String()))
	if err != nil {
		return result, err
	}

	for _, node := range data.SelectElements("//FileLocation") {
		result = append(result, node.SelectElement("Url").InnerText())
	}

	return result, nil
}

func getProductBundle(url string) (*bundleData, error) {
	uri, err := net.Parse(url)
	if err != nil {
		return nil, err
	}

	url = fmt.Sprintf("http://%s%s?%s", uri.Host, uri.EscapedPath(), uri.Query().Encode())

	res, err := client.
		SetCommonHeader("Connection", "Keep-Alive").
		SetCommonHeader("Accept", "*/*").
		SetCommonHeader("User-Agent", "Microsoft-Delivery-Optimization/10.0").
		R().
		Head(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("requiest error: %s", res.Status)
	}

	header := res.Header.Get("Content-Disposition")
	if header == "" {
		return nil, fmt.Errorf("can not get file name")
	}

	fileNameRegexp := regexp.MustCompile(`filename=(\S+)`)
	matches := fileNameRegexp.FindStringSubmatch(header)
	if len(matches) != 2 {
		return nil, fmt.Errorf("can not get file name")
	}

	data, err := newBundleData(matches[1])
	if err != nil {
		return nil, err
	}

	data.URL = url

	return data, nil
}
