package validation

import "regexp"

func GenericEmailFormat(email string) bool {
	pattern := `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
