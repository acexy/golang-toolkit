package http

type ContentType string
type HttpMethod uint

// header key
const (
	HeadContentType = "content-type"
)

// contentType
const (
	ContentTypeJson ContentType = "application/json"
)

const (
	HttpMethodGet = iota
	HttpMethodPost
	HttpMethodPut
	HttpMethodDelete
	HttpMethodHead
	HttpMethodOptions
	HttpMethodPatch
)
