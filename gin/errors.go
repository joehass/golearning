package gin

import "reflect"

type ErrorType uint64

const (
	//表示一个私有的error
	ErrorTypePrivate ErrorType = 1 << 0
)

type Error struct {
	Err  error
	Type ErrorType
	Meta interface{}
}

func (msg Error) Error() string {
	return msg.Err.Error()
}

func (msg *Error) JSON() interface{} {
	jsonData := H{}
	if msg.Meta != nil {
		value := reflect.ValueOf(msg.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return msg.Meta
		case reflect.Map:
			for _, key := range value.MapKeys() {
				jsonData[key.String()] = value.MapIndex(key).Interface()
			}
		default:
			jsonData["meta"] = msg.Meta
		}
	}

	if _, ok := jsonData["error"]; !ok {
		jsonData["error"] = msg.Error()
	}

	return jsonData
}
