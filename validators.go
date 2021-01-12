package groudon

import (
	"regexp"
)

const (
	UUID_PATTERN       = `[0-9a-f]{8}-[0-9a-f]{4}-[0-5][0-9a-f]{3}-[089ab][0-9a-f]{3}-[0-9a-f]{12}`
	UUID_ONLY_PATTERN  = `^` + UUID_PATTERN + `$`
	EMAIL_PATTERN      = "[^@]+@[^@]+"
	EMAIL_ONLY_PATTERN = `^` + EMAIL_PATTERN + `$`
)

var (
	uuid_regex       *regexp.Regexp = regexp.MustCompile(UUID_PATTERN)
	uuid_only_regex  *regexp.Regexp = regexp.MustCompile(UUID_ONLY_PATTERN)
	email_regex      *regexp.Regexp = regexp.MustCompile(EMAIL_PATTERN)
	email_only_regex *regexp.Regexp = regexp.MustCompile(EMAIL_ONLY_PATTERN)
)

func ValidUUID(it interface{}) (ok bool, _ error) {
	var id string
	if id, ok = it.(string); !ok {
		return
	}

	ok = uuid_only_regex.MatchString(id)
	return
}

func ValidEmail(it interface{}) (ok bool, _ error) {
	var email string
	if email, ok = it.(string); !ok {
		return
	}

	ok = email_only_regex.MatchString(email)
	return
}

func ValidString(it interface{}) (ok bool, _ error) {
	_, ok = it.(string)
	return
}

func ValidStringSlice(it interface{}) (ok bool, _ error) {
	if ok = it == nil; ok {
		return
	}

	if _, ok = it.([]string); ok {
		return
	}

	var them []interface{}
	if them, ok = it.([]interface{}); !ok {
		return
	}

	for _, it = range them {
		if _, ok = it.(string); !ok {
			break
		}
	}
	return
}

func ValidNumber(it interface{}) (ok bool, _ error) {
	if _, ok = it.(float64); ok {
		return
	}

	_, ok = it.(int)
	return
}

func ValidBool(it interface{}) (ok bool, _ error) {
	_, ok = it.(bool)
	return
}

func OptionalString(it interface{}) (ok bool, _ error) {
	if ok = it == nil; ok {
		return
	}

	_, ok = it.(string)
	return
}

func OptionalNumber(it interface{}) (ok bool, _ error) {
	if ok = it == nil; ok {
		return
	}

	if _, ok = it.(float64); ok {
		return
	}

	_, ok = it.(int)
	return
}

func OptionalBool(it interface{}) (ok bool, _ error) {
	if ok = it == nil; ok {
		return
	}

	_, ok = it.(bool)
	return
}
