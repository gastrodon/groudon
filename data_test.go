package groudon

import (
	"bytes"
	"encoding/json"
	"testing"
)

type Something struct {
	Name   string `json:"name"`
	Friend string `json:"friend"`
	Age    int    `json:"age"`
}

func a_string(it interface{}) (ok bool, _ error) {
	_, ok = it.(string)
	return
}

func a_number(it interface{}) (ok bool, _ error) {
	_, ok = it.(float64)
	return
}

func (some *Something) Validators() (data map[string]func(interface{}) (bool, error)) {
	data = map[string]func(interface{}) (bool, error){
		"name":   a_string,
		"friend": a_string,
		"age":    a_number,
	}

	return
}

func (some *Something) Defaults() (data map[string]interface{}) {
	data = map[string]interface{}{
		"friend": "jim",
	}

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
