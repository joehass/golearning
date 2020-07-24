package gin

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

const defaultMultipartMemory = 32 << 20 // 32 MB

var (
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")
)

var defaultappEngine bool

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

//表示路由请求规范，包含了方法和路径
type RouteInfo struct {
	Method      string
	Path        string
	Handler     string
	HandlerFunc HandlerFunc
}

type RoutesInfo []RouteInfo

//容器实例，包含多路复用器，中间件和基础配置
type Engine struct {
	RouterGroup
	//重定向，如果当前路由没有匹配到任何已经存在的路由
	//例如，/foo/被请求，但是只有/foo存在，客户端会重定向到/foo，针对get请求会返回301,其他请求会返回307
	RedirectTrailingSlash bool

	//如果启用，路由将会尝试去修复当前的请求路径，如果没有找到任何配置的路由
	//首先，多余的路由元素，如 ../和//将会被移除
	//然后，会进行不区分大小写的查找，
	//如果能找到匹配的路由，会进行重定向，对于get请求返回301，其他请求会返回307
	//例如，/FOO和/..//FOO会被重定向到/foo
	RedirectFixPath bool

	//如果启用，路由器会检查当前路由是否允许其他方法，
	//如果请求没有找到匹配路由，会返回"Method Not Allowed"和405
	//如果没有其他方法被允许，这个请求会委托给 NotFound处理器处理
	HandleMethodNotAllowed bool
	ForwardedByClientIP    bool

	//如果启用，http head会强制加上'X-AppEngine...'，以便更好的和PaaS平台集成
	AppEngine bool

	//如果启用，将会使用url.RawPath查找参数
	UseRawPath bool

	//true：将不转义路径值，如果useRawPath是false，则UnescapePathValues为true
	//作为url.Path使用
	UnescapePathValues bool

	//赋予http请求的ParseMultipartForm调用maxMemory参数的值
	MaxMultipartMemory int64

	//去除多余/
	RemoveExtraSlash bool

	secureJSONPrefix string
	pool             sync.Pool
	trees            methodTrees
	maxParams        uint16
	noRoute          HandlersChain
	noMethod         HandlersChain
	allNoMethod      HandlersChain
	allNoRoute       HandlersChain
}

//创建一个新的空白的，没有中间件engine
func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		RedirectTrailingSlash:  true,
		RedirectFixPath:        false,
		HandleMethodNotAllowed: false,
		ForwardedByClientIP:    true,
		AppEngine:              defaultappEngine,
		UseRawPath:             false,
		RemoveExtraSlash:       false,
		UnescapePathValues:     true,
		MaxMultipartMemory:     defaultMultipartMemory,
		trees:                  make(methodTrees, 0, 0),
		secureJSONPrefix:       "while(1);",
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}

	return engine
}

func (engine *Engine) Routes() (routes RoutesInfo) {
	for _, tree := range engine.trees {
		routes = iterate("", tree.method, routes, tree.root)
	}
	return routes
}

func iterate(path, method string, routes RoutesInfo, root *node) RoutesInfo {
	path += root.path
	if len(root.handlers) > 0 {
		handlerFunc := root.handlers.Last()
		routes = append(routes, RouteInfo{
			Method:      method,
			Path:        path,
			Handler:     nameOfFunction(handlerFunc),
			HandlerFunc: handlerFunc,
		})
	}
	for _, child := range root.children {
		routes = iterate(path, method, routes, child)
	}
	return routes
}

func (engine *Engine) allocateContext() *Context {
	v := make(Params, 0, engine.maxParams)
	return &Context{engine: engine, params: &v}
}

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")
	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)

	if paramCount := countParams(path); paramCount > engine.maxParams {
		engine.maxParams = paramCount
	}
}

//重新输入已经被重写的上下文
//可以把c.Request.URL.PATH中的内容清空
func (engine *Engine) handleContext(c *Context) {
	oldIndexValue := c.index
	c.reset()
	engine.handleHTTPRequest(c)

	c.index = oldIndexValue
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	fmt.Println("handleHttpRequest")
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	unescape := false
	if engine.UseRawPath && len(c.Request.URL.RawPath) > 0 {
		rPath = c.Request.URL.RawPath
		unescape = engine.UnescapePathValues
	}

	//找到http method的根
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		value := root.getValue(rPath, c.params, unescape)
		if value.params != nil {

		}
		if value.handlers != nil {
			c.handlers = value.handlers
			c.fullPath = value.fullPath
			c.Next()
			c.writermem.WriteHeaderNow()
			return
		}
		//if httpMethod != "CONNECT" && rPath != "/"{
		//	if value.tsr && engine.RedirectTrailingSlash {
		//
		//	}
		//}
	}
	c.handlers = engine.allNoRoute
	serveError(c, http.StatusNotFound, default404Body)
}

var mimePlain = []string{MIMEPlain}

func serveError(c *Context, code int, defaultMessage []byte) {
	c.writermem.status = code
	c.Next()
	if c.writermem.Written() {
		return
	}
	if c.writermem.Status() == code {
		c.writermem.Header()["Content-Type"] = mimePlain
		_, err := c.Writer.Write(defaultMessage)
		if err != nil {
			log.Printf("cannot write message to writer during serve error: %v\n", err)
		}
		return
	}
	c.writermem.WriteHeaderNow()
}

func (engine *Engine) Run(addr ...string) (err error) {
	address := resolveAddress(addr)
	fmt.Printf("Listrning and serving HTTP on %s\n", address)
	err = http.ListenAndServe(address, engine)
	return
}

//实现http.Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.reset()
	engine.handleHTTPRequest(c)
	engine.pool.Put(c)
}

func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.noRoute = handlers
	engine.rebuild404Handlers()
}

func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
	engine.noMethod = handlers
	engine.rebuild405Handlers()
}

func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.Use(middleware...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
	return engine
}

func (engine *Engine) rebuild404Handlers() {
	engine.allNoRoute = engine.combineHandlers(engine.noRoute)
}

func (engine *Engine) rebuild405Handlers() {
	engine.allNoMethod = engine.combineHandlers(engine.noMethod)
}
