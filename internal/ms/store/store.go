package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/log"
	"github.com/blbrdv/ezstore/internal/ms"
	"github.com/blbrdv/ezstore/internal/ms/windows"
	"github.com/blbrdv/ezstore/internal/utils"
	"net/http"
	"os"
	"regexp"
)

var fileNameRegexp = regexp.MustCompile(`filename=(\S+)`)

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
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetching product name failed: HEAD %s: server responded with error: %s", url, res.Status)
	}

	header := res.Header.Get("Content-Disposition")
	if header == "" {
		return "", fmt.Errorf("fetching product name failed: can not get file name: response header \"Content-Disposition\" is empty")
	}

	matches := fileNameRegexp.FindStringSubmatch(header)
	if len(matches) != 2 {
		return "", fmt.Errorf("fetching product name failed: can not get file name: response header \"Content-Disposition\" has invalid format: %s", header)
	}

	return matches[1], nil
}

func downloadFile(destinationPath string, data *bundle) (*ms.FileInfo, error) {
	fullPath := utils.Join(
		destinationPath,
		fmt.Sprintf("%s_%s_%s.%s", data.Name, data.Version.String(), data.Arch, data.Format),
	)

	file := windows.OpenFile(fullPath, os.O_CREATE)

	_, err := client.R().SetOutput(file).Get(data.URL)
	file.Close()
	if err != nil {
		return nil, fmt.Errorf("can not download file: GET %s: %s", data.URL, err.Error())
	}

	return ms.NewFileInfo(fullPath, data.Name, data.Version), nil
}

// Download backage and its dependencies from MS Store by id, version and locale to destination directory
// and returns array of backage and its dependencies paths.
func Download(id string, version *ms.Version, locale *ms.Locale, destinationPath string) (*ms.BundleFileInfo, error) {
	log.Debug("Fetching cookie...")
	cookie, err := getCookie()
	if err != nil {
		return nil, fmt.Errorf("can not fetch cookie: %s", err.Error())
	}
	log.Tracef("Cookie: %s", cookie)
	log.Info("Cookie fetched")

	log.Debug("Fetching product info...")
	apps, wuid, err := getAppInfo(id, locale)
	if err != nil {
		return nil, fmt.Errorf("can not fetch product info: %s", err.Error())
	}
	log.Tracef("Apps: %s", apps.String())
	log.Tracef("WUID: %s", wuid)
	log.Info("Product info fetched")

	log.Debug("Fetching product files...")
	productsInfo, err := getProducts(cookie, wuid)
	if err != nil {
		return nil, fmt.Errorf("can not fetch file: %s", err.Error())
	}
	log.Tracef("Products: %s", PrettyString(productsInfo))

	bundles := newBundles()
	for _, info := range productsInfo {
		productURLs, err2 := getURL(info)
		if err2 != nil {
			return nil, fmt.Errorf("can not fetch file: %s", err2.Error())
		}

		for _, url := range productURLs {
			bundleName, err3 := getProductName(url)
			if err3 != nil {
				return nil, fmt.Errorf("can not fetch file: %s", err3.Error())
			}

			bundle, err3 := newBundle(bundleName, url)
			if err3 != nil {
				return nil, fmt.Errorf("can not fetch file: %s", err3.Error())
			}

			if bundle.Format == "blockmap" {
				continue
			}

			bundles.Add(bundle)
		}
	}
	log.Tracef("Bundles: %s", bundles.String())

	files := newFiles()
	for _, app := range apps.Values() {
		appBundle, err2 := bundles.GetAppBundle(app)
		if err2 != nil {
			return nil, fmt.Errorf("can not fetch file: %s", err2.Error())
		}

		file := newFile(appBundle)

		for _, dep := range app.Dependencies() {
			depBundle, err3 := bundles.GetDependency(dep, app.DepArch)
			if err3 != nil {
				return nil, fmt.Errorf("can not fetch file: %s", err3.Error())
			}
			file.Add(depBundle)
		}
		files.Add(file)
	}
	log.Tracef("Files: %s", files.String())

	appFile, err := files.Get(version, ms.Arch)
	if err != nil {
		return nil, fmt.Errorf("can not fetch file: %s", err.Error())
	}
	log.Tracef("App file: %s", appFile.String())
	log.Info("Product files fetched")

	log.Debug("Download product files...")
	file, err := downloadFile(destinationPath, appFile.GetBundle())
	if err != nil {
		return nil, err
	}

	result := ms.NewBundleFileInfo(file)
	for _, dependency := range appFile.Dependencies() {
		depFile, err2 := downloadFile(destinationPath, dependency)
		if err2 != nil {
			return nil, err2
		}
		result.AddDependency(depFile)
	}
	log.Tracef("Downloaded files: %s", result)
	log.Info("Product files downloaded")

	return result, nil
}
