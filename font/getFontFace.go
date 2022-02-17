package font

import (
	"embed"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

//go:embed ttf
var fontDir embed.FS

// GetFontFace
// @Description: 读取、解析字体
// @param fontSize
// @return face
// @return err
func GetFontFace(fontSize float64) (face font.Face, err error) {
	readFile, err := fontDir.ReadFile("ttf/FZHei.TTF")
	if err != nil {
		return nil, err
	}

	f, err := truetype.Parse(readFile)
	if err != nil {
		return nil, err
	}

	face = truetype.NewFace(f, &truetype.Options{
		Size: fontSize,
	})

	return
}
