package rest

import "strings"

// ConcatURL handle slashes between segments automatically.
func ConcatURL(segments ...string) string {
	if len(segments) <= 0 {
		return ""
	}
	if len(segments) == 1 {
		return segments[0]
	}
	urls := make([]string, len(segments))
	for i, segment := range segments {
		switch i {
		case 0:
			urls[i] = strings.TrimRight(segment, "/")
		case len(segment) - 1:
			urls[i] = strings.TrimLeft(segment, "/")
		default:
			urls[i] = strings.Trim(segment, "/")
		}
	}
	return strings.Join(urls, "/")
}
