package model

type FFProbeStream struct {
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	PixFmt             string `json:"pix_fmt"`
	CodecName          string `json:"codec_name"`
	CodecLongName      string `json:"codec_long_name"`
	DisplayAspectRatio string `json:"display_aspect_ratio"`
}

type FFProbeFormat struct {
	Filename       string `json:"filename"`
	FormatName     string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	Duration       string `json:"duration"`
	Size           string `json:"size"`
	SizeFormat     string `json:"sizeFormat"`
}

type FFProbe struct {
	Streams []*FFProbeStream `json:"streams"`
	Format  *FFProbeFormat   `json:"format"`
}
