package binding

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//EnableDecoderUseNumber用于在JSON
//解码器实例上调用UseNumber方法。 UseNumber导致解码器将一个数字作为数字（而不是float64）解组到//接口{}中。
var EnableDecodeUseNumber = false

var EnableDecoderDisallowUnknownFields = false

type jsonBinding struct{}

func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return fmt.Errorf("invalid request")
	}
	return decodeJSON(req.Body,obj)
}

func decodeJSON(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if EnableDecodeUseNumber {
		decoder.UseNumber() //将会把数字解码为interface
	}
	if EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}

	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
