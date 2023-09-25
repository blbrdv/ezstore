package msstore

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
	"github.com/blbrdv/ezstore/msver"
	"github.com/go-resty/resty/v2"
)

const (
	clientUrl        = "https://fe3.delivery.mp.microsoft.com/ClientWebService/client.asmx"
	clientSecuredUrl = "https://fe3.delivery.mp.microsoft.com/ClientWebService/client.asmx/secured"
	wuidInfoUrl      = "https://displaycatalog.mp.microsoft.com/v7.0/products/"
	msaToken         = "<Device>dAA9AEUAdwBBAHcAQQBzAE4AMwBCAEEAQQBVADEAYgB5AHMAZQBtAGIAZQBEAFYAQwArADMAZgBtADcAbwBXAHkASAA3AGIAbgBnAEcAWQBtAEEAQQBMAGoAbQBqAFYAVQB2AFEAYwA0AEsAVwBFAC8AYwBDAEwANQBYAGUANABnAHYAWABkAGkAegBHAGwAZABjADEAZAAvAFcAeQAvAHgASgBQAG4AVwBRAGUAYwBtAHYAbwBjAGkAZwA5AGoAZABwAE4AawBIAG0AYQBzAHAAVABKAEwARAArAFAAYwBBAFgAbQAvAFQAcAA3AEgAagBzAEYANAA0AEgAdABsAC8AMQBtAHUAcgAwAFMAdQBtAG8AMABZAGEAdgBqAFIANwArADQAcABoAC8AcwA4ADEANgBFAFkANQBNAFIAbQBnAFIAQwA2ADMAQwBSAEoAQQBVAHYAZgBzADQAaQB2AHgAYwB5AEwAbAA2AHoAOABlAHgAMABrAFgAOQBPAHcAYQB0ADEAdQBwAFMAOAAxAEgANgA4AEEASABzAEoAegBnAFQAQQBMAG8AbgBBADIAWQBBAEEAQQBpAGcANQBJADMAUQAvAFYASABLAHcANABBAEIAcQA5AFMAcQBhADEAQgA4AGsAVQAxAGEAbwBLAEEAdQA0AHYAbABWAG4AdwBWADMAUQB6AHMATgBtAEQAaQBqAGgANQBkAEcAcgBpADgAQQBlAEUARQBWAEcAbQBXAGgASQBCAE0AUAAyAEQAVwA0ADMAZABWAGkARABUAHoAVQB0AHQARQBMAEgAaABSAGYAcgBhAGIAWgBsAHQAQQBUAEUATABmAHMARQBGAFUAYQBRAFMASgB4ADUAeQBRADgAagBaAEUAZQAyAHgANABCADMAMQB2AEIAMgBqAC8AUgBLAGEAWQAvAHEAeQB0AHoANwBUAHYAdAB3AHQAagBzADYAUQBYAEIAZQA4AHMAZwBJAG8AOQBiADUAQQBCADcAOAAxAHMANgAvAGQAUwBFAHgATgBEAEQAYQBRAHoAQQBYAFAAWABCAFkAdQBYAFEARQBzAE8AegA4AHQAcgBpAGUATQBiAEIAZQBUAFkAOQBiAG8AQgBOAE8AaQBVADcATgBSAEYAOQAzAG8AVgArAFYAQQBiAGgAcAAwAHAAUgBQAFMAZQBmAEcARwBPAHEAdwBTAGcANwA3AHMAaAA5AEoASABNAHAARABNAFMAbgBrAHEAcgAyAGYARgBpAEMAUABrAHcAVgBvAHgANgBuAG4AeABGAEQAbwBXAC8AYQAxAHQAYQBaAHcAegB5AGwATABMADEAMgB3AHUAYgBtADUAdQBtAHAAcQB5AFcAYwBLAFIAagB5AGgAMgBKAFQARgBKAFcANQBnAFgARQBJADUAcAA4ADAARwB1ADIAbgB4AEwAUgBOAHcAaQB3AHIANwBXAE0AUgBBAFYASwBGAFcATQBlAFIAegBsADkAVQBxAGcALwBwAFgALwB2AGUATAB3AFMAawAyAFMAUwBIAGYAYQBLADYAagBhAG8AWQB1AG4AUgBHAHIAOABtAGIARQBvAEgAbABGADYASgBDAGEAYQBUAEIAWABCAGMAdgB1AGUAQwBKAG8AOQA4AGgAUgBBAHIARwB3ADQAKwBQAEgAZQBUAGIATgBTAEUAWABYAHoAdgBaADYAdQBXADUARQBBAGYAZABaAG0AUwA4ADgAVgBKAGMAWgBhAEYASwA3AHgAeABnADAAdwBvAG4ANwBoADAAeABDADYAWgBCADAAYwBZAGoATAByAC8ARwBlAE8AegA5AEcANABRAFUASAA5AEUAawB5ADAAZAB5AEYALwByAGUAVQAxAEkAeQBpAGEAcABwAGgATwBQADgAUwAyAHQANABCAHIAUABaAFgAVAB2AEMAMABQADcAegBPACsAZgBHAGsAeABWAG0AKwBVAGYAWgBiAFEANQA1AHMAdwBFAD0AJgBwAD0A</Device>"
)

