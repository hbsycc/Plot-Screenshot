package model

type File struct {
	Name      string
	XxHash    string
	Ext       string
	Dir       string
	TempDir   string
	Path      string
	RePath    string
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
