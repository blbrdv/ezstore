package types

import (
	"fmt"
	"regexp"
)

const pattern = `^(?P<lang>[a-z][a-z][a-z]?)(?:-(?:\d\d\d|[A-Z][a-z]+))?(?:-(?P<country>[A-Z][A-Z]))?(?:_\w+|-(?:[a-z]+|\d+))?$`

// Locale represents ISO 639-1 language tag and optional ISO 3166-1 country code.
type Locale struct {
	Language string
	Country  string
}

// String returns MS Store compatible locale literal
func (l Locale) String() string {
	if l.Country == "" {
		return l.Language
	}

	return fmt.Sprintf("%s_%s", l.Language, l.Country)
}

// Parse returns Locale from RFC 5646 input string or error if invalid format.
// See Appendix A for examples.
func Parse(input string) (Locale, error) {
	regex := regexp.MustCompile(pattern)
	match := regex.FindStringSubmatch(input)
	if match == nil {
		return Locale{}, fmt.Errorf("%s is not a valid locale", input)
	}

	return Locale{
		Language: match[regex.SubexpIndex("lang")],
		Country:  match[regex.SubexpIndex("country")],
	}, nil
}
