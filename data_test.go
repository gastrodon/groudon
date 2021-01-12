package groudon

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

type Something struct {
	Name   string `json:"name"`
	Friend string `json:"friend"`
	Age    int    `json:"age"`
}

func longString(it interface{}) (ok bool, err error) {
	if ok, err = ValidString(it); !ok || err != nil {
		return
	}

	ok = len(it.(string)) >= 3
	return
}

func (some *Something) Validators() (data map[string]func(interface{}) (bool, error)) {
	data = map[string]func(interface{}) (bool, error){
		"name":   ValidString,
		"friend": longString,
		"age":    ValidNumber,
	}

	return
}

func (some *Something) Defaults() (data map[string]interface{}) {
	data = map[string]interface{}{
		"friend": "jim",
	}

	return
}

type verifyCloser struct {
	io.Reader
	Closed bool
}

func (it *verifyCloser) Close() (err error) {
	it.Closed = true
	return
}

func Test_something(test *testing.T) {
	var data map[string]interface{} = map[string]interface{}{
		"name": "foo",
		"age":  20,
	}

	var err error
	var data_bytes []byte
	if data_bytes, err = json.Marshal(data); err != nil {
		test.Fatal(err)
	}

	var target Something
	var external error
	if err, external = SerializeBody(bytes.NewReader(data_bytes), &target); err != nil {
		test.Fatal(err)
	}

	if external != nil {
		test.Fatal(external)
	}

	if target.Name != data["name"].(string) {
		test.Errorf("data name mismatch! have: %s, want: %s", target.Name, data["name"])
	}

	var defaults map[string]interface{} = target.Defaults()
	if target.Friend != defaults["friend"].(string) {
		test.Errorf("default friend not taken! have: %s, want: %s", target.Friend, defaults["friend"])
	}

}

func Test_something_badRequest(test *testing.T) {
	var target Something
	var err, external error
	if err, external = SerializeBody(bytes.NewReader([]byte("?")), &target); err != nil {
		test.Fatal(err)
	}

	if external == nil {
		test.Errorf("invalid body was parsed %#v", target)
	}
}

func Test_something_nil(test *testing.T) {
	var target Something
	var err, external error
	if err, external = SerializeBody(nil, &target); err != nil {
		test.Fatal(err)
	}

	if external != ErrNilBody {
		test.Errorf("nil body was parsed %#v", target)
	}
}

func Test_something_badTypes(test *testing.T) {
	var data map[string]interface{} = map[string]interface{}{
		"name":   "zero",
		"age":    0,
		"friend": ":(",
	}

	var dataBytes []byte
	dataBytes, _ = json.Marshal(data)

	var target Something
	var err, external error
	if err, external = SerializeBody(bytes.NewReader(dataBytes), &target); err != nil {
		test.Fatal(err)
	}

	if external != ErrInvalidTyping {
		test.Fatalf("incorrect err, %v != %v", external, ErrInvalidTyping)
	}
}

func Test_something_closed(test *testing.T) {
	var data map[string]interface{} = map[string]interface{}{
		"name": "zero",
		"age":  0,
	}

	var dataBytes []byte
	dataBytes, _ = json.Marshal(data)

	var reader *verifyCloser = &verifyCloser{bytes.NewReader(dataBytes), false}

	var target Something
	var err, external error
	if err, external = SerializeBody(reader, &target); err != nil {
		test.Fatal(err)
	}

	if external != nil {
		test.Fatal(external)
	}

	if !reader.Closed {
		test.Fatal("reader.Close was never called")
	}
}
