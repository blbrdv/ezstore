package store

import (
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/blbrdv/ezstore/internal/ms"
	"strings"
)

const displaycatalogURL = "https://displaycatalog.mp.microsoft.com/v7.0/products"

func canRedeem(skuAvailability *jsonquery.Node) (bool, error) {
	availabilities := jsonquery.FindOne(skuAvailability, "Availabilities")
	if availabilities == nil {
		return false, fmt.Errorf("can not get app info: can not get sku availabilities")
	}

	for _, availability := range availabilities.ChildNodes() {
		actions := jsonquery.FindOne(availability, "Actions")
		if actions == nil {
			return false, fmt.Errorf("can not get app info: can not get actions from sku availability")
		}

		for _, action := range actions.ChildNodes() {
			value := strings.ToLower(fmt.Sprintf("%v", action.Value()))
			if value == "redeem" {
				return true, nil
			}
		}
	}

	return false, nil
}

func getAppInfo(id string, locale *ms.Locale) (*apps, string, error) {
	url := fmt.Sprintf(
		"%s/%s?market=%s&languages=%s,%s,neutral",
		displaycatalogURL,
		id,
		locale.Country,
		locale.String(),
		locale.Language,
	)
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("can not get app info: GET %s: %s", url, err.Error())
	}
	if resp.StatusCode == 404 {
		return nil, "", fmt.Errorf(`product with id "%s" and locale "%s" not found`, id, locale.String())
	}
	if resp.IsErrorState() {
		return nil, "", fmt.Errorf("can not get app info: GET %s: server returns error: %s", url, resp.Status)
	}

	data, err := jsonquery.Parse(strings.NewReader(resp.String()))
	if err != nil {
		return nil, "", fmt.Errorf("can not get app info: can not parse result: %s", err.Error())
	}

	skusAvailabilities := jsonquery.FindOne(data, "Product/DisplaySkuAvailabilities")
	if skusAvailabilities == nil {
		return nil, "", fmt.Errorf("can not get app info: can not find availabilities data in response body")
	}
	if len(skusAvailabilities.ChildNodes()) == 0 {
		return nil, "", fmt.Errorf("can not get app info: no availabilities for this app")
	}

	var packages []*jsonquery.Node
	for _, availability := range skusAvailabilities.ChildNodes() {
		redeemable, err := canRedeem(availability)
		if err != nil {
			return nil, "", err
		}

		if redeemable {
			skuPackages := jsonquery.FindOne(availability, "Sku/Properties/Packages")
			if skuPackages == nil {
				return nil, "", fmt.Errorf("can not get app info: can not get sku packages")
			}

			for _, skuPackage := range skuPackages.ChildNodes() {
				platformDeps := jsonquery.FindOne(skuPackage, "PlatformDependencies").ChildNodes()
				for _, platform := range platformDeps {
					identity := jsonquery.FindOne(platform, "PlatformName")
					if identity != nil {
						value := strings.ToLower(fmt.Sprintf("%v", identity.Value()))
						if value == "windows.desktop" || value == "windows.universal" {
							packages = append(packages, skuPackage)
						}
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
		fullName := jsonquery.FindOne(pkg, "PackageFullName")
		if fullName == nil {
			return nil, "", fmt.Errorf("can not get app info: can not get full name")
		}

		app, err := newApp(fmt.Sprintf("%v", fullName.Value()))
		if err != nil {
			return nil, "", err
		}

		deps := jsonquery.FindOne(pkg, "FrameworkDependencies").ChildNodes()
		for _, dep := range deps {
			name := jsonquery.FindOne(dep, "PackageIdentity")
			app.Add(fmt.Sprintf("%v", name.Value()))
		}

		apps.Add(app)

		if wuid == "" {
			value := fmt.Sprintf("%v", jsonquery.FindOne(pkg, "FulfillmentData/WuCategoryId").Value())
			if value == "" {
				return nil, "", fmt.Errorf("can not get app info: can not find WUID in fulfillment data")
			}

			wuid = value
		}
	}

	return apps, wuid, nil
}
