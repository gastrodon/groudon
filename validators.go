package groudon

import (
	"regexp"
)

const (
	UUID_PATTERN      = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
	UUID_ONLY_PATTERN = `^` + UUID_PATTERN + `$`
)

var (
	uuid_regex      *regexp.Regexp = regexp.MustCompile(UUID_PATTERN)
	uuid_only_regex *regexp.Regexp = regexp.MustCompile(UUID_ONLY_PATTERN)
)

func ValidUUID(it interface{}) (ok bool, _ error) {
	var id string
	if id, ok = it.(string); !ok {
		return
	}

	ok = uuid_only_regex.MatchString(id)
	return
}

func ValidString(it interface{}) (ok bool, _ error) {
	_, ok = it.(string)
	return
}

func ValidNumber(it interface{}) (ok bool, _ error) {
	_, ok = it.(float64)
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

	_, ok = it.(float64)
	return
}

func OptionalBool(it interface{}) (ok bool, _ error) {
	if ok = it == nil; ok {
		return
	}

	_, ok = it.(bool)
	return
}
