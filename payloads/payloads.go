package payloads

import (
	_ "embed"
	"fmt"
)

//go:embed FE3FileUrl.xml
var fe3FileUrl string

//go:embed GetCookie.xml
var getCookie string

//go:embed WUIDRequest.xml
var wuidRequiest string

func FE3FileUrl(ticketType string, id string, revisionNumber string) string {
	return fmt.Sprintf(fe3FileUrl, ticketType, id, revisionNumber)
}

func GetCookie() string {
	return getCookie
}

func WUIDRequest(ticketType string, cookie string, categoryIdentifier string) string {
	return fmt.Sprintf(wuidRequiest, ticketType, cookie, categoryIdentifier)
}
