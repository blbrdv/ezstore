package windows

import (
	"github.com/gonutz/w32/v2"
)

type WindowsVersion int

const (
	Unsupported WindowsVersion = 0
	W8                         = 8
	W8d1                       = 9
	W10                        = 10
	W11                        = 11
)

var Version = GetVersion()

func GetVersion() WindowsVersion {
	v := w32.RtlGetVersion()

	if v.MajorVersion == 10 {
		if v.MinorVersion >= 22000 {
			return W11
		} else {
			return W10
		}
	} else if v.MajorVersion == 6 {
		if v.MinorVersion >= 9600 {
			return W8d1
		} else {
			return W8
		}
	} else {
		return Unsupported
	}
}
