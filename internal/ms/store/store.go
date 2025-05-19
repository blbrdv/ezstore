package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/blbrdv/ezstore/internal/ms"
	"github.com/blbrdv/ezstore/internal/utils"
	"os"
	"regexp"
)

func getProductName(url string) (string, error) {
	res, err := client.
		SetCommonHeader("Connection", "Keep-Alive").
		SetCommonHeader("Accept", "*/*").
		SetCommonHeader("User-Agent", "Microsoft-Delivery-Optimization/10.0").
		R().
		Head(url)
	if err != nil {
		return "", fmt.Errorf("fetching product name failed: HEAD %s: %s", url, err.Error())
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("fetching product name failed: HEAD %s: server responded with error: %s", url, res.Status)
	}

	header := res.Header.Get("Content-Disposition")
	if header == "" {
		return "", fmt.Errorf("fetching product name failed: can not get file name: response header \"Content-Disposition\" is empty")
	}

	fileNameRegexp := regexp.MustCompile(`filename=(\S+)`)
	matches := fileNameRegexp.FindStringSubmatch(header)
	if len(matches) != 2 {
		return "", fmt.Errorf("fetching product name failed: can not get file name: response header \"Content-Disposition\" has invalid format: %s", header)
	}

	return matches[1], nil
}

// Download backage and its dependencies from MS Store by id, version and locale to destination directory
// and returns array of backage and its dependencies paths.
func Download(id string, version *ms.Version, locale *ms.Locale, destinationPath string) ([]ms.FileInfo, error) {
	log.Debug("Fetching cookie...")
	cookie, err := getCookie()
	if err != nil {
		return nil, fmt.Errorf("can not fetch cookie: %s", err.Error())
	}
	log.Info("Cookie fetched")

	log.Debug("Fetching product info...")
	apps, wuid, err := getAppInfo(id, locale)
	if err != nil {
		return nil, fmt.Errorf("can not fetch product info: %s", err.Error())
	}
	log.Info("Product info fetched")

	log.Debug("Fetching product files...")
	productsInfo, err := getProducts(cookie, wuid)
	if err != nil {
		return nil, fmt.Errorf("can not fetch file: %s", err.Error())
	}

	bundles := newBundles()
	for _, info := range productsInfo {
		productURLs, err := getURL(info)
		if err != nil {
			return nil, fmt.Errorf("can not fetch file: %s", err.Error())
		}

		for _, url := range productURLs {
			bundleName, err := getProductName(url)
			if err != nil {
				return nil, fmt.Errorf("can not fetch file: %s", err.Error())
			}

			bundle, err := newBundle(bundleName, url)
			if err != nil {
				return nil, fmt.Errorf("can not fetch file: %s", err.Error())
			}

			if bundle.Format == "blockmap" {
				continue
			}

			bundles.Add(bundle)
		}
	}

	files := newFiles()
	for _, app := range apps.Values() {
		appBundle, err := bundles.GetAppBundle(app)
		if err != nil {
			return nil, fmt.Errorf("can not fetch file: %s", err.Error())
		}

		file := newFile(appBundle)

		for _, dep := range app.Dependencies() {
			depBundle, err := bundles.GetDependency(dep)
			if err != nil {
				return nil, fmt.Errorf("can not fetch file: %s", err.Error())
			}
			file.Add(depBundle)
		}
		files.Add(file)
	}

	appFile, err := files.Get(version, ms.Arch)
	if err != nil {
		return nil, fmt.Errorf("can not fetch file: %s", err.Error())
	}

	bundlesToDownload := appFile.Bundles()
	log.Info("Product files fetched")

	log.Debug("Download product files...")
	var result []ms.FileInfo
	for _, data := range bundlesToDownload {
		fullPath := utils.Join(
			destinationPath,
			fmt.Sprintf("%s-%s.%s", data.Name, data.Version.String(), data.Format),
		)

		file, err := os.OpenFile(fullPath, os.O_CREATE, 0660)
		if err != nil {
			return nil, fmt.Errorf("can not download file: can not open file \"%s\": %s", fullPath, err.Error())
		}

		_, err = client.R().SetOutput(file).Get(data.URL)
		_ = file.Close()
		if err != nil {
			return nil, fmt.Errorf("can not download file: GET %s: %s", data.URL, err.Error())
		}

		result = append(result, ms.FileInfo{Path: fullPath, Name: data.Name, Version: data.Version})
	}
	log.Info("Product files downloaded")

	return result, nil
}
