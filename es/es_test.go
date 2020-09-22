package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/olivere/elastic.v5"
	"reflect"
	"testing"
)

var (
	client *elastic.Client
	eIndex = "users"
	eType  = "user"
	err    error
)

type User struct {
	Name      string   `json:"name"`
	Age       int      `json:"age"`
	About     string   `json:"about"`
	Interests []string `json:"interests"`
}

func init() {
	client, err = elastic.NewClient(
		elastic.SetURL("http://10.10.114.123:9200/"),
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	info, code, err := client.Ping("http://10.10.114.123:9200/").Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := client.ElasticsearchVersion("http://10.10.114.123:9200/")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
}

//增加索引
func TestPostIndex(t *testing.T) {
	result, err := client.CreateIndex(eIndex).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

//查询索引
func TestGetIndex(t *testing.T) {
	indexs, err := client.IndexNames()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(indexs)
}

//增加记录
func TestPost(t *testing.T) {
	e1 := User{"Jane", 32, "I like to collect rock albums", []string{"music"}}
	put1, err := client.Index().Index(eIndex).Type(eType).Id("1").BodyJson(e1).Do(context.TODO())
	assert.Nil(t, err)
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)

	e2 := `{"name":"John","age":25,"about":"I love to go rock climbing","interests":["sports","music"]}`
	put2, err := client.Index().Index(eIndex).Type(eType).Id("2").BodyJson(e2).Do(context.TODO())
	assert.Nil(t, err)
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put2.Id, put2.Index, put2.Type)

	e3 := `{"name":"Douglas","age":35,"about":"I like to build cabinets","interests":["forestry"]}`
	put3, err := client.Index().Index(eIndex).Type(eType).Id("3").BodyJson(e3).Do(context.TODO())
	assert.Nil(t, err)
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put3.Id, put3.Index, put3.Type)
}

//自增id
func TestNewPost(t *testing.T) {
	e1 := User{"huye", 32, "I like to collect rock albums", []string{"music"}}
	put, err := client.Index().Index(eIndex).Type(eType).BodyJson(e1).Do(context.TODO())
	assert.Nil(t, err)
	fmt.Println(put.Id, put.Status)
}

//查询
func TestGet(t *testing.T) {
	//通过id查找
	get1, err := client.Get().Index(eIndex).Type(eType).Id("1").Do(context.TODO())
	if err != nil {
		panic(err)
	}
	if get1.Found {
		fmt.Println(string(*get1.Source))
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}
}

func TestSearch(t *testing.T) {
	var res *elastic.SearchResult

	//获取所有
	res, err = client.Search(eIndex).Type(eType).Do(context.TODO())
	assert.Nil(t, err)
	printUser(res)

	//字段相等
	q := elastic.NewQueryStringQuery("name:John")
	res, err = client.Search(eIndex).Type(eType).Query(q).Do(context.TODO())
	assert.Nil(t, err)
	printUser(res)

	if res.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d Employee \n", res.Hits.TotalHits)

		for _, hit := range res.Hits.Hits {

			var u User
			err := json.Unmarshal(*hit.Source, &u)
			assert.Nil(t, err)

			fmt.Printf("user name %s\n", u.Name)
		}
	} else {
		fmt.Println("not found")
	}

	//条件查询
	//年龄大于20
	boolQ := elastic.NewBoolQuery()
	boolQ.Must(elastic.NewMatchQuery("name", "John"))
	boolQ.Filter(elastic.NewRangeQuery("age").Gt(20))
	res, err = client.Search(eIndex).Type(eType).Query(boolQ).Do(context.TODO())
	assert.Nil(t, err)
	printUser(res)

	fmt.Println()
	//短语搜索
	matchphraseQuery := elastic.NewMatchPhraseQuery("about", "rock climbing")
	res, err = client.Search(eIndex).Type(eType).Query(matchphraseQuery).Do(context.TODO())
	assert.Nil(t, err)
	printUser(res)
}

func printUser(res *elastic.SearchResult) {
	var user User

	for _, item := range res.Each(reflect.TypeOf(user)) {
		t := item.(User)
		fmt.Printf("%#v\n", t)
	}
}

//删除文档
func TestDelete(t *testing.T) {
	result, err := client.Delete().Index(eIndex).Type(eType).Id("AXOZBSyuTi5QnmIB_yvp").Do(context.TODO())
	assert.Nil(t, err)
	fmt.Println(result.Id, result.Result)
}

//检测version
func TestPostVersion(t *testing.T) {
	str := "{\"title\":\"My first blog entry\",\"text\":\"Just trying this out...\"}"
	res, err := client.Index().Index(eIndex).Type(eType).BodyString(str).Do(context.TODO())
	assert.Nil(t, err)
	fmt.Println(res.Id)
	fmt.Println(res.Version)
}

func TestGetVersion(t *testing.T) {
	res, err := client.Get().Index(eIndex).Type(eType).Id("AXOZEZKjTi5QnmIB_yvt").Do(context.TODO())
	assert.Nil(t, err)

	fmt.Println(*res.Version)
	fmt.Println(res.Id)
	fmt.Println(string(*res.Source))
}

//指定version和id进行修改，第一次执行成功，version=2，第二次执行则会失败，因为version变了
func TestPutVersion(t *testing.T) {
	res, err := client.Update().Index(eIndex).Type(eType).Id("AXOZEZKjTi5QnmIB_yvt").Version(1).Doc(map[string]interface{}{"title": "hah"}).Do(context.TODO())
	assert.Nil(t, err)

	fmt.Println(res.Version)
}

//同时搜索多个文档
func TestGetMore(t *testing.T) {
	res, err := client.Get().Index(eIndex).Type(eType).Id("1").Id("2").Do(context.TODO())
	assert.Nil(t, err)

	fmt.Println(string(*res.Source))
}
