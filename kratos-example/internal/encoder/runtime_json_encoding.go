// https://github.com/go-kratos/kratos/blob/main/encoding/json/json.go
package encoder

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-kratos/kratos/v2/encoding"
)

// Name is the name registered for the json codec.
const Name = "json-custom"

func init() {
	encoding.RegisterCodec(codec{})
}

// codec is a Codec implementation with json.
type codec struct{}

func (codec) Marshal(v interface{}) ([]byte, error) {

	fmt.Println("========================================================")
	fmt.Println("使用默认的 json 解码，不依赖 kratos 的 message type 判断")
	fmt.Println("========================================================")
	marshal, _ := json.Marshal(v)
	return marshal, nil
}

func (codec) Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	for rv := rv; rv.Kind() == reflect.Ptr; {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	return json.Unmarshal(data, v)
}

func (codec) Name() string {
	return Name
}
