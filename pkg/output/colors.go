package output

import (
	"fmt"

	"github.com/fatih/color"
)

var Info func(a ...interface{}) string = color.New(color.FgGreen).SprintFunc()
var Notice func(a ...interface{}) string = color.New(color.FgMagenta).SprintFunc()
var Warning func(a ...interface{}) string = color.New(color.FgYellow).SprintFunc()
var Danger func(a ...interface{}) string = color.New(color.FgRed).SprintFunc()

func ErrorMsg(err error) {
	fmt.Printf("%s %s\n", Danger("Error:"), err)
}
