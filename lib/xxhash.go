package lib

import (
	"github.com/cespare/xxhash/v2"
	"io"
	"os"
	"strconv"
)

// XxHash
// @Description: 计算文件 xxhash 值
// @param filePath
// @return sum
// @return err
func XxHash(filePath string) (sum string, err error) {
	file, _ := os.Open(filePath)
	h := xxhash.New()

	_, err = io.Copy(h, file)
	if err != nil {
		return
	}

	err = file.Close()
	if err != nil {
		return
	}

	sum = strconv.FormatUint(h.Sum64(), 10)
	return sum, err
}
