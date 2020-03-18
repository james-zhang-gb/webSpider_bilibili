package main

import (
	//"fmt"
	//"os"
	"fmt"
	"strconv"
	"webCrawler/spider"
)

func client(address string, ch chan<- int, num int) {
	l := spider.GetList(address)
	spider.SaveImgByTag(l)
	ch <- num
}
func main() {
	var ch chan int
	//并发爬取各页面
	for i := 1; i <= 20; i++ {
		go client("https://api.vc.bilibili.com/link_draw/v2/Doc/list?category=all&type=hot&page_num="+strconv.Itoa(i)+"page_size=20", ch, i)
	}
	for i := 1; i <= 20; i++ {
		fmt.Printf("第%d页爬取完成\n", <-ch)
	}

}
