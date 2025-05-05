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
