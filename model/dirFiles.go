package model

type File struct {
	Name      string
	ReName    string
	Ext       string
	Dir       string
	Path      string
	RePath    string
	SHA1      string
	MediaInfo struct {
		DurationSeconds    int64
		DurationFormat     string
		Width              int
		Height             int
		DisplayAspectRatio string
		CodecName          string
		PixFmt             string
		RFrameRate         string
	}
}
