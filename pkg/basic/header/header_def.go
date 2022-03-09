package header

// header define
const (
	HeaderContentType  = "Content-Type"
	HeaderUserID       = "User-Id"
	HeaderUserName     = "User-Name"
	HeaderDepartmentID = "Department-Id"
	HeaderRole         = "Role"
	HeaderRequestID    = "Request-Id"
	HeaderXRequestID   = "X-Request-Id"
	HeaderAccessToken  = "Access-Token"
	HeaderAccessKeyID  = "Access-Key-Id"
	HeaderXTimezone    = "X-Timezone"
	HeaderTimezone     = "Timezone"
)

// header prefix
const (
	HeaderPrefixAccessKeyID = "key_" // key_rsYO687sdJ=
)

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
)
