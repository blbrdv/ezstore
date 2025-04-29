package msstore

import (
	"fmt"
	net "net/url"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
	types "github.com/blbrdv/ezstore/internal"
	"github.com/go-resty/resty/v2"
	"github.com/pterm/pterm"
)

const (
	clientURL        = "https://fe3.delivery.mp.microsoft.com/ClientWebService/client.asmx"
	clientSecuredURL = "https://fe3.delivery.mp.microsoft.com/ClientWebService/client.asmx/secured"
	wuidInfoURL      = "https://displaycatalog.mp.microsoft.com/v7.0/products/"
	msaToken         = "<Device>dAA9AEUAdwBBAHcAQQBzAE4AMwBCAEEAQQBVADEAYgB5AHMAZQBtAGIAZQBEAFYAQwArADMAZgBtADcAbwBXAHkASAA3AGIAbgBnAEcAWQBtAEEAQQBMAGoAbQBqAFYAVQB2AFEAYwA0AEsAVwBFAC8AYwBDAEwANQBYAGUANABnAHYAWABkAGkAegBHAGwAZABjADEAZAAvAFcAeQAvAHgASgBQAG4AVwBRAGUAYwBtAHYAbwBjAGkAZwA5AGoAZABwAE4AawBIAG0AYQBzAHAAVABKAEwARAArAFAAYwBBAFgAbQAvAFQAcAA3AEgAagBzAEYANAA0AEgAdABsAC8AMQBtAHUAcgAwAFMAdQBtAG8AMABZAGEAdgBqAFIANwArADQAcABoAC8AcwA4ADEANgBFAFkANQBNAFIAbQBnAFIAQwA2ADMAQwBSAEoAQQBVAHYAZgBzADQAaQB2AHgAYwB5AEwAbAA2AHoAOABlAHgAMABrAFgAOQBPAHcAYQB0ADEAdQBwAFMAOAAxAEgANgA4AEEASABzAEoAegBnAFQAQQBMAG8AbgBBADIAWQBBAEEAQQBpAGcANQBJADMAUQAvAFYASABLAHcANABBAEIAcQA5AFMAcQBhADEAQgA4AGsAVQAxAGEAbwBLAEEAdQA0AHYAbABWAG4AdwBWADMAUQB6AHMATgBtAEQAaQBqAGgANQBkAEcAcgBpADgAQQBlAEUARQBWAEcAbQBXAGgASQBCAE0AUAAyAEQAVwA0ADMAZABWAGkARABUAHoAVQB0AHQARQBMAEgAaABSAGYAcgBhAGIAWgBsAHQAQQBUAEUATABmAHMARQBGAFUAYQBRAFMASgB4ADUAeQBRADgAagBaAEUAZQAyAHgANABCADMAMQB2AEIAMgBqAC8AUgBLAGEAWQAvAHEAeQB0AHoANwBUAHYAdAB3AHQAagBzADYAUQBYAEIAZQA4AHMAZwBJAG8AOQBiADUAQQBCADcAOAAxAHMANgAvAGQAUwBFAHgATgBEAEQAYQBRAHoAQQBYAFAAWABCAFkAdQBYAFEARQBzAE8AegA4AHQAcgBpAGUATQBiAEIAZQBUAFkAOQBiAG8AQgBOAE8AaQBVADcATgBSAEYAOQAzAG8AVgArAFYAQQBiAGgAcAAwAHAAUgBQAFMAZQBmAEcARwBPAHEAdwBTAGcANwA3AHMAaAA5AEoASABNAHAARABNAFMAbgBrAHEAcgAyAGYARgBpAEMAUABrAHcAVgBvAHgANgBuAG4AeABGAEQAbwBXAC8AYQAxAHQAYQBaAHcAegB5AGwATABMADEAMgB3AHUAYgBtADUAdQBtAHAAcQB5AFcAYwBLAFIAagB5AGgAMgBKAFQARgBKAFcANQBnAFgARQBJADUAcAA4ADAARwB1ADIAbgB4AEwAUgBOAHcAaQB3AHIANwBXAE0AUgBBAFYASwBGAFcATQBlAFIAegBsADkAVQBxAGcALwBwAFgALwB2AGUATAB3AFMAawAyAFMAUwBIAGYAYQBLADYAagBhAG8AWQB1AG4AUgBHAHIAOABtAGIARQBvAEgAbABGADYASgBDAGEAYQBUAEIAWABCAGMAdgB1AGUAQwBKAG8AOQA4AGgAUgBBAHIARwB3ADQAKwBQAEgAZQBUAGIATgBTAEUAWABYAHoAdgBaADYAdQBXADUARQBBAGYAZABaAG0AUwA4ADgAVgBKAGMAWgBhAEYASwA3AHgAeABnADAAdwBvAG4ANwBoADAAeABDADYAWgBCADAAYwBZAGoATAByAC8ARwBlAE8AegA5AEcANABRAFUASAA5AEUAawB5ADAAZAB5AEYALwByAGUAVQAxAEkAeQBpAGEAcABwAGgATwBQADgAUwAyAHQANABCAHIAUABaAFgAVAB2AEMAMABQADcAegBPACsAZgBHAGsAeABWAG0AKwBVAGYAWgBiAFEANQA1AHMAdwBFAD0AJgBwAD0A</Device>"
)

