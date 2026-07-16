package httpclient

import "mime"

type ContentType string

const (
	HeaderContentType = "Content-Type"
)

const (
	ContentTypeJSON ContentType = "application/json"
	ContentTypeForm ContentType = "application/x-www-form-urlencoded"
	ContentTypeText ContentType = "text/plain"
	ContentTypeXML  ContentType = "application/xml"
	ContentTypeHTML ContentType = "text/html"
)

func getContentType(contentType ContentType, charset ...string) string {
	if len(charset) == 0 || charset[0] == "" {
		return string(contentType)
	}
	return mime.FormatMediaType(string(contentType), map[string]string{"charset": charset[0]})
}
