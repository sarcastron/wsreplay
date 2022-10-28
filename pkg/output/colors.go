package output

import "github.com/fatih/color"

var Info func(a ...interface{}) string = color.New(color.FgGreen).SprintFunc()
var Notice func(a ...interface{}) string = color.New(color.FgMagenta).SprintFunc()
var Warning func(a ...interface{}) string = color.New(color.FgYellow).SprintFunc()
var Danger func(a ...interface{}) string = color.New(color.FgRed).SprintFunc()
