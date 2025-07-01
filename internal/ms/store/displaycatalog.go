package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"net/http"
	"strings"
)

const displaycatalogURL = "https://displaycatalog.mp.microsoft.com/v7.0/products"

type framework struct {
	Name string `json:"PackageIdentity"`
	Min  uint64 `json:"MinVersion"`
	Max  uint64 `json:"MaxTested"`
}

type platform struct {
	Name string `json:"PlatformName"`
}

type fulfillmentData struct {
	ID string `json:"WuCategoryId"`
}

type jsonPkg struct {
	Architectures         []string        `json:"Architectures"`
	Name                  string          `json:"PackageFullName"`
	FrameworkDependencies []framework     `json:"FrameworkDependencies"`
	PlatformDependencies  []platform      `json:"PlatformDependencies"`
	FulfillmentData       fulfillmentData `json:"FulfillmentData"`
}

type properties struct {
	Packages []jsonPkg `json:"Packages"`
}

type sku struct {
	Properties properties `json:"Properties"`
}

type availability struct {
	Actions []string `json:"Actions"`
}

type skuAvailability struct {
	Sku            sku            `json:"Sku"`
	Availabilities []availability `json:"Availabilities"`
}

type product struct {
	SkuAvailabilities []skuAvailability `json:"DisplaySkuAvailabilities"`
}

type appInfo struct {
	Product product `json:"Product"`
}

func canRedeem(skuAvailability skuAvailability) (bool, error) {
	for _, availability := range skuAvailability.Availabilities {
		for _, action := range availability.Actions {
			if strings.ToLower(action) == "redeem" {
				return true, nil
			}
		}
	}

	return false, nil
}

func getAppInfo(id string, locale *ms.Locale) (*apps, string, error) {
	var info appInfo
	url := fmt.Sprintf(
		"%s/%s?market=%s&languages=%s,%s,neutral",
		displaycatalogURL,
		id,
		locale.Country,
		locale.String(),
		locale.Language,
	)
	resp, err := client.R().SetSuccessResult(&info).Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("can not get app info: GET %s: %s", url, err.Error())
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, "", fmt.Errorf(`product with id "%s" and locale "%s" not found`, id, locale.String())
	}
	if resp.IsErrorState() {
		return nil, "", fmt.Errorf("can not get app info: GET %s: server returns error: %s", url, resp.Status)
	}

	if len(info.Product.SkuAvailabilities) == 0 {
		return nil, "", fmt.Errorf("can not get app info: no availabilities for this app")
	}

	var packages []jsonPkg
	for _, availability := range info.Product.SkuAvailabilities {
		redeemable, err2 := canRedeem(availability)
		if err2 != nil {
			return nil, "", err2
		}

		if redeemable {
			for _, skuPackage := range availability.Sku.Properties.Packages {
				for _, platform := range skuPackage.PlatformDependencies {
					dep := strings.ToLower(platform.Name)
					if dep == "windows.desktop" || dep == "windows.universal" {
						packages = append(packages, skuPackage)
					}
				}
			}
		}
	}

	if len(packages) == 0 {
		return nil, "", fmt.Errorf("can not get app info: no available packages found")
	}

	apps := newApps()
	wuid := ""
	for _, pkg := range packages {
		app, err2 := newApp(pkg.Name, pkg.Architectures[0])
		if err2 != nil {
			return nil, "", err2
		}

		for _, dep := range pkg.FrameworkDependencies {
			var minVersion *ms.Version = nil
			var maxVersion *ms.Version = nil
			if dep.Min != 0 {
				minVersion = ms.NewVersionFromNumber(dep.Min)
			}
			if dep.Max != 0 {
				maxVersion = ms.NewVersionFromNumber(dep.Max)
			}

			app.Add(dep.Name, minVersion, maxVersion)
		}

		apps.Add(app)

		if wuid == "" {
			wuid = pkg.FulfillmentData.ID
		}
	}

	return apps, wuid, nil
}
