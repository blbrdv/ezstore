package store

import (
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/blbrdv/ezstore/internal/ms"
	"strings"
)

func getWUID(id string, locale *ms.Locale) (*bundleInfo, string, error) {
	url := fmt.Sprintf(
		"%s%s?market=%s&languages=%s,%s,neutral",
		"https://displaycatalog.mp.microsoft.com/v7.0/products/",
		id,
		locale.Country,
		locale.String(),
		locale.Language,
	)

	resp, err := client.R().Get(url)
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
