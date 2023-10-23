package msstore

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
	"github.com/blbrdv/ezstore/internal/types"
	"github.com/go-resty/resty/v2"
	"github.com/pterm/pterm"
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

	if resp.IsError() {
		return "", fmt.Errorf("server error: %s", resp.Error())
	}

	xml, err := xmlquery.Parse(strings.NewReader(resp.String()))

	if err != nil {
		return "", err
	}

	return xml.
		SelectElement("//EncryptedData").
		InnerText(), nil
}

func getWUID(id string, locale string) (string, error) {
	localeRaw := strings.Split(locale, "_")
	resp, err := fe3Client().
		R().
		Get(fmt.Sprintf("%s%s?market=%s&languages=%s-%s,%s,neutral", wuidInfoUrl, id, localeRaw[1], localeRaw[0], localeRaw[1], localeRaw[0]))

	if err != nil {
		return "", err
	}

	if resp.StatusCode() == 404 {
		return "", fmt.Errorf(`product with id "%s" not found`, id)
	}

	if resp.IsError() {
		return "", fmt.Errorf("server error: %s", resp.Error())
	}

	json, err := jsonquery.Parse(strings.NewReader(resp.String()))

	if err != nil {
		return "", err
	}

	wuid := jsonquery.
		FindOne(json, "//WuCategoryId").
		Value()

	return fmt.Sprintf("%v", wuid), nil
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

	if resp.IsError() {
		return list, fmt.Errorf("server error: %s", resp.Error())
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

	if resp.IsError() {
		return "", fmt.Errorf("server error: %s", resp.Error())
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

func Download(id string, version string, arch string, locale string, destinationPath string) ([]string, error) {
	sCoockie, _ := pterm.DefaultSpinner.Start("Fetching cookie...")
	cookie, err := getCookie()
	if err != nil {
		return nil, err
	}
	sCoockie.Success("Cookie fetched")

	sWUID, _ := pterm.DefaultSpinner.Start("Fetching product WUID...")
	wuid, err := getWUID(id, locale)
	if err != nil {
		return nil, err
	}
	sWUID.Success("WUID fetched")

	sLinks, _ := pterm.DefaultSpinner.Start("Fetching product links...")
	productInfos, err := getProducts(cookie, wuid)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, info := range productInfos {
		urlstr, err := getUrl(info)
		if err != nil {
			return nil, err
		}
		// we don't need .BlockMap files
		if !strings.HasPrefix(urlstr, "http://dl.delivery.mp.microsoft.com") {
			urls = append(urls, urlstr)
		}
	}
	sLinks.Success("Product links fetched")

	productsBar, _ := pterm.DefaultProgressbar.WithTotal(len(urls)).WithTitle("Fetching product files info...").Start()
	regex := regexp.MustCompile(`^([0-9a-zA-Z.-]+)_([\d\.]+)_([a-z0-9]+)_~?_[a-z0-9]+.([a-zA-Z]+)`)
	var bundles types.Bundles
	for _, urlobj := range urls {
		name, err := getFileName(urlobj)
		if err != nil {
			return nil, err
		}

		regexData := regex.FindStringSubmatch(name)
		v, err := types.New(regexData[2])
		if err != nil {
			return nil, err
		}
		bundle := types.BundleData{Version: v, Name: regexData[1], Url: urlobj, Arch: regexData[3], Format: strings.ToLower(regexData[4])}
		bundles = append(bundles, bundle)
		productsBar.Increment()
	}

	var files types.Bundles
	for _, bundle := range bundles {
		if bundle.Format == "appx" {
			if bundle.Arch == arch {
				found := false
				for index, file := range files {
					if bundle.Name == file.Name {
						if bundle.Version.Compare(*file.Version) >= 0 {
							files[index] = bundle
						}
						found = true
						break
					}
				}

				if !found {
					files = append(files, bundle)
				}
			}
		} else {
			found := false
			for index, file := range files {
				if bundle.Name == file.Name {
					if bundle.Version.Compare(*file.Version) >= 0 {
						files[index] = bundle
					}
					found = true
					break
				}
			}

			if !found {
				files = append(files, bundle)
			}
		}
	}

	filesBar, _ := pterm.DefaultProgressbar.WithTotal(len(files)).WithTitle("Downloading product files...").Start()
	var result []string
	for _, file := range files {
		fullPath := destinationPath + "\\" + file.Name + "-" + file.Version.String() + "." + file.Format

		_, err = http().R().
			SetOutput(fullPath).
			Get(file.Url)

		if err != nil {
			return nil, err
		}

		result = append(result, fullPath)
		filesBar.Increment()
	}

	return result, nil
}
