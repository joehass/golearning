package main

import "fmt"

//配合文章学习：https://colobu.com/2019/11/18/how-is-the-bigcache-is-fast/
//https://pengrl.com/p/35302/
//bigCache主要思想

//缓存分片和避免gc
//缓存分片：提高缓存并发
func main() {
	cache := newCache()
	cache.set("key", []byte("the value"))

	value, err := cache.get("key")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(value))
}
