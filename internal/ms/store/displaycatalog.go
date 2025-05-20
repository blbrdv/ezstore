package store

import (
	"fmt"
	"github.com/blbrdv/ezstore/internal/ms"
	"strings"
)

const displaycatalogURL = "https://displaycatalog.mp.microsoft.com/v7.0/products"

type framework struct {
	Name string `json:"PackageIdentity"`
}

type platform struct {
	Name string `json:"PlatformName"`
}

type fulfillmentData struct {
	ID string `json:"WuCategoryId"`
}

type jsonPkg struct {
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
	var appInfo appInfo
	url := fmt.Sprintf(
		"%s/%s?market=%s&languages=%s,%s,neutral",
		displaycatalogURL,
		id,
		locale.Country,
		locale.String(),
		locale.Language,
	)
	resp, err := client.R().SetSuccessResult(&appInfo).Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("can not get app info: GET %s: %s", url, err.Error())
	}
	if resp.StatusCode == 404 {
		return nil, "", fmt.Errorf(`product with id "%s" and locale "%s" not found`, id, locale.String())
	}
	if resp.IsErrorState() {
		return nil, "", fmt.Errorf("can not get app info: GET %s: server returns error: %s", url, resp.Status)
	}

	if len(appInfo.Product.SkuAvailabilities) == 0 {
		return nil, "", fmt.Errorf("can not get app info: no availabilities for this app")
	}

	var packages []jsonPkg
	for _, availability := range appInfo.Product.SkuAvailabilities {
		redeemable, err := canRedeem(availability)
		if err != nil {
			return nil, "", err
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
		app, err := newApp(pkg.Name)
		if err != nil {
			return nil, "", err
		}

		for _, dep := range pkg.FrameworkDependencies {
			app.Add(dep.Name)
		}

		apps.Add(app)

		if wuid == "" {
			wuid = pkg.FulfillmentData.ID
		}
	}

	return apps, wuid, nil
}
