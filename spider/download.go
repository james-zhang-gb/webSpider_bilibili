package spider

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

//检测文件后缀名
func detectExtensionName(url string) string {
	l := len(url)
	for i := l - 1; i >= l-4; i-- {
		if url[i] == '.' {
			return url[i:]
		}
	}
	return ".png"
}

// DownloadImg 下载图片
func DownloadImg(fileName, url string) {
	fileName += detectExtensionName(url)
	r, err := chromeGet(url)
	
	if err != nil {
		fmt.Println("chromeGet err=",err)
		return
	}
	defer r.Body.Close()
	s, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err=",err)
		return
	}
	err = ioutil.WriteFile(fileName, s, 0666)
	if err != nil {
		fmt.Println("ioutil.WriteFile err=",err)
		fmt.Printf("保存图片[%s]失败", fileName)
		return
	}
	fmt.Printf("下载并保存图片[%s] 成功\n", fileName)
}

// SaveImgByTag ... 保存图片到对应标签文件夹
func SaveImgByTag(m []ImgListInfo){
	var ch chan int
	ch=make(chan int)
	chNum:=0
	for _,v :=range m{
		os.MkdirAll("image/recommed/"+v.Tag,os.ModePerm)
		num:=1
		for _,url:=range v.ImgSrc{
			fileName:="image/recommed/"+v.Tag+"/"+v.Title+"_"+strconv.Itoa(num)
			num++
			chNum++
			go func (fileName string,url string){
				DownloadImg(fileName,url)
				ch<-1
				}(fileName,url)
		}
		
	}
	for i:=0;i<chNum;i++{
		<-ch	
	}

}
