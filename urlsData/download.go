package urls

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
		fmt.Println(err)
		return
	}
	s, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile(fileName, s, 0666)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("保存图片[%s]失败", fileName)
		return
	}
	fmt.Printf("下载并保存图片[%s] 成功\n", fileName)
}

// SaveImgByTag ... 保存图片到对应标签文件夹
func SaveImgByTag(m []ImgListInfo){
	for _,v :=range m{
		os.MkdirAll("image/"+v.Tag,os.ModePerm)
		num:=1
		for _,url:=range v.ImgSrc{
			fileName:="image/"+v.Tag+"/"+v.Title+"_"+strconv.Itoa(num)
			num++
			DownloadImg(fileName,url)
		}
	}

}
