package groudon

import (
	"github.com/mitchellh/mapstructure"

	"encoding/json"
	"fmt"
	"io"
	"strings"
)

var (
	ErrInvalidTyping = fmt.Errorf("json body types do not match expected")
)

type Fillable interface {
	Types() map[string]string
	Defaults() map[string]interface{}
}

func withDefaults(data, defaults map[string]interface{}) (populated map[string]interface{}) {
	populated = data

	var key string
	var it interface{}
	for key, it = range defaults {
		if populated[key] == nil {
			populated[key] = it
		}
	}

	return
}

func validateTypes(data map[string]interface{}, types map[string]string) (valid bool) {
	var key, desired string
	for key, desired = range types {
		if strings.HasPrefix(desired, "[]") {
			continue
		}

		switch desired {
		case "float64", "number":
			_, valid = data[key].(float64)
		case "float", "float32":
			_, valid = data[key].(float32)
		case "int8", "byte":
			_, valid = data[key].(int8)
		case "int16":
			_, valid = data[key].(int16)
		case "int", "int32", "rune":
			_, valid = data[key].(int32)
		case "int64":
			_, valid = data[key].(int64)
		case "bool":
			_, valid = data[key].(bool)
		case "string":
			_, valid = data[key].(string)
		case "complex64":
			_, valid = data[key].(complex64)
		case "complex128":
			_, valid = data[key].(complex128)
		}

		if !valid {
			break
		}

	}

	return
}

func encodeBody(data map[string]interface{}, target Fillable) (err error) {
	var config mapstructure.DecoderConfig = mapstructure.DecoderConfig{
		Metadata: nil,
		TagName:  "json",
		Result:   target,
	}

	var decoder *mapstructure.Decoder
	if decoder, err = mapstructure.NewDecoder(&config); err == nil {
		err = decoder.Decode(data)
	}

	return
}

func SerializeBody(reader io.Reader, target Fillable) (internal, external error) {
	var closable bool
	if _, closable = reader.(io.Closer); closable {
		defer reader.(io.Closer).Close()
	}

	var data map[string]interface{}
	if external = json.NewDecoder(reader).Decode(&data); external != nil {
		return
	}

	data = withDefaults(data, target.Defaults())

	if !validateTypes(data, target.Types()) {
		external = ErrInvalidTyping
		return
	}

	internal = encodeBody(data, target)
	return
}
