package main

import (
	"a.resources.cc/config"
	"a.resources.cc/lib"
	"a.resources.cc/model"
	"fmt"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var Dirs []string

var Files []model.File

func MediaDirWalk() (err error) {
	for _, dir := range config.GetConfig().MediaDirs {
		err = walk(dir)
		if err != nil {
			return
		}
	}

	return
}

func walk(dir string) (err error) {

	err = filepath.Walk(dir, func(infoPath string, info os.FileInfo, err error) error {
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

			isMedia := lib.StringsContains(config.GetConfig().FilterExt, strings.ToLower(f.Ext))
			if isMedia {
				rand.Seed(time.Now().UnixNano())
				rename := fmt.Sprintf("%v%v", rand.Intn(9999999999999999), f.Ext)
				f.ReName = rename
				f.RePath = fmt.Sprintf("%v\\%v", f.Dir, rename)

				//if sha1, err := checksum.SHA1sum(f.Path); err != nil {
				//	return err
				//} else {
				//	f.SHA1 = sha1
				//}

				Files = append(Files, f)
			}
		}

		return err
	})

	return
}
