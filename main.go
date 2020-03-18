package main

import (
	//"fmt"
	//"os"
	"fmt"
	"strconv"
	urls "webCrawler/urlsData"
)


func client(address string,ch chan<- int,num int){
	l:=urls.GetList(address)
	urls.SaveImgByTag(l)
	ch<-num
}
func main() {
	var ch chan int
	for i:=1;i<=20;i++{
		go client("https://api.vc.bilibili.com/link_draw/v2/Doc/list?category=all&type=hot&page_num="+strconv.Itoa(i)+"page_size=20",ch,i)
	}
	for i:=1;i<=20;i++{
		fmt.Printf("第%d页爬取完成\n",<-ch)
	}
	
}

// func getFirstRepeatChar(strSrc string) byte {
// 	m := make(map[rune]bool)
// 	for i, v := range strSrc {
// 		if _, ok := m[v]; ok {

// 			return strSrc[i]
// 		}
// 		m[v] = true
// 	}
// 	return 0
// 	// write code here
// }
// func stringFilter( strSrc string ,  strPat string ) string {
// 	m:=make(map[rune]bool,len(strPat))
// 	for _,v:= range (strPat){
// 		if v==' '{
// 			continue
// 		}
// 		m[v]=true
// 	}
// 	s:=make([]rune,0)
// 	for _,v:=range(strSrc){
// 		if  _, ok := m[v]; !ok{
// 			s=append(s,v)
// 		}
// 	}
// 	return string(s)
//     // write code here
// }
