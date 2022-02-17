package model

type Config struct {
	Debug bool `json:"debug"` // 是否输出debug信息
	Media struct {
		Dir []string `json:"dir"` // 媒体文件所在路径
		Ext []string `json:"ext"` // 匹配媒体文件后缀
	} `json:"media"`
	Capture struct {
		Dir         string `json:"dir"`     // 截图保存路径，空值则自动在媒体文件下创建 screenshot 目录
		Thread      int    `json:"thread"`  // 生成截图时的线程数量
		Quality     int    `json:"quality"` // 截图质量
		ResizeWidth int    `json:"width"`   // 截图缩放最大宽度,默认0不缩放
		Grid        struct {
			Column    int `json:"column"`
			Row       int `json:"row"`
			ColumnGap int `json:"columnGap"`
			RowGap    int `json:"rowGap"`
		} `json:"grid"`
	} `json:"capture"`
}
