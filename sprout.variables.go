package sprout

/*
 + SPROUT VARIABLES
 */

type sprout_variables struct {
    // GENERAL DEFINITIONS
    app_name string
    version string
    os string
    http_status_messages map[int] string
}
var  static_sprout_variables = &sprout_variables{
    // GENERAL DEFINITIONS
    app_name: "sprout",
    version: "pre-alpha 0.7",
    os: runtime.GOOS,
    http_status_messages: map[int] string{
        // 1xx
        100: "Continue",
        101: "Switching Protocols",
        102: "Processing",
        103: "Early Hints",
        // 2xx
        200: "OK",
        201: "Created",
        202: "Accepted",
        203: "Non-Authoritative Information",
        204: "No Content",
        205: "Reset Content",
        206: "Partial Content",
        207: "Multi-Status",
        208: "Already Reported",
        226: "IM Used",
        // 3xx
        300: "Multiple Choices",
        301: "Moved Permanently",
        302: "Found",
        303: "See Other",
        304: "Not Modified",
        305: "Use Proxy",
        306: "Switch Proxy",
        307: "Temporary Redirect",
        308: "Permanent Redirect",
        //4xx Client Errors
        400: "Bad Request",
        401: "Unathorized",
        402: "Payment Required",
        403: "Forbidden",
        404: "Not Found",
        405: "Method Not Allowed",
        406: "Not Acceptable",
        407: "Proxy Authentication Required",
        408: "Request Timeout",
        409: "Conflict",
        410: "Gone",
        411: "Length Required",
        412: "Precondition Failed",
        413: "Payload Too Large",
        414: "URI Too Long",
        415: "Unsupported Media Type",
        416: "Range Not Satisfiable",
        417: "Expectation Failed",
        418: "I'm a teapot",
        421: "Misdirected Request",
        422: "Unprocessable Entity",
        423: "Locked",
        424: "Failed Dependency",
        426: "Upgrade Required",
        428: "Precondition Required",
        429: "Too Many Requests",
        431: "Request Header Fields Too Large",
        451: "Unavailable For Legal Reasons",
        // 5xx Server Errors
        500: "Internal Server Error",
        501: "Not Implemented",
        502: "Bad Gateway",
        503: "Service Unavailable",
        504: "Gateway Timeout",
        505: "HTTP Version Not Supported",
        506: "Variant Also Negotiates",
        507: "Insufficient Storage",
        508: "Loop Detected",
        510: "Not Extended",
        511: "Network Authentication Required",
    },
}

func SproutVariables() *sprout_variables {}
func( _sprvar *sprout_variables ) AppName() string {}
func( _sprvar *sprout_variables ) Version() string {}
func( _sprvar *sprout_variables ) OS() string {}
func( _sprvar *sprout_variables ) HTTPStatusMessages() func( int ) string {}