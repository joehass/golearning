package gin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddRoute(t *testing.T) {
	router := New()
	router.addRoute("GET", "/", HandlersChain{func(context *Context) {}})

	assert.Len(t, router.trees, 1)
	assert.NotNil(t, router.trees.get("GET"))
	assert.Nil(t, router.trees.get("POST"))

	router.addRoute("POST", "/", HandlersChain{func(_ *Context) {}})

	assert.Len(t, router.trees, 2)

	router.addRoute("POST", "/post", HandlersChain{func(_ *Context) {}})

	router.addRoute("POST", "/post2", HandlersChain{func(_ *Context) {}})

	fmt.Println(router.trees)
}

func TestListOfRoutes(t *testing.T) {
	router := New()

	router.GET("/favicon.ico", handlerTest1)
	router.GET("/", handlerTest1)
	group := router.Group("/users")
	{
		group.GET("/", handlerTest2)
		group.GET("/:id", handlerTest1)
	}
	router.Static("/static", ".")

	list := router.Routes()
	assert.Len(t, list, 6)

}

func TestEngineHandleContext(t *testing.T) {
	r := New()
	r.GET("/", func(c *Context) {
		c.Request.URL.Path = "/v2"

	})
}

func handlerTest1(c *Context) {}
func handlerTest2(c *Context) {}
