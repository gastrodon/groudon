package groudon

import (
	"github.com/mitchellh/mapstructure"

	"encoding/json"
	"fmt"
	"io"
)

var (
	ErrInvalidTyping = fmt.Errorf("json body types do not match expected")
)

type Fillable interface {
	Validators() map[string]func(interface{}) (bool, error)
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

func validateValues(data map[string]interface{}, validators map[string]func(interface{}) (bool, error)) (valid bool, err error) {
	var key string
	var validator func(interface{}) (bool, error)
	for key, validator = range validators {
		if valid, err = validator(data[key]); !valid || err != nil {
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

	var valid bool
	if valid, internal = validateValues(data, target.Validators()); !valid || internal != nil {
		external = ErrInvalidTyping
		return
	}

	internal = encodeBody(data, target)
	return
}
