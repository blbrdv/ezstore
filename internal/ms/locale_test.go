package ms_test

import (
	"fmt"
	"testing"

	. "github.com/blbrdv/ezstore/internal/ms"
)

var langOnlyData = []struct {
	Name string
	Raw  string
}{
	// English locale
	{"TestLocaleLang", "en"},
	// Pedi locale
	{"TestLocaleLangLong", "nso"},
}

func TestLocaleLangOnly(t *testing.T) {
	for _, data := range langOnlyData {
		t.Run(data.Name, func(t *testing.T) {
			expected := Locale{Language: data.Raw}
			actual, err := NewLocale(data.Raw)

			if err != nil {
				t.Fatalf(`Can not parse locale: %s`, err.Error())
			}

			expectedStr := expected.String()
			actualStr := actual.String()

			if actualStr != expectedStr {
				t.Fatalf(`Incorrect Locale, expected: "%s", actual: "%s"`, expectedStr, actualStr)
			}
		})
	}
}

var langOnlyWithNoiseData = []struct {
	Name     string
	Raw      string
	Language string
}{
	// English (Caribbean) locale
	{"TestLocaleLangWithRegionCode", "en-029", "en"},
	// Tachelhit (Latin) locale
	{"TestLocaleLangWithScript", "shi-Latn", "shi"},
}

func TestLocaleLangOnlyWithNoise(t *testing.T) {
	for _, data := range langOnlyWithNoiseData {
		t.Run(data.Name, func(t *testing.T) {
			expected := Locale{Language: data.Language}
			actual, err := NewLocale(data.Raw)

			if err != nil {
				t.Fatalf(`Can not parse locale: %s`, err.Error())
			}

			expectedStr := expected.String()
			actualStr := actual.String()

			if actualStr != expectedStr {
				t.Fatalf(`Incorrect Locale, expected: "%s", actual: "%s"`, expectedStr, actualStr)
			}
			if actualStr != data.Language {
				t.Fatalf(`Incorrect Locale, expected: "%s", actual: "%s"`, data.Language, actualStr)
			}
		})
	}
}

var langWithCountryData = []struct {
	Name     string
	Raw      string
	Language string
	Country  string
}{
	// English (United States) locale
	{"TestLocaleLangWithCountry", "en-US", "en", "US"},
	// Yangben (Cameroon) locale
	{"TestLocaleLongLangWithCountry", "yav_CM", "yav", "CM"},
	// German (Germany) locale with phone book order
	{"TestLocaleLangWithCountryAndAlternateSort", "de-DE_phoneb", "de", "DE"},
	// Valencian (Spain) locale
	{"TestLocaleLangWithCountryAndVariant", "ca_ES-valencia", "ca", "ES"},
	// German (Switzerland) locale using orthography variant
	{"TestLocaleLangWithCountryAndNumericVariant", "de-CH-1901", "de", "CH"},
	// Bosnian (Cyrillic, Bosnia and Herzegovina) locale
	{"TestLocaleLangWithScriptAndCountry", "bs-Cyrl-BA", "bs", "BA"},
	// Vai (Latin, Liberia) locale
	{"TestLocaleLongLangWithScriptAndCountry", "vai-Latn_LR", "vai", "LR"},
}

func TestLocaleLangWithCountry(t *testing.T) {
	for _, data := range langWithCountryData {
		t.Run(data.Name, func(t *testing.T) {
			expected := Locale{Language: data.Language, Country: data.Country}
			actual, err := NewLocale(data.Raw)

			if err != nil {
				t.Fatalf(`Can not parse locale: %s`, err.Error())
			}

			expectedStr := expected.String()
			actualStr := actual.String()

			if actualStr != expectedStr {
				t.Fatalf(`Incorrect Locale, expected: "%s", actual: "%s"`, expectedStr, actualStr)
			}
			if actual.Language != data.Language {
				t.Fatalf(`Incorrect Locale, expected language: "%s", actual: "%s"`, data.Language, actual.Language)
			}
			if actual.Country != data.Country {
				t.Fatalf(`Incorrect Locale, expected country: "%s", actual: "%s"`, data.Country, actual.Country)
			}
		})
	}
}

var incorrectLocaleData = []struct {
	Name  string
	Value string
}{
	{"TestEmptyString", ""},
	{"TestInvalidFormat", "foo bar 42"},
	{"TestInvalidLocale", "bababooey"},
	{"TestInvalidSeparator", "en/US"},
	{"TestInvalidOrder", "CM_yav"},
	{"TestInvalidVariantSeparator", "ca_ES/valencia"},
	{"TestInvalidScriptFormat", "bs-xXxCyrlx1337Xx-BA"},
	{"TestInvalidCountryAndNumericVariantOrder", "de-1901-CH"},
}

func TestInvalidLocale(t *testing.T) {
	for _, data := range incorrectLocaleData {
		t.Run(data.Name, func(t *testing.T) {
			expected := fmt.Sprintf("\"%s\" is not a valid locale", data.Value)
			result, err := NewLocale(data.Value)

			if err == nil {
				t.Fatalf(`Function must return error "%s", but return result "%s"`, expected, result.String())
			}
			if err.Error() != expected {
				t.Fatalf(`Incorrect error message, expected "%s", actual "%s"`, expected, err.Error())
			}
		})
	}
}
