package lib

import "golang.org/x/text/encoding/simplifiedchinese"

const (
	UTF8    = "UTF-8"
	GB18030 = "GB18030"
)

func ConvertCodeToUTF8(sourceStr string, sourceCharsetCode string) string {
	var str string
	switch sourceCharsetCode {
	case GB18030:
		var decode, _ = simplifiedchinese.GB18030.NewDecoder().String(sourceStr)
		str = decode
	case UTF8:
		fallthrough
	default:
		str = sourceStr
	}
	return str
}

func ConvertCodeToGBK(sourceStr string, sourceCharsetCode string) string {
	var str string
	switch sourceCharsetCode {
	case UTF8:
		var decode, _ = simplifiedchinese.GB18030.NewEncoder().String(sourceStr)
		str = decode
	case GB18030:
		fallthrough
	default:
		str = sourceStr
	}
	return str
}
