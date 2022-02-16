package model

type FFProbe struct {
	Streams []struct {
		Width              int    `json:"width"`
		Height             int    `json:"height"`
		PixFmt             string `json:"pix_fmt"`
		Duration           string `json:"duration"`
		BitRate            string `json:"bit_rate"`
		CodecName          string `json:"codec_name"`
		RFrameRate         string `json:"r_frame_rate"`
		AvgFrameRate       string `json:"avg_frame_rate"`
		CodecLongName      string `json:"codec_long_name"`
		DisplayAspectRatio string `json:"display_aspect_ratio"`
	} `json:"streams"`

	Format struct {
		Filename       string `json:"filename"`
		FormatName     string `json:"format_name"`
		FormatLongName string `json:"format_long_name"`
		Size           string `json:"size"`
	} `json:"format"`
}
