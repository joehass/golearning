package gin

import (
	"crypto/tls"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func testRequest(t *testing.T, url string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	fmt.Println("start request")
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	assert.NoError(t, err)
	defer resp.Body.Close()
	body, ioerr := ioutil.ReadAll(resp.Body)
	assert.NoError(t, ioerr)
	assert.Equal(t, "it worked", string(body), "resp body should match")
	assert.Equal(t, "200 OK", resp.Status, "should get a 200")
}

func TestRunEmpty(t *testing.T) {
	router := New()

	go func() {
		router.GET("/example", func(c *Context) {
			c.String(http.StatusOK, "it worked")
		})
		router.Run(":8080")
	}()

	time.Sleep(5 * time.Millisecond)

	testRequest(t, "http://localhost:8080/example")
}

func TestConcurrentHandleContext(t *testing.T) {
	router := New()
	router.GET("/", func(c *Context) {
		c.Request.URL.Path = "/example"
		router.handleContext(c)
	})
	router.GET("/example", func(c *Context) { c.String(http.StatusOK, "it worked") })

	var wg sync.WaitGroup
	iterations := 200
	wg.Add(200)
	for i := 0; i < iterations; i++ {
		go func() {
			testGetRequestHandler(t, router, "/")
			wg.Done()
		}()
	}
	wg.Wait()
}

func testGetRequestHandler(t *testing.T, h http.Handler, url string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)
	time.Sleep(5 * time.Second)
	w := httptest.NewRecorder()
	fmt.Println("startRequest")
	h.ServeHTTP(w, req)

	assert.Equal(t, "it worked", w.Body.String(), "resp body should match")
	assert.Equal(t, 200, w.Code, "should get a 200")
}