type ProductInfo struct {
	UpdateID       string
	RevisionNumber string
}

func fe3Client() *resty.Client {
	return http().
		SetHeader("Content-Type", "application/soap+xml")
}

func getCookie() (string, error) {
	resp, err :=
		execute("post", clientURL,
			fe3Client().
				R().
				SetBody(getCookiePayload),
		)

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("server error: %s", resp.Error())
	}

	data, err := xmlquery.Parse(strings.NewReader(resp.String()))

	if err != nil {
		return "", err
	}

	return data.
		SelectElement("//EncryptedData").
		InnerText(), nil
}

func getWUID(id string, locale *types.Locale) (string, error) {
	resp, err :=
		execute(
			"get",
			fmt.Sprintf(
				"%s%s?market=%s&languages=%s,%s,neutral",
				wuidInfoURL,
				id,
				locale.Country,
				locale.String(),
				locale.Language,
			),
			fe3Client().R(),
		)

	if err != nil {
		return "", err
	}

	if resp.StatusCode() == 404 {
		return "", fmt.Errorf(`product with id "%s" not found`, id)
	}

	if resp.IsError() {
		return "", fmt.Errorf("server error: %s", resp.Error())
	}

	data, err := jsonquery.Parse(strings.NewReader(resp.String()))

	if err != nil {
		return "", err
	}

	wuid := jsonquery.
		FindOne(data, "//WuCategoryId").
		Value()

	return fmt.Sprintf("%v", wuid), nil
}

func getProducts(cookie string, categoryIdentifier string) ([]ProductInfo, error) {
	var list []ProductInfo

	resp, err :=
		execute("post", clientURL,
			fe3Client().R().
				SetBody(wuidRequest(msaToken, cookie, categoryIdentifier)),
		)

	if err != nil {
		return list, err
	}

	if resp.IsError() {
		return list, fmt.Errorf("server error: %s", resp.Error())
	}

	data, err := xmlquery.Parse(strings.NewReader(resp.String()))

	if err != nil {
		return list, err
	}

	undeformedXMLStr := strings.Replace(
		strings.Replace(
			strings.Replace(data.OutputXML(true), "&lt;", "<", -1),
			"&gt;", ">", -1),
		"&#34;", "\"", -1)

	data, err = xmlquery.Parse(strings.NewReader(undeformedXMLStr))

	if err != nil {
		return list, err
	}

	for _, element := range data.SelectElements("//SecuredFragment/../../UpdateIdentity") {
		revisionNumber := element.SelectAttr("RevisionNumber")
		updateID := element.SelectAttr("UpdateID")

		if revisionNumber != "" {
			list = append(list, ProductInfo{updateID, revisionNumber})
		}
	}

	return list, nil
}

func getURL(info ProductInfo) (string, error) {
	resp, err :=
		execute("post", clientSecuredURL,
			fe3Client().R().
				SetBody(fe3FileURL(msaToken, info.UpdateID, info.RevisionNumber)),
		)

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("server error: %s", resp.Error())
	}

	data, err := xmlquery.Parse(strings.NewReader(resp.String()))

	if err != nil {
		return "", err
	}

	return data.SelectElement("//FileLocation/Url").InnerText(), nil
}

func getFileName(url string) (string, error) {
	uri, err := net.Parse(url)

	if err != nil {
		return "", err
	}

	url = fmt.Sprintf("http://%s%s?%s", uri.Host, uri.EscapedPath(), uri.Query().Encode())

	name, err :=
		execute("head", url,
			http().R().
				SetHeader("Connection", "Keep-Alive").
				SetHeader("Accept", "*/*").
				SetHeader("User-Agent", "Microsoft-Delivery-Optimization/10.0"),
		)

	if err != nil {
		return "", err
	}

	header := name.Header().Get("Content-Disposition")
	fileNameRegexp := regexp.MustCompile(`filename=(\S+)`)

	return fileNameRegexp.FindStringSubmatch(header)[1], nil
}

