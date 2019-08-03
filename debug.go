// +build debug

package pkger

func Debug(format string, a ...interface{}) {
		fmt.Println("[PKGER] ", fmt.Sprintf(format, a...)
}
