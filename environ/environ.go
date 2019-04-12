package environ

import (
    "../util"
)

const (
    AppName = "sprout"
    AppVersion = "pre-alpha"
)

var Logger = util.NewLogger()
var Debug = false