// Download backage and its dependencies from MS Store by id, version and locale to destination directory
// and returns array of backage and its dependencies paths.
func Download(id string, version *types.Version, arch types.Architecture, locale *types.Locale, destinationPath string) ([]string, error) {
	cookieSpinner, _ := pterm.DefaultSpinner.Start("Fetching cookie...")
	cookie, err := getCookie()
	if err != nil {
		cookieSpinner.Fail(err.Error())
		return nil, err
	}
	cookieSpinner.Success("Cookie fetched")

	WUIDSpinner, _ := pterm.DefaultSpinner.Start("Fetching product WUID...")
	wuid, err := getWUID(id, locale)
	if err != nil {
		WUIDSpinner.Fail(err.Error())
		return nil, err
	}
	WUIDSpinner.Success("WUID fetched")

	linksSpinner, _ := pterm.DefaultSpinner.Start("Fetching product links...")
	productInfos, err := getProducts(cookie, wuid)
	if err != nil {
		linksSpinner.Fail(err.Error())
		return nil, err
	}

	var urls []string
	for _, info := range productInfos {
		productURL, err := getURL(info)
		if err != nil {
			linksSpinner.Fail(err.Error())
			return nil, err
		}
		// we don't need .BlockMap files
		if !strings.HasPrefix(productURL, "http://dl.delivery.mp.microsoft.com") {
			urls = append(urls, productURL)
		}
	}
	linksSpinner.Success("Product links fetched")

	productsBar, _ := pterm.DefaultProgressbar.WithTotal(len(urls)).WithTitle("Fetching product files info...").Start()
	bundleRegexp := regexp.MustCompile(`^([0-9a-zA-Z.-]+)_([\d\.]+)_([a-z0-9]+)_~?_[a-z0-9]+.([a-zA-Z]+)`)
	var bundles Bundles
	for _, productURL := range urls {
		fileName, err := getFileName(productURL)
		if err != nil {
			_, _ = productsBar.Stop()
			return nil, err
		}

		bundleData := bundleRegexp.FindStringSubmatch(fileName)
		v, err := types.NewVersion(bundleData[2])
		if err != nil {
			_, _ = productsBar.Stop()
			return nil, err
		}

		newArch, err := types.NewArchitecture(bundleData[3])
		if err == nil {
			bundle := BundleData{Version: v, Name: bundleData[1], URL: productURL, Arch: newArch, Format: strings.ToLower(bundleData[4])}
			bundles = append(bundles, bundle)
		}
		productsBar.Increment()
	}
	_, _ = productsBar.Stop()

	var filteredBundles Bundles
	for _, bundle := range bundles {
		if bundle.Format == "appx" {
			if bundle.Arch == arch {
				found := false
				for index, file := range filteredBundles {
					if bundle.Name == file.Name {
						if bundle.Version.Compare(file.Version) >= 0 {
							filteredBundles[index] = bundle
						}
						found = true
						break
					}
				}

				if !found {
					filteredBundles = append(filteredBundles, bundle)
				}
			}
		} else {
			found := false
			for index, file := range filteredBundles {
				if bundle.Name == file.Name {
					if bundle.Version.Compare(file.Version) >= 0 {
						filteredBundles[index] = bundle
					}
					found = true
					break
				}
			}

			if !found {
				filteredBundles = append(filteredBundles, bundle)
			}
		}
	}

	sort.Slice(filteredBundles, func(i, j int) bool {
		return filteredBundles[i].Format == "appx"
	})

	filesBar, _ := pterm.DefaultProgressbar.WithTotal(len(filteredBundles)).WithTitle("Downloading product files...").Start()
	var result []string
	for _, bundle := range filteredBundles {
		fullPath := path.Join(
			destinationPath,
			fmt.Sprintf("%s-%s.%s", bundle.Name, bundle.Version.String(), bundle.Format),
		)

		_, err = execute("get", bundle.URL, http().R().SetOutput(fullPath))

		if err != nil {
			_, _ = filesBar.Stop()
			return nil, err
		}

		result = append(result, fullPath)
		filesBar.Increment()
	}
	_, _ = filesBar.Stop()

	return result, nil
}
