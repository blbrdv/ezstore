package store

import (
	"fmt"
	"github.com/imroc/req/v3"
	"iter"
	"maps"
	net "net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
	"github.com/blbrdv/ezstore/internal/ms"
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

func (p ProductInfo) String() string {
	return fmt.Sprintf("[ \"%s\"; \"%s\" ]", p.UpdateID, p.RevisionNumber)
}

type bundleInfo struct {
	Name string
	Id   string
}

func newBundleInfo(input string) (*bundleInfo, error) {
	bundleRegexp := regexp.MustCompile(`^([0-9a-zA-Z.-]+)_([a-z0-9]+)$`)
	matches := bundleRegexp.FindStringSubmatch(input)
	if len(matches) == 0 {
		return nil, fmt.Errorf("%s is not valid bundle info", input)
	}

	return &bundleInfo{Name: matches[1], Id: matches[2]}, nil
}

type bundleData struct {
	*bundleInfo

	Version *ms.Version
	Arch    string
	Format  string
	Url     string
}

func newBundleData(input string) (*bundleData, error) {
	bundleRegexp := regexp.MustCompile(`^([0-9a-zA-Z.-]+)_([\d\.]+)_([a-z0-9]+)_~?_([a-z0-9]+).([a-zA-Z]+)`)
	matches := bundleRegexp.FindStringSubmatch(input)
	if len(matches) == 0 {
		return nil, fmt.Errorf("%s is not valid bundle data", input)
	}

	info := &bundleInfo{Name: matches[1], Id: matches[4]}
	version, err := ms.NewVersion(matches[2])
	if err != nil {
		return nil, err
	}

	return &bundleData{
			bundleInfo: info,
			Version:    version,
			Arch:       strings.ToLower(matches[3]),
			Format:     strings.ToLower(matches[5]),
		},
		nil
}

func (bd *bundleData) String() string {
	return fmt.Sprintf(
		"{ %s, Name: %s, Version: %s, Architecture: %s, Format: %s }",
		bd.Id,
		bd.Name,
		bd.Version.String(),
		bd.Arch,
		bd.Format,
	)
}

type bundlesList []*bundleData

type bundles struct {
	bundlesList
}

func newBundles(bundle *bundleData) *bundles {
	return &bundles{bundlesList{bundle}}
}

func initBundles() *bundles {
	return &bundles{bundlesList{}}
}

func (b *bundles) Append(bundle *bundleData) {
	for _, value := range b.bundlesList {
		if value.String() == bundle.String() {
			return
		}
	}

	b.bundlesList = append(b.bundlesList, bundle)
}

func (b *bundles) GetSupported(arch ms.Architecture) (*bundleData, error) {
	for _, supported := range arch.CompatibleWith() {
		for _, data := range b.bundlesList {
			if data.Arch == supported.String() {
				return data, nil
			}
		}
	}

	for _, data := range b.bundlesList {
		if data.Arch == "neutral" {
			return data, nil
		}
	}

	return nil, fmt.Errorf("%s architecture is not supported by this app", arch.String())
}

type bundlesByVersion map[ms.Version]*bundles

type bundlesGroup struct {
	bundlesByVersion
}

func newBundlesGroup(bundle *bundleData) *bundlesGroup {
	return &bundlesGroup{bundlesByVersion{*bundle.Version: newBundles(bundle)}}
}

func initBundlesGroup() *bundlesGroup {
	return &bundlesGroup{bundlesByVersion{}}
}

func (bg *bundlesGroup) Add(bundle *bundleData) {
	version := *bundle.Version
	b := bg.bundlesByVersion[version]

	if b == nil {
		bg.bundlesByVersion[version] = newBundles(bundle)
	} else {
		b.Append(bundle)
	}
}

func (bg *bundlesGroup) Get(version *ms.Version, arch ms.Architecture) (*bundleData, error) {
	var searchVersion ms.Version
	if version == nil {
		versions := toSlice(maps.Keys(bg.bundlesByVersion))
		searchVersion = versions[0]
		for _, key := range versions[1:] {
			if key.Compare(&searchVersion) == 1 {
				searchVersion = key
			}
		}
	} else {
		searchVersion = *version
	}

	list := bg.bundlesByVersion[searchVersion]
	if list == nil {
		return nil, fmt.Errorf("can not get bundle by version %s", searchVersion.String())
	}

	return list.GetSupported(arch)
}

func (bg *bundlesGroup) GetLatest(arch ms.Architecture) (*bundleData, error) {
	return bg.Get(nil, arch)
}

type bundlesById map[string]*bundlesGroup

type bundlesMap struct {
	bundlesById
}

func (bm *bundlesMap) Add(bundle *bundleData) {
	group := bm.bundlesById[bundle.Id]

	if group == nil {
		bm.bundlesById[bundle.Id] = newBundlesGroup(bundle)
	} else {
		group.Add(bundle)
	}
}

func initBundleMap() *bundlesMap {
	return &bundlesMap{bundlesById{}}
}

func fe3Client() *req.Client {
	return client.SetCommonHeader("Content-Type", "application/soap+xml")
}

func getCookie() (string, error) {
	resp, err := fe3Client().R().SetBody(getCookiePayload).Post(clientURL)
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

func getWUID(id string, locale *ms.Locale) (*bundleInfo, string, error) {
	url := fmt.Sprintf(
		"%s%s?market=%s&languages=%s,%s,neutral",
		wuidInfoURL,
		id,
		locale.Country,
		locale.String(),
		locale.Language,
	)

	resp, err := fe3Client().R().Get(url)
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode == 404 {
		return nil, "", fmt.Errorf(`product with id "%s" not found`, id)
	}
	if resp.IsErrorState() {
		return nil, "", fmt.Errorf("server error: %s", resp.ErrorResult())
	}

	data, err := jsonquery.Parse(strings.NewReader(resp.String()))
	if err != nil {
		return nil, "", err
	}

	fulfillmentData := jsonquery.FindOne(data, "Product/DisplaySkuAvailabilities/*[1]/Sku/Properties/FulfillmentData")
	if fulfillmentData == nil {
		return nil, "", fmt.Errorf("can not find fulfillment data")
	}

	wuid := fmt.Sprintf("%v", jsonquery.FindOne(fulfillmentData, "WuCategoryId").Value())
	if wuid == "" {
		return nil, "", fmt.Errorf("can not find WUID")
	}

	packageName := fmt.Sprintf("%v", jsonquery.FindOne(fulfillmentData, "PackageFamilyName").Value())
	if packageName == "" {
		return nil, "", fmt.Errorf("can not find package name")
	}

	info, err := newBundleInfo(packageName)
	if err != nil {
		return nil, "", err
	}

	return info, wuid, nil
}

func getProducts(cookie string, categoryIdentifier string) ([]ProductInfo, error) {
	var list []ProductInfo

	resp, err := fe3Client().R().SetBody(wuidRequest(msaToken, cookie, categoryIdentifier)).Post(clientURL)
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
			list = append(list, ProductInfo{updateID, revisionNumber})
		}
	}

	return list, nil
}

