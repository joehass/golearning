package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

func Test_JQ(t *testing.T) {
	html := `<html>
            <body>
                <h1 id="title">春晓</h1>
                <p class="content1">
                春眠不觉晓，
                处处闻啼鸟。
                夜来风雨声，
                花落知多少。
                </p>
            </body>
            </html>
            `
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}

	dom.Find("p").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
}

func Test_JQ1(t *testing.T) {
	html := `<html>
            <body>
                <h1 id="title">春晓</h1>
                <p class="content1">
               	<span class="a">春眠不觉晓，</span>
                <span class="a">处处闻啼鸟。</span>
                <span class="a">夜来风雨声，</span>
                <span class="a">花落知多少。</span>
                </p>
            </body>
            </html>
            `
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}

	dom.Find("span[class=a]").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
}

func TestSpider(t *testing.T) {
	url := "https://www.caijinle.com/pailie5/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}
	fmt.Printf("%s", b)
}

func TestQUeryDom(t *testing.T) {
	dom, err := Grab("https://www.caijinle.com/pailie5/")
	if err != nil {
		log.Fatal(err)
		return
	}
	dom.Find("tbody tr:first-child td[class]").Each(func(i int, selection *goquery.Selection) {
		str := selection.Text()
		strArr := strings.Split(str, "\n")[1:6]
		fmt.Println(len(strArr))
		fmt.Println(strArr)

		node := selection.Nodes
		fmt.Println(node)
	})
}

func TestGetGrabText(t *testing.T) {
	str, err := GetGrabText("https://www.caijinle.com/pailie5/", "tbody tr:first-child td[class]")
	assert.Nil(t, err)
	//fmt.Println(str)

	str, err = GetGrabText("https://www.8200.cn/kjh/p5/", ".ballBox")
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)

	fmt.Println(str)
}

func GetGrabText(url, selection string) (string, error) {
	var str string
	dom, err := Grab(url)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	dom.Find(selection).Each(func(i int, selection *goquery.Selection) {
		str = selection.Text()
	})
	return str, nil
}

func Grab(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		return nil, err
	}

	return dom, nil
}
