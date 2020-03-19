package spider

import (
	json "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)
import URL "net/url"

type Robot struct {
	RequestHeader map[string]string
	JSONFileName  string
	SaveFilePath  string
	URLList       []string
	Proxy         string
	ImgDetailURL  string
}

func NewRobot() *Robot {
	return &Robot{
		RequestHeader: make(map[string]string),
		JSONFileName:  "default.json",
		SaveFilePath:  "",
		URLList:       make([]string, 0),
		Proxy:         "",
		ImgDetailURL:  "",
	}
}

func (r *Robot) httpGet(url string) (res *http.Response, err error) {
	clt := http.DefaultClient
	// 添加代理地址
	if len(r.Proxy) > 0 {
		proxy, _ := URL.Parse(r.Proxy)
		tr := http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		clt.Transport = &tr
	}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	for k, v := range r.RequestHeader {
		req.Header.Add(k, v)
	}
	res, err = clt.Do(req)
	return
}
func (s *Robot) GetList(url string) ([]ImgListInfo, error) {
	r, err := s.httpGet(url)
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
	ch = make(chan int)
	for k, v := range list.Data.Items {
		go func(k int, id int) {
			ImgList[k].Tag = s.GetTag(id)
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
	for range list.Data.Items {
		<-ch
	}
	return ImgList, nil
}
func (s *Robot) GetTag(id int) string {
	origin := s.ImgDetailURL
	origin += strconv.Itoa(id)
	r, err := s.httpGet(origin)
	if err != nil {
		fmt.Println("GetTag httpGet err=", err)
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

func (s *Robot) SaveJSON() {
	list := make([]ImgListInfo, 0)
	var ch chan []ImgListInfo
	ch = make(chan []ImgListInfo)
	var errCh chan int
	errCh = make(chan int)
	for _, v := range s.URLList {
		go func(i string) {
			l, err := s.GetList(v)
			if err != nil {
				fmt.Println("func (s *Robot)SaveJSON() GetList err=", err)
				errCh <- 1
			}
			ch <- l
		}(v)
	}
	errNum := 0
	for range s.URLList {
		select {
		case l := <-ch:
			list = append(list, l...)
			outputData, err := json.Marshal(list)
			if err != nil {
				fmt.Println("json.Marshal err=", err)
				return
			}
			if err := ioutil.WriteFile(s.JSONFileName, outputData, 0666); err != nil {
				fmt.Println("ioutil.WriteFile err=", err)
			}
		case _ = <-errCh:
			errNum++
			fmt.Println(errNum, "个表出错")
		}
	}
	fmt.Println("已保存 ", s.JSONFileName)
}
func (s *Robot) Download() {
	list := make([]ImgListInfo, 0)
	f, _ := ioutil.ReadFile(s.JSONFileName)
	err := json.Unmarshal(f, &list)
	if err != nil {
		fmt.Println("func (s *Robot)Download() json.Unmarshal err=", err)
	}
	SaveImgByTag(list, s.SaveFilePath)
}
func (s *Robot) Run() {
	s.SaveJSON()
	s.Download()
}
