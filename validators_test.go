package groudon

import (
	"github.com/google/uuid"

	"testing"
)

type validatorCase struct {
	Validator func(interface{}) (bool, error)
	Argument  interface{}
	Ok        bool
	Error     error
}

func Test_validators(test *testing.T) {
	var cases []validatorCase = []validatorCase{
		validatorCase{ValidUUID, uuid.New().String(), true, nil},
		validatorCase{ValidUUID, "nothing good", false, nil},
		validatorCase{ValidUUID, 42069, false, nil},
		validatorCase{ValidEmail, "foo@bar.io", true, nil},
		validatorCase{ValidEmail, "foo@bar", true, nil},
		validatorCase{ValidEmail, "foobaz", false, nil},
		validatorCase{ValidEmail, "foo@baz@bar", false, nil},
		validatorCase{ValidEmail, 0, false, nil},
		validatorCase{ValidString, "", true, nil},
		validatorCase{ValidString, 0, false, nil},
		validatorCase{ValidStringSlice, []string{"foo", "bar"}, true, nil},
		validatorCase{ValidStringSlice, []interface{}{"foo", "bar"}, true, nil},
		validatorCase{ValidStringSlice, nil, true, nil},
		validatorCase{ValidStringSlice, []byte{1, 2, 3}, false, nil},
		validatorCase{ValidStringSlice, []interface{}{"foo", "bar", 42}, false, nil},
		validatorCase{ValidStringSlice, "", false, nil},
		validatorCase{ValidNumber, 400, true, nil},
		validatorCase{ValidNumber, 400.69, true, nil},
		validatorCase{ValidNumber, -400, true, nil},
		validatorCase{ValidNumber, -400.69, true, nil},
		validatorCase{ValidNumber, nil, false, nil},
		validatorCase{ValidNumber, "nil", false, nil},
		validatorCase{ValidBool, true, true, nil},
		validatorCase{ValidBool, false, true, nil},
		validatorCase{ValidBool, nil, false, nil},
		validatorCase{ValidBool, 0, false, nil},
		validatorCase{ValidBool, 1, false, nil},
		validatorCase{OptionalString, "", true, nil},
		validatorCase{OptionalString, nil, true, nil},
		validatorCase{OptionalString, 0, false, nil},
		validatorCase{OptionalNumber, 0, true, nil},
		validatorCase{OptionalNumber, 0.123, true, nil},
		validatorCase{OptionalNumber, -10, true, nil},
		validatorCase{OptionalNumber, -0.123, true, nil},
		validatorCase{OptionalNumber, nil, true, nil},
		validatorCase{OptionalNumber, "0", false, nil},
		validatorCase{OptionalBool, false, true, nil},
		validatorCase{OptionalBool, true, true, nil},
		validatorCase{OptionalBool, nil, true, nil},
		validatorCase{OptionalBool, "nil", false, nil},
	}

	var index int
	var single validatorCase
	for index, single = range cases {
		var ok bool
		var err error
		if ok, err = single.Validator(single.Argument); single.Error != err {
			test.Fatalf("error incorrect at %d, %v != %v", index, single.Error, err)
		}

		if single.Ok != ok {
			test.Fatalf("ok incorrect at %d, %t != %t", index, single.Ok, ok)
		}
	}
}
