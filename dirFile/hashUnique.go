package dirFile

import (
	"a.resources.cc/config"
	"a.resources.cc/lib"
	"a.resources.cc/model"
	"fmt"
	"time"
)

// HashUnique
// @Description: 对文件进行 XxHash 作为唯一标识
// @param file
// @return err
func HashUnique(file *model.File) (err error) {
	var hash string
	startTime := time.Now()
	if hash, err = lib.XxHash(file.Path); err != nil {
		return err
	} else {
		file.XxHash = hash
		file.RePath = fmt.Sprintf("%v\\%v", file.Dir, hash+file.Ext)
		if len(config.GetConfig().Capture.Dir) == 0 {
			file.TempDir = fmt.Sprintf("%v\\screenshot\\%v", file.Dir, hash)
		} else {
			file.TempDir = fmt.Sprintf("%v\\%v", config.GetConfig().Capture.Dir, hash)
		}
	}
	lib.DebugLog(fmt.Sprintf("对文件xxhash,耗时：%v", time.Since(startTime)), "hash")

	return
}
