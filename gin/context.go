package gin

import (
	"fmt"
	"golearning/gin/binding"
	"golearning/gin/render"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

const (
	MIMEPlain = binding.MIMEPlain
)

const abortIndex int8 = math.MaxInt8 / 2

//gin中最重要的部分，他允许我们在中间件中传递变量，组织流，校验请求json和返回响应
type Context struct {
	engine   *Engine
	params   *Params
	handlers HandlersChain
	index    int8
	fullPath string

	Request   *http.Request
	Writer    ResponseWriter
	writermem responseWriter
	Params    Params

	mu sync.RWMutex

	//使用url.ParseQuery 来获取url中的value
	queryCache url.Values

	//针对每个请求的key/value键值对
	Keys map[string]interface{}
}

//返回upl中的参数
//这是c.Params.ByName(key)的快捷方式
func (c *Context) Param(key string) string {
	return c.params.ByName(key)
}

func (c *Context) reset() {
	c.Writer = &c.writermem
	c.Params = c.Params[0:0]
	c.handlers = nil
	c.index = -1

	c.fullPath = ""
	c.Keys = nil
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.Render(code, render.String{Format: format, Data: values})
}

//将数据写入请求头，并回调相应的数据格式来渲染数据
func (c *Context) Render(code int, r render.Render) {
	c.Status(code)

	if !bodyAllowedForStatus(code) {

		return
	}

	if err := r.Render(c.Writer); err != nil {
		panic(err)
	}

}

//设置http响应code
func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

//在内部中间件使用
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}

//获取url中的value
func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) QueryArray(key string) []string {
	values, _ := c.GetQueryArray(key)
	return values
}

func (c *Context) GetQueryArray(key string) ([]string, bool) {
	c.initQueryCache()
	if values, ok := c.queryCache[key]; ok && len(values) > 0 {
		return values, true
	}

	return []string{}, false
}

func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		c.queryCache = c.Request.URL.Query()
	}
}

//解析表单上传
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Request.ParseMultipartForm(c.engine.MaxMultipartMemory)
	return c.Request.MultipartForm, err
}

func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err

}

func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
	c.mu.Unlock()
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done always returns nil (chan which will wait forever),
// if you want to abort your work when the connection was closed
// you should use Request.Context().Done() instead.
func (c *Context) Done() <-chan struct{} {
	return nil
}

// Err always returns nil, maybe you want to use Request.Context().Err() instead.
func (c *Context) Err() error {
	return nil
}

func (c *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return c.Request
	}
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}

func (c *Context) JSON(code int, obj interface{}) {
	c.Render(code, render.JSON{Data: obj})
}

//FileAttachment以一种有效的方式将指定的文件写入主体流中
//在客户端，通常将使用给定的文件名下载文件
func (c *Context) FileAttachment(filepath, filename string) {
	c.Writer.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	http.ServeFile(c.Writer, c.Request, filepath)
}

func (c *Context) requestHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Context) Header(key, value string) {
	if value == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Set(key, value)
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.Writer.WriteHeaderNow()
	c.Abort()
}
