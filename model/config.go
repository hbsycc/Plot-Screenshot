package model

type Config struct {
	Debug      bool   `json:"debug"`      // 是否输出debug信息
	OutDirName string `json:"outDirName"` // 输出文件夹名称，为空则输出在媒体文件相同目录
	Media      struct {
		Dir []string `json:"dir"` // 媒体文件所在路径
		Ext []string `json:"ext"` // 匹配媒体文件后缀
	} `json:"media"`
	Capture struct {
		Thread      int `json:"thread"`  // 生成截图时的线程数量
		Quality     int `json:"quality"` // 截图质量
		ResizeWidth int `json:"width"`   // 截图缩放最大宽度,默认0不缩放
		Grid        struct {
			Column      int     `json:"column"`
			Row         int     `json:"row"`
			ColumnGap   int     `json:"columnGap"`
			RowGap      int     `json:"rowGap"`
			BorderWidth float64 `json:"borderWidth"` // 截图边框宽度，默认0不显示
		} `json:"grid"`
	} `json:"capture"`
}
