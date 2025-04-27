package msstore

import (
	_ "embed"
	"fmt"
)

//go:embed payloads/FE3FileUrl.xml
var fileUrlPayload string

//go:embed payloads/GetCookie.xml
var getCookiePayload string

//go:embed payloads/WUIDRequest.xml
var wuidRequestPayload string

func fe3FileUrl(ticketType string, id string, revisionNumber string) string {
	return fmt.Sprintf(fileUrlPayload, ticketType, id, revisionNumber)
}

func wuidRequest(ticketType string, cookie string, categoryIdentifier string) string {
	return fmt.Sprintf(wuidRequestPayload, ticketType, cookie, categoryIdentifier)
}
