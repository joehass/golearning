package render

import "net/http"

//继承json
type Render interface {
	//用自定义格式写数据
	Render(w http.ResponseWriter) error

	//写自定义的contentType
	WriteContentType(w http.ResponseWriter)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
