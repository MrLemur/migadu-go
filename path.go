package migadu

import "net/url"

func escapePathSegment(value string) string {
	return url.PathEscape(value)
}
