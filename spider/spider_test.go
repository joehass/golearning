package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
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

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(b)))
	if err != nil {
		log.Fatalln(err)
	}
	dom.Find("tbody").Each(func(i int, selection *goquery.Selection) {
		selection.Find("tr:first-child").Each(func(i int, selection *goquery.Selection) {
			selection.Find("td[class]").Each(func(i int, selection *goquery.Selection) {
				str := selection.Text()
				strArr := strings.Split(str, "\n")[1:6]
				fmt.Println(len(strArr))
				fmt.Println(strArr)
			})
		})
	})
}
