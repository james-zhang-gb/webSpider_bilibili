package spider

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
