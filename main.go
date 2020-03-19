package main

import (
	"strconv"
	"webCrawler/spider"
)

//var url1 string = "https://api.vc.bilibili.com/link_draw/v2/Doc/index?type=recommend&page_num=0&page_size=45"
//var dosDetail string = "https://api.vc.bilibili.com/link_draw/v1/doc/detail?doc_id=1093475"
//var rank string = "https://api.vc.bilibili.com/link_draw/v2/Doc/ranklist?biz=1&category=&rank_type=month&date=2020-03&page_num=0&page_size=50"
//var photoRecommed string = "https://api.vc.bilibili.com/link_draw/v2/Photo/index?type=recommend&page_num=0&page_size=45"

func main() {
	robot := spider.NewRobot()
	robot.JSONFileName = "recommend_2.json"
	robot.SaveFilePath = "image/recommend/"
	robot.Proxy = "http://127.0.0.1:1080"
	robot.ImgDetailURL = "https://api.vc.bilibili.com/link_draw/v1/doc/detail?doc_id="
	robot.RequestHeader = map[string]string{
		"User-Agent": " Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
		"accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Origin":     "https://h.bilibili.com",
		"Referer":    "https://h.bilibili.com/eden/draw_area",
	}
	for i := 0; i < 5; i++ {
		robot.URLList = append(robot.URLList, "https://api.vc.bilibili.com/link_draw/v2/Doc/index?type=recommend&page_num="+strconv.Itoa(i)+"&page_size=45")
	}
	robot.Run()
}

//func client(address string, ch chan<- int, num int) {
//	l, err := spider.GetList(address)
//	if err != nil {
//		fmt.Println("spider.GetList err=", err)
//		ch <- num
//		return
//	}
//	spider.SaveImgByTag(l, "image/")
//	ch <- num
//}
//func downloadImageFromJSON(jsonName string, path string) {
//	list := make([]spider.ImgListInfo, 0)
//	f, _ := ioutil.ReadFile(jsonName)
//	json.Unmarshal(f, &list)
//	spider.SaveImgByTag(list, path)
//}
//func saveImgSrcJSON(pagea, pageb int, jsonName string) {
//	list := make([]spider.ImgListInfo, 0)
//	var ch chan []spider.ImgListInfo
//	ch = make(chan []spider.ImgListInfo)
//	var errCh chan int
//	errCh = make(chan int)
//	for i := pagea; i <= pageb; i++ {
//		go func(i int) {
//			l, err := spider.GetList("https://api.vc.bilibili.com/link_draw/v2/Doc/index?type=recommend&page_num=" + strconv.Itoa(i) + "&page_size=45")
//			if err != nil {
//				fmt.Println("spider.GetList err=", err)
//				errCh <- 1
//			}
//			ch <- l
//		}(i)
//	}
//	for i := pagea; i <= pageb; i++ {
//		select {
//		case l := <-ch:
//			list = append(list, l...)
//			outputData, err := json.Marshal(list)
//			if err != nil {
//				fmt.Println("json.Marshal err=", err)
//				return
//			}
//			if err := ioutil.WriteFile(jsonName, outputData, 0666); err != nil {
//				fmt.Println("ioutil.WriteFile err=", err)
//			}
//			fmt.Println("已完成", i+1-pagea)
//		case _ = <-errCh:
//			fmt.Println("一个表出错")
//		}
//	}
//
//}
//
//func saveImgSrcToJSONSlow(pagea, pageb int) {
//	list := make([]spider.ImgListInfo, 0)
//	for i := pagea; i <= pageb; i++ {
//		l, err := spider.GetList(rank)
//		if err != nil {
//			fmt.Println("spider.GetList err=", err)
//		}
//		list = append(list, l...)
//		outputData, err := json.Marshal(list)
//		if err != nil {
//			fmt.Println("json.Marshal err=", err)
//			return
//		}
//		if err := ioutil.WriteFile("rank_month.json", outputData, 0666); err != nil {
//			fmt.Println("ioutil.WriteFile err=", err)
//		}
//	}
//}
