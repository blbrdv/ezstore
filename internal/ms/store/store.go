package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"github.com/pterm/pterm"
	net "net/url"
	"os"
	"path"
	"regexp"
)

func getProductBundle(url string) (*bundleData, error) {
	uri, err := net.Parse(url)
	if err != nil {
		return nil, err
	}

	res, err := client.
		SetCommonHeader("Connection", "Keep-Alive").
		SetCommonHeader("Accept", "*/*").
		SetCommonHeader("User-Agent", "Microsoft-Delivery-Optimization/10.0").
		R().
		SetPathParam("host", uri.Host).
		SetPathParam("path", uri.EscapedPath()).
		SetPathParam("query", uri.Query().Encode()).
		Head("https://{host}{path}?{query}")
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
	var sp *pterm.SpinnerPrinter
	var pb *pterm.ProgressbarPrinter

	sp, _ = pterm.DefaultSpinner.Start("Fetching cookie...")
	cookie, err := getCookie()
	if err != nil {
		return nil, err
	}
	sp.Success("Cookie fetched")

	sp, _ = pterm.DefaultSpinner.Start("Fetching product info...")
	appIndo, wuid, err := getWUID(id, locale)
	if err != nil {
		return nil, err
	}
	sp.Success("Product info fetched")

	sp, _ = pterm.DefaultSpinner.Start("Fetching product files...")
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
	sp.Success("Product files fetched")

	pb, _ = pterm.DefaultProgressbar.WithTotal(len(bundlesToDownload.bundlesList)).WithTitle("Fetching product files info...").Start()
	var result []ms.FileInfo
	for _, data := range bundlesToDownload.bundlesList {
		fullPath := path.Join(
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
		pb.Increment()
	}
	//TODO remove progress bar when done

	return result, nil
}
