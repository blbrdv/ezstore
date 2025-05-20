package ms

import (
	"fmt"
	"regexp"
)

const pattern = `^(?P<lang>[a-z][a-z][a-z]?)(?:-(?:\d\d\d|[A-Z][a-z]+))?(?:[_-](?P<country>[A-Z][A-Z]))?(?:_\w+|-(?:[a-z]+|\d+))?$`

// Locale represents ISO 639-1 language tag and optional ISO 3166-1 country code.
type Locale struct {
	Language string
	Country  string
}

// String returns MS Store compatible locale literal.
func (l *Locale) String() string {
	if l.Country == "" {
		return l.Language
	}

	return fmt.Sprintf("%s-%s", l.Language, l.Country)
}

func (l *Locale) Equal(other *Locale) bool {
	return l.Country == other.Country && l.Language == other.Language
}

// NewLocale returns [Locale] from RFC 5646 input string or error if invalid format.
// See Appendix A for examples.
func NewLocale(input string) (*Locale, error) {
	localeRegexp := regexp.MustCompile(pattern)
	matches := localeRegexp.FindStringSubmatch(input)
	if matches == nil {
		return nil, fmt.Errorf("\"%s\" is not a valid locale", input)
	}

	return &Locale{
		Language: matches[localeRegexp.SubexpIndex("lang")],
		Country:  matches[localeRegexp.SubexpIndex("country")],
	}, nil
}
