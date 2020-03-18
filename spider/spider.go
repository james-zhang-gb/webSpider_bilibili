package spider

import (
	//"fmt"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	URL "net/url"
	"strconv"
	//"time"
)

// GetList ... 获取下载列表
func GetList(url string) ([]ImgListInfo, error) {
	r, err := chromeGet(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	list := JSONDecode(b)
	ImgList := make([]ImgListInfo, len(list.Data.Items))
	fmt.Println("ImgList len =", len(ImgList))
	var ch chan int
	ch=make(chan int)
	for k, v := range list.Data.Items {
		go func(k int, s int) {
			ImgList[k].Tag = GetTag(s)
			ch <- 1
		}(k, v.Item.DocID)
		ImgList[k].Title = v.Item.Title
		//time.Sleep(time.Second)
		ImgList[k].DocID = v.Item.DocID
		ImgList[k].ImgSrc = make([]string, len(v.Item.Picture))
		for i, j := range v.Item.Picture {
			ImgList[k].ImgSrc[i] = j.ImgSrc
		}
	}
	fmt.Println("完成信息转移")
	for range list.Data.Items {
		<-ch
	}
	return ImgList, nil
}

// GetTag ... 获取图片标签
func GetTag(id int) string {
	origin := "https://api.vc.bilibili.com/link_draw/v1/doc/detail?doc_id="
	origin += strconv.Itoa(id)
	r, err := chromeGet(origin)
	if err != nil {
		fmt.Println("GetTag chromeGet err=", err)
		return "unknow"
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			fmt.Println("GetTag r.Body.Close() err=", err)
			return
		}
	}()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("GetTag ioutil.ReadAll err=", err)
		return "unknow"
	}
	var list tag
	if err := json.Unmarshal(b, &list); err != nil {
		fmt.Println("json 解码出错（getTag）", err)
		return "unknow"
	}
	if len(list.Datas.ItemStruct.Tags) == 0 {
		return "unknow"
	}
	fmt.Println(list.Datas.ItemStruct.Tags[0].TagText)
	return list.Datas.ItemStruct.Tags[0].TagText
}

// UseProxy ... 使用代理,留空则不使用代理
var UseProxy string = ""

//模仿浏览器访问链接
func chromeGet(url string) (r *http.Response, err error) {
	clt := http.DefaultClient
	// 添加代理地址
	if len(UseProxy) > 0 {
		proxy, _ := URL.Parse(UseProxy)
		tr := http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		clt.Transport = &tr
	}
	//clt.Timeout =90 * time.Second
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("User-Agent", " Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Origin", "https://h.bilibili.com")
	req.Header.Add("Referer", "https://h.bilibili.com/eden/draw_area")
	r, err = clt.Do(req)
	return
}

//JSONDecode ... 将json解析到结构体
func JSONDecode(body []byte) *URLList {
	var list URLList
	if err := json.Unmarshal(body, &list); err != nil {
		log.Fatal("json 解码出错", err)
	}
	return &list
}

// ImgListInfo ... 图片下载链接、标签、标题
type ImgListInfo struct {
	ImgSrc []string
	Title  string
	Tag    string
	DocID  int
}
