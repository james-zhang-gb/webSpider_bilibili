package spider

import (
	//"fmt"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	URL "net/url"
)

// GetList ... 获取下载列表
func GetList(url string) ([]ImgListInfo) {
	r, err := chromeGet(url)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
	}
	b,err:=ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	list:=JsonDecode(b)
	ImgList:=make([]ImgListInfo,len(list.Data.Items))
	fmt.Println("ImgList len =",len(ImgList))
	for k,v :=range list.Data.Items{
		ImgList[k].Title=v.Item.Title
		//time.Sleep(time.Second)
		ImgList[k].Tag=GetTag(v.Item.DocID)
		ImgList[k].ImgSrc=make([]string,len(v.Item.Picture) )
		for i,j:=range v.Item.Picture{
			ImgList[k].ImgSrc[i]=j.ImgSrc
		}
	}
	return ImgList
}
type tagList struct{
	TagText string `json:"tag"`
}
type tag struct{
	Datas struct{
		ItemStruct struct{
			Tags []tagList `json:"tags"`
		} `json:"item"`
	} `json:"data"`
}

// GetTag ... 获取图片标签
func GetTag(id int)string{
	origin:="https://api.vc.bilibili.com/link_draw/v1/doc/detail?doc_id="
	origin+=strconv.Itoa(id)
	r,err:=chromeGet(origin)
	defer func() {
		err:=r.Body.Close()
		if err!=nil{
			fmt.Println(err)
			return
		}
	}()
	b,err:=ioutil.ReadAll(r.Body)
	if err!=nil{
		fmt.Println(err)
		return "unknow"
	}
	var list tag
	if err := json.Unmarshal(b, &list); err != nil {
		log.Fatal(err)
	}
	if len(list.Datas.ItemStruct.Tags)==0{
		return "unknow"
	}
	fmt.Println(list.Datas.ItemStruct.Tags[0].TagText)
	return list.Datas.ItemStruct.Tags[0].TagText
}
//模仿浏览器访问链接
func chromeGet(url string) (r *http.Response, err error) {
	clt := http.DefaultClient
	// 添加代理地址
	proxy, _ := URL.Parse("http://127.0.0.1:1080")
	tr := http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	clt.Transport=&tr
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("User-Agent", " Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	r, err = clt.Do(req)
	return
}

// URLList 图片集列表及标签
type URLList struct {
	Data data `json:"data"`
}
type data struct {
	Items []items `json:"items`
}
type items struct {
	Item item `json:"item"`
}

// Item ... 图片列表
type item struct {
	Picture []picture `json:"pictures"`
	Title   string    `json:"title"`
	DocID int `json:"doc_id"`
}

//picture ... picture 层结构
type picture struct {
	ImgSrc string `json:"img_src"`
}

//JsonDecode ... 将json解析到结构体
func JsonDecode(body []byte) *URLList {
	var list URLList
	if err := json.Unmarshal(body, &list); err != nil {
		log.Fatal(err)
	}
	return &list
}

type ImgListInfo struct{
	ImgSrc []string
	Title string
	Tag string
}
