package encoder

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// DefaultResponseEncoder 默认实现
func DefaultResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// 通过Request Header的Accept中提取出对应的编码器
	// 如果找不到则忽略报错，并使用默认json编码器
	codec, _ := http.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}
	// 在Response Header中写入编码器的scheme
	w.Header().Set("Content-Type", ContentType(codec.Name()))
	w.Write(data)
	return nil
}

func CustomResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// 通过Request Header的Accept中提取出对应的编码器
	// 如果找不到则忽略报错，并使用默认json编码器
	//codec, _ := http.CodecForRequest(r, "Accept")

	// 使用自定义的 json 编码器
	customJsonCodec := encoding.GetCodec("json-custom")
	//customJsonCodecPlus := encoding.GetCodec("json-custom-plus")
	fmt.Printf("%+v \n", customJsonCodec)
	data, err := customJsonCodec.Marshal(v)
	if err != nil {
		return err
	}

	// 在Response Header中写入编码器的scheme
	w.Header().Set("Content-Type", ContentType(customJsonCodec.Name()))
	w.Write(data)
	return nil
}
