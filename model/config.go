package model

type Config struct {
	Debug bool `json:"debug"` // 是否输出debug信息
	Media struct {
		Dir []string `json:"dir"` // 媒体文件所在路径
		Ext []string `json:"ext"` // 匹配媒体文件后缀
	} `json:"media"`
	Capture struct {
		Dir     string `json:"dir"`    // 截图保存路径
		Count   int    `json:"count"`  // 生成截图数量
		Thread  int    `json:"thread"` // 生成截图时的线程数量
		Quality int    `json:"quality"`
		Grid    struct {
			Column    int `json:"column"`
			Row       int `json:"row"`
			ColumnGap int `json:"columnGap"`
			RowGap    int `json:"rowGap"`
		} `json:"grid"`
	} `json:"capture"`
}
