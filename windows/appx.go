package windows

import (
	"fmt"
)

func Install(path string) {
	execute(fmt.Sprintf(`Add-AppxPackage -Path "%s"`, path))
}