func getURL(info ProductInfo) ([]string, error) {
	var result []string

	resp, err := fe3Client().R().
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

func getFileData(url string) (*bundleData, error) {
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

	data.Url = url

	return data, nil
}

func toSlice[T any](iter iter.Seq[T]) []T {
	var result []T

	for value := range iter {
		result = append(result, value)
	}

	return result
}

// Download backage and its dependencies from MS Store by id, version and locale to destination directory
// and returns array of backage and its dependencies paths.
func Download(id string, version *ms.Version, arch ms.Architecture, locale *ms.Locale, destinationPath string) ([]ms.FileInfo, error) {
	var sp *pterm.SpinnerPrinter
	var pb *pterm.ProgressbarPrinter

	sp, _ = pterm.DefaultSpinner.Start("Fetching cookie...")
	cookie, err := getCookie()
	if err != nil {
		return nil, err
	}
	sp.Success("Cookie fetched")

	sp, _ = pterm.DefaultSpinner.Start("Fetching product info...")
	appInfo, wuid, err := getWUID(id, locale)
	if err != nil {
		return nil, err
	}
	sp.Success("Product info fetched")

	sp, _ = pterm.DefaultSpinner.Start("Fetching product files...")
	productInfos, err := getProducts(cookie, wuid)
	if err != nil {
		return nil, err
	}

	appFiles := initBundlesGroup()
	depFiles := initBundleMap()
	for _, info := range productInfos {
		productURL, err := getURL(info)
		if err != nil {
			return nil, err
		}

		for _, url := range productURL {
			fileData, err := getFileData(url)
			if err != nil {
				return nil, err
			}

			if fileData.Format == "blockmap" {
				continue
			}

			if fileData.Id == appInfo.Id {
				appFiles.Add(fileData)
			} else {
				depFiles.Add(fileData)
			}
		}
	}

	downloads := initBundles()
	appFile, err := appFiles.Get(version, arch)
	if err != nil {
		return nil, err
	}
	downloads.Append(appFile)
	for _, deps := range depFiles.bundlesById {
		depFile, err := deps.GetLatest(arch)
		if err != nil {
			return nil, err
		}

		downloads.Append(depFile)
	}
	sp.Success("Product files fetched")

	pb, _ = pterm.DefaultProgressbar.WithTotal(len(downloads.bundlesList)).WithTitle("Fetching product files info...").Start()
	var result []ms.FileInfo
	for _, data := range downloads.bundlesList {
		fullPath := path.Join(
			destinationPath,
			fmt.Sprintf("%s-%s.%s", data.Name, data.Version.String(), data.Format),
		)

		file, err := os.OpenFile(fullPath, os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}

		_, err = client.R().SetOutput(file).Get(data.Url)
		if err != nil {
			return nil, err
		}

		result = append(result, ms.FileInfo{Path: fullPath, Name: data.Name, Version: data.Version})
		pb.Increment()
	}
	//TODO remove progress bar when done

	return result, nil
}
