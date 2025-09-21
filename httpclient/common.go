package httpclient

type ContentType string

const (
	HeadContentType = "content-type"
)

const (
	ContentTypeJson ContentType = "application/json"
	ContentTypeForm             = "application/x-www-form-urlencoded"
	ContentTypeText             = "text/plain"
	ContentTypeXml              = "application/xml"
	ContentTypeHtml             = "text/html"
)

func getContentType(contentType ContentType, charset ...string) string {
	if len(charset) > 0 {
		return string(contentType) + "; charset=" + charset[0]
	}
	return string(contentType)
}
