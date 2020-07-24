package gin

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

var _ context.Context = &Context{}

func TestContextMultipartForm(t *testing.T) {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	assert.NoError(t, mw.WriteField("foo", "bar"))
	w, err := mw.CreateFormFile("file", "test")
	if assert.NoError(t, err) {
		_, err = w.Write([]byte("test"))
		assert.NoError(t, err)
	}
	mw.Close()
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("POST", "/", buf)
	c.Request.Header.Set("Content-Type", mw.FormDataContentType())
	f, err := c.MultipartForm()
	if assert.NoError(t, err) {
		assert.NotNil(t, f)
	}
	assert.NoError(t, c.SaveUploadedFile(f.File["file"][0], "test"))

}

func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine) {
	r = New()
	c = r.allocateContext()
	c.reset()
	c.writermem.reset(w)
	return
}

func TestContextSetGet(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Set("foo", "bar")

	value, err := c.Get("foo")
	assert.Equal(t, "bar", value)
	assert.True(t, err)

	value, err = c.Get("foo2")
	assert.Nil(t, value)
	assert.False(t, err)
}

func TestContextQuery(t *testing.T) {
	c, _ := CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest("GET", "http://example.com/?foo=bar&page=10&id=", nil)

	value, ok := c.GetQuery("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", value)
	assert.Equal(t, "bar", c.Query("foo"))

}

func TestContetRenderJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)

	c.JSON(http.StatusCreated, H{"foo": "bar", "html": "<b>"})

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "{\"foo\":\"bar\",\"html\":\"\\u003cb\\u003e\"}", w.Body.String())
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestContextRenderAttachment(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := CreateTestContext(w)
	newFilename := "new_filename.go"

	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.FileAttachment("./gin.go", newFilename)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "func New() *Engine {")
	assert.Equal(t, fmt.Sprintf("attachment; filename=\"%s\"", newFilename), w.HeaderMap.Get("Content-Disposition"))
}
