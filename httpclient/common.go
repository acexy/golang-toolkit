package httpclient

type ContentType string

const (
	HeaderContentType = "content-type"
)

const (
	ContentTypeJSON ContentType = "application/json"
	ContentTypeForm ContentType = "application/x-www-form-urlencoded"
	ContentTypeText ContentType = "text/plain"
	ContentTypeXML  ContentType = "application/xml"
	ContentTypeHTML ContentType = "text/html"
)

func getContentType(contentType ContentType, charset ...string) string {
	if len(charset) > 0 {
		return string(contentType) + "; charset=" + charset[0]
	}
	return string(contentType)
}
