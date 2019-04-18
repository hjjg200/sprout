package environ

import (
    "../util"
)

const (
    AppName = "sprout"
    AppVersion = "pre-alpha"

    ErrorPageTemplatePath = "template/error_page.html"
)

var Logger = util.NewLogger()