package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/blbrdv/ezstore/internal/ms"
	"github.com/blbrdv/ezstore/internal/utils"
	net "net/url"
	"os"
	"regexp"
)

func getProductBundle(url string) (*bundleData, error) {
	uri, err := net.Parse(url)
	if err != nil {
		return nil, err
	}

	query := uri.Query().Encode()
	var queryStr string
	if query == "" {
		queryStr = query
	} else {
		queryStr = fmt.Sprintf("?%s", query)
	}
	url = fmt.Sprintf("https://%s%s%s", uri.Host, uri.EscapedPath(), queryStr)
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

// Download backage and its dependencies from MS Store by id, version and locale to destination directory
// and returns array of backage and its dependencies paths.
func Download(id string, version *ms.Version, arch ms.Architecture, locale *ms.Locale, destinationPath string) ([]ms.FileInfo, error) {
	log.Debug("Fetching cookie...")
	cookie, err := getCookie()
	if err != nil {
		return nil, err
	}
	log.Info("Cookie fetched")

	log.Debug("Fetching product info...")
	appIndo, wuid, err := getWUID(id, locale)
	if err != nil {
		return nil, err
	}
	log.Info("Product info fetched")

	log.Debug("Fetching product files...")
	productsInfo, err := getProducts(cookie, wuid)
	if err != nil {
		return nil, err
	}

	appBundles := initBundlesGroup()
	depBundles := initBundleMap()
	for _, info := range productsInfo {
		productURLs, err := getURL(info)
		if err != nil {
			return nil, err
		}

		for _, url := range productURLs {
			productBundle, err := getProductBundle(url)
			if err != nil {
				return nil, err
			}

			if productBundle.Format == "blockmap" {
				continue
			}

			if productBundle.ID == appIndo.ID {
				appBundles.Add(productBundle)
			} else {
				depBundles.Add(productBundle)
			}
		}
	}

	bundlesToDownload := initBundles()
	for _, deps := range depBundles.bundlesByID {
		depBundle, err := deps.GetLatest(arch)
		if err != nil {
			return nil, err
		}

		bundlesToDownload.Append(depBundle)
	}
	appBundle, err := appBundles.Get(version, arch)
	if err != nil {
		return nil, err
	}
	bundlesToDownload.Append(appBundle)
	log.Info("Product files fetched")

	log.Debug("Fetching product files info...")
	var result []ms.FileInfo
	for _, data := range bundlesToDownload.bundlesList {
		fullPath := utils.Join(
			destinationPath,
			fmt.Sprintf("%s-%s.%s", data.Name, data.Version.String(), data.Format),
		)

		file, err := os.OpenFile(fullPath, os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}

		_, err = client.R().SetOutput(file).Get(data.URL)
		_ = file.Close()
		if err != nil {
			return nil, err
		}

		result = append(result, ms.FileInfo{Path: fullPath, Name: data.Name, Version: data.Version})
	}
	log.Info("Product files info fetched")

	return result, nil
}
