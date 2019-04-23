package environ

import (
    "github.com/hjjg200/sprout/util"
)

const (
    AppName = "sprout"
    AppVersion = "pre-alpha"

    ErrorPageTemplatePath = "template/error_page.html"
    IndexPageTemplatePath = "template/index_page.html"
)

var (
    Debug = false
    Logger = util.NewLogger()
)