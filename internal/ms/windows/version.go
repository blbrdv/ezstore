package windows

import (
	"github.com/gonutz/w32/v2"
)

type WindowsVersion int

const (
	W8 WindowsVersion = iota + 8
	W8d1
	W10
	W11
	Unsupported WindowsVersion = 0
)

var Version = GetVersion()

func GetVersion() WindowsVersion {
	v := w32.RtlGetVersion()

	switch v.MajorVersion {
	case 10:
		if v.MinorVersion >= 22000 {
			return W11
		} else {
			return W10
		}
	case 6:
		if v.MinorVersion >= 9600 {
			return W8d1
		} else {
			return W8
		}
	default:
		return Unsupported
	}
}
