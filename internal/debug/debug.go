// +build debug

package debug

import "fmt"

func Debug(format string, a ...interface{}) {
	fmt.Println("[PKGER]", fmt.Sprintf(format, a...))
}
