package model

type Config struct {
	Debug     bool     `json:"debug"`     // 是否输出debug信息
	MediaDirs []string `json:"mediaDirs"` // 媒体文件所在路径
	FilterExt []string `json:"filterExt"` // 匹配媒体文件后缀
	Capture   struct {
		Count  int    `json:"count"`  // 生成截图数量
		Thread int    `json:"thread"` // 生成截图时的线程数量
		Dir    string `json:"dir"`    // 截图保存路径
	} `json:"capture"`
}
