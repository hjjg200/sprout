package sprout

import (
    "errors"
)

const (
    default_localizer_threshold = 5

    cookie_locale = "locale"
)

var (
    defaultWhitelistedExtensions = []string{
        ".css", ".js", ".jpg", ".jpeg", ".png", ".gif", ".ico", ".icn",
    }
    template_extensions = []string{
        ".html", ".htm",
    }
)

var (
    ErrInvalidLocale    = errors.New( "sprout: the given locale is invalid" )
    ErrNoAvailableCache = errors.New( "sprout: no available cache" )
)