type ProductInfo struct {
	UpdateId       string
	RevisionNumber string
}

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
		FindOne(json, "//WuCategoryId").
		Value()

	return fmt.Sprintf("%v", aboba), nil
}

func getProducts(cookie string, categoryIdentifier string) ([]ProductInfo, error) {
	var list []ProductInfo

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
		id := element.SelectAttr("UpdateID")

		if num != "" {
			list = append(list, ProductInfo{id, num})
		}
	}

	return list, nil
}

func getUrl(info ProductInfo) (string, error) {
	resp, err := fe3Client().
		R().
		SetBody(fe3FileUrl(msaToken, info.UpdateId, info.RevisionNumber)).
		Post(clientSecuredUrl)

	if err != nil {
		return "", err
	}

	xml, err := xmlquery.Parse(strings.NewReader(resp.String()))

	if err != nil {
		return "", err
	}

	return xml.SelectElement("//FileLocation/Url").InnerText(), nil
}

func getFileName(urlraw string) (string, error) {
	uri, err := url.Parse(urlraw)

	if err != nil {
		return "", err
	}

	fullurl := "http://" + uri.Host + uri.EscapedPath() + "?" + uri.Query().Encode()

	name, err := http().R().
		SetHeader("Connection", "Keep-Alive").
		SetHeader("Accept", "*/*").
		SetHeader("User-Agent", "Microsoft-Delivery-Optimization/10.0").
		Head(fullurl)

	if err != nil {
		return "", err
	}

	header := name.Header().Get("Content-Disposition")
	r := regexp.MustCompile(`filename=(\S+)`)

	return r.FindStringSubmatch(header)[1], nil
}

func Download(id string, version string, destinationPath string) (string, error) {
	fmt.Print("Getting cookies ...\n")

	cookie, err := getCookie()

	if err != nil {
		return "", err
	}

	fmt.Print("Getting product WUID ...\n")

	wuid, err := getWUID(id, "US", "en")

	if err != nil {
		return "", err
	}

	fmt.Print("Getting product urls ...\n")

	productInfos, err := getProducts(cookie, wuid)

	if err != nil {
		return "", err
	}

	var result []string

	for _, info := range productInfos {
		urlstr, err := getUrl(info)

		if err != nil {
			return "", err
		}

		// we don't need .BlockMap files
		if !strings.HasPrefix(urlstr, "http://dl.delivery.mp.microsoft.com") {
			result = append(result, urlstr)
		}
	}

	r := regexp.MustCompile(`^[a-zA-Z.-]+_([\d\.]+)_`)
	var bundles msver.Bundles

	for _, urlobj := range result {
		name, err := getFileName(urlobj)

		if err != nil {
			return "", err
		}

		if strings.HasSuffix(strings.ToLower(name), "bundle") {
			v, err := msver.New(r.FindStringSubmatch(name)[1])

			if err != nil {
				return "", err
			}

			bundles = append(bundles, msver.BundleData{v, name, urlobj})
		}
	}

	sort.Sort(bundles)
	var product msver.BundleData

	if version == "latest" {
		product = bundles[bundles.Len()-1]
	} else {
		var prodIndex int
		found := false

		for index, productInfo := range bundles {
			if productInfo.Version.String() == "v"+version {
				prodIndex = index
				found = true
				break
			}
		}

		if !found {
			return "", fmt.Errorf(`version "%s" not found`, version)
		}

		product = bundles[prodIndex]
	}

	fullPath := destinationPath + "\\" + product.Name

	if err != nil {
		return "", err
	}

	fmt.Printf(`Downloading product "%s"`, product.Name)
	fmt.Println("")

	_, err = http().R().
		SetOutput(fullPath).
		Get(product.Url)

	if err != nil {
		return "", err
	}

	return fullPath, nil
}
