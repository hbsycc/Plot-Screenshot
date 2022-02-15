package main

import (
	"a.resources.cc/config"
	"a.resources.cc/lib"
	"a.resources.cc/model"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var Dirs []string

var Files []model.File

func MediaDirWalk() (err error) {
	for _, dir := range config.GetConfig().Media.Dir {
		if _, err := os.Stat(dir); err != nil {
			return err
		}

		err = walk(dir)
		if err != nil {
			return
		}
	}

	return
}

func walk(dir string) (err error) {
	fmt.Println(dir)
	err = filepath.Walk(dir, func(infoPath string, info os.FileInfo, err error) error {
		fmt.Println(info)
		if info.IsDir() {
			dirName := info.Name()
			childDir := infoPath + "\\" + dirName
			Dirs = append(Dirs, childDir)
		} else {
			name := info.Name()
			infoPaths := strings.Split(infoPath, "\\")
			f := model.File{
				Name: name,
				Ext:  path.Ext(name),
				Dir:  strings.Join(infoPaths[0:len(infoPaths)-1], "\\"),
				Path: infoPath,
			}

			isMedia := lib.StringsContains(config.GetConfig().Media.Ext, strings.ToLower(f.Ext))
			if isMedia {
				startTime := time.Now()
				if hash, err := lib.XxHash(f.Path); err != nil {
					return err
				} else {
					f.XxHash = hash
					f.TempDir = fmt.Sprintf("%v\\%v", config.GetConfig().Capture.Dir, hash)
					f.RePath = fmt.Sprintf("%v\\%v", f.Dir, hash+f.Ext)
				}
				lib.DebugLog(fmt.Sprintf("对文件xxhash,耗时：%v\n", time.Since(startTime)), "hash")

				Files = append(Files, f)
			}
		}

		return err
	})

	return
}
