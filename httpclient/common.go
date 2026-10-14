package httpclient

import "mime"

// ContentType identifies a MIME content type.
type ContentType string

const (
	// HeaderContentType is the Content-Type header name.
	HeaderContentType = "Content-Type"
)

const (
	// ContentTypeJSON identifies JSON content.
	ContentTypeJSON ContentType = "application/json"
	// ContentTypeForm identifies URL-encoded form content.
	ContentTypeForm ContentType = "application/x-www-form-urlencoded"
	// ContentTypeText identifies plain-text content.
	ContentTypeText ContentType = "text/plain"
	// ContentTypeXML identifies XML content.
	ContentTypeXML ContentType = "application/xml"
	// ContentTypeHTML identifies HTML content.
	ContentTypeHTML ContentType = "text/html"
)

func getContentType(contentType ContentType, charset ...string) string {
	if len(charset) == 0 || charset[0] == "" {
		return string(contentType)
	}
	return mime.FormatMediaType(string(contentType), map[string]string{"charset": charset[0]})
}
