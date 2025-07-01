package store

import (
	_ "embed"
	"fmt"
	"github.com/imroc/req/v3"
	"strings"

	"github.com/antchfx/xmlquery"
)

const (
	clientURL        = "https://fe3.delivery.mp.microsoft.com/ClientWebService/client.asmx"
	clientSecuredURL = "https://fe3.delivery.mp.microsoft.com/ClientWebService/client.asmx/secured"
)

//go:embed token.txt
var msaToken string

type productInfo struct {
	UpdateID       string
	RevisionNumber string
}

func (p productInfo) String() string {
	return fmt.Sprintf("%s %s", p.UpdateID, p.RevisionNumber)
}

func getXMLClient() *req.Client {
	return client.SetCommonHeader("Content-Type", "application/soap+xml")
}

var xmlClient = getXMLClient()

func getCookie() (string, error) {
	resp, err := xmlClient.R().SetBody(getCookiePayload).Post(clientURL)
	if err != nil {
		return "", fmt.Errorf("can not get cookie: POST %s: %s", clientURL, err.Error())
	}
	if resp.IsErrorState() {
		return "", fmt.Errorf("can not get cookie: POST %s: server returns error: %s", clientURL, resp.Status)
	}

	data, err := xmlquery.Parse(strings.NewReader(resp.String()))
	if err != nil {
		return "", fmt.Errorf("can not get cookie: can not parse result: %s", err.Error())
	}

	return data.SelectElement("//EncryptedData").InnerText(), nil
}

func getProducts(cookie string, categoryIdentifier string) ([]productInfo, error) {
	var list []productInfo

	resp, err := xmlClient.R().SetBody(wuidRequest(msaToken, cookie, categoryIdentifier)).Post(clientURL)
	if err != nil {
		return list, fmt.Errorf("can not get products files: POST %s: %s", clientURL, err.Error())
	}
	if resp.IsErrorState() {
		return list, fmt.Errorf("can not get products files: POST %s: server returns error: %s", clientURL, resp.Status)
	}

	data, err := xmlquery.Parse(strings.NewReader(resp.String()))
	if err != nil {
		return list, fmt.Errorf("can not get products files: can not parse result: %s", err.Error())
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
		return list, fmt.Errorf("can not get products files: can not parse result: %s", err.Error())
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
		return result, fmt.Errorf("can not get file url: POST %s: %s", clientSecuredURL, err.Error())
	}
	if resp.IsErrorState() {
		return result, fmt.Errorf("can not get file url: POST %s: server returns error: %s", clientSecuredURL, resp.Status)
	}

	data, err := xmlquery.Parse(strings.NewReader(resp.String()))
	if err != nil {
		return result, fmt.Errorf("can not get file url: can not parse result: %s", err.Error())
	}

	for _, node := range data.SelectElements("//FileLocation") {
		result = append(result, node.SelectElement("Url").InnerText())
	}

	return result, nil
}
