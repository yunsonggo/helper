package free

import "github.com/golang-module/carbon/v2"

func FormatToRFC3339(data carbon.Carbon) string {
	return data.Layout("Jan-02-2006 03:04:05 PM")
}
