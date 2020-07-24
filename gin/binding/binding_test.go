package binding

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type FooStruct struct {
	Foo string `msgpack:"foo" json:"foo" form:"foo" xml:"foo" binding:"required"`
}

func TestBindingJSONNilBody(t *testing.T) {
	var obj FooStruct
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	err := JSON.Bind(req,&obj)
	assert.Error(t,err)
}
