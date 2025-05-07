package store

import (
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
)

func TestFE3FileUrl(t *testing.T) {
	expectedTicketType := "aaaaa"
	expectedID := "bbbbb"
	expectedRevisionNumber := "ccccc"

	xml := strings.NewReader(fe3FileURL(expectedTicketType, expectedID, expectedRevisionNumber))

	doc, err := xmlquery.Parse(xml)
	if err != nil {
		t.Fatalf(`Can not parse XML file: %s`, err.Error())
	}

	root := xmlquery.FindOne(doc, "//s:Envelope")

	actualTicketType := root.
		SelectElement("//s:Header/o:Security/wuws:WindowsUpdateTicketsToken/TicketType").
		InnerText()
	if actualTicketType != expectedTicketType {
		t.Fatalf(`Incorrect TycketType, expected: "%s", actual: "%s"`, expectedTicketType, actualTicketType)
	}

	identity := root.
		SelectElement("//s:Body/GetExtendedUpdateInfo2/updateIDs/UpdateIdentity")

	actualID := identity.
		SelectElement("//UpdateID").
		InnerText()
	if actualID != expectedID {
		t.Fatalf(`Incorrect TycketType, expected: "%s", actual: "%s"`, expectedID, actualID)
	}

	actualRevisionNumber := identity.
		SelectElement("//RevisionNumber").
		InnerText()
	if actualRevisionNumber != expectedRevisionNumber {
		t.Fatalf(`Incorrect TycketType, expected: "%s", actual: "%s"`, expectedRevisionNumber, actualRevisionNumber)
	}
}

func TestWUIDRequest(t *testing.T) {
	expectedTycketType := "aaaaa"
	expectedCookie := "bbbbb"
	expectedCategoryIdentifier := "ccccc"

	xml := strings.NewReader(wuidRequest(expectedTycketType, expectedCookie, expectedCategoryIdentifier))

	doc, err := xmlquery.Parse(xml)
	if err != nil {
		t.Fatalf(`Can not parse XML file: %s`, err.Error())
	}

	root := xmlquery.FindOne(doc, "//s:Envelope")

	actualTicketType := root.
		SelectElement("//s:Header/o:Security/wuws:WindowsUpdateTicketsToken/TicketType").
		InnerText()
	if actualTicketType != expectedTycketType {
		t.Fatalf(`Incorrect TycketType, expected: "%s", actual: "%s"`, expectedTycketType, actualTicketType)
	}

	syncUpdates := root.
		SelectElement("//s:Body/SyncUpdates")

	actualCookie := syncUpdates.
		SelectElement("//cookie/EncryptedData").
		InnerText()
	if actualCookie != expectedCookie {
		t.Fatalf(`Incorrect TycketType, expected: "%s", actual: "%s"`, expectedCookie, actualCookie)
	}

	actualCategoryIdentifier := syncUpdates.
		SelectElement("//parameters/FilterAppCategoryIds/CategoryIdentifier/Id").
		InnerText()
	if actualCategoryIdentifier != expectedCategoryIdentifier {
		t.Fatalf(`Incorrect TycketType, expected: "%s", actual: "%s"`, expectedCategoryIdentifier, actualCategoryIdentifier)
	}
}
