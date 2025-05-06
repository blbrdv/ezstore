package store

import (
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/blbrdv/ezstore/internal/ms"
	"strings"
)

const displaycatalogURL = "https://displaycatalog.mp.microsoft.com/v7.0/products"

func getAppInfo(id string, locale *ms.Locale) (*bundleInfo, string, error) {
	url := fmt.Sprintf("%s/%s", displaycatalogURL, id)
	resp, err := client.R().
		SetQueryParam("market", locale.Country).
		SetQueryParam("languages", fmt.Sprintf("%s,%s,neutral", locale.String(), locale.Language)).
		Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("can not get app info: GET %s: %s", url, err.Error())
	}
	if resp.StatusCode == 404 {
		return nil, "", fmt.Errorf(`product with id "%s" and locale "%s" not found`, id, locale.String())
	}
	if resp.IsErrorState() {
		return nil, "", fmt.Errorf("can not get app info: GET %s: server returns error: %s", url, resp.ErrorResult())
	}

	data, err := jsonquery.Parse(strings.NewReader(resp.String()))
	if err != nil {
		return nil, "", fmt.Errorf("can not get app info: can not parse result: %s", err.Error())
	}

	fulfillmentData := jsonquery.FindOne(data, "Product/DisplaySkuAvailabilities/*[1]/Sku/Properties/FulfillmentData")
	if fulfillmentData == nil {
		return nil, "", fmt.Errorf("can not get app info: can not find fulfillment data in response body")
	}

	wuid := fmt.Sprintf("%v", jsonquery.FindOne(fulfillmentData, "WuCategoryId").Value())
	if wuid == "" {
		return nil, "", fmt.Errorf("can not get app info: can not find WUID in fulfillment data")
	}

	packageName := fmt.Sprintf("%v", jsonquery.FindOne(fulfillmentData, "PackageFamilyName").Value())
	if packageName == "" {
		return nil, "", fmt.Errorf("can not get app info: can not find package name in fulfillment data")
	}

	info, err := newBundleInfo(packageName)
	if err != nil {
		return nil, "", fmt.Errorf("can not get app info: can not get bundle info: %s", err.Error())
	}

	return info, wuid, nil
}
