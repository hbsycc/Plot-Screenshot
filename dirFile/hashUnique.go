package dirFile

import (
	"a.resources.cc/config"
	"a.resources.cc/lib"
	"a.resources.cc/model"
	"fmt"
	"path/filepath"
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
		file.RePath = filepath.Join(file.Dir, hash+file.Ext)
		file.TempDir = filepath.Join(file.Dir, config.GetConfig().OutDirName, hash)
	}
	lib.DebugLog(fmt.Sprintf("对文件xxhash,耗时：%v", time.Since(startTime)), "hash")

	return
}
