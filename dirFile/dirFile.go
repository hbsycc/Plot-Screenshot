package dirFile

import (
	"a.resources.cc/config"
	"a.resources.cc/lib"
	"a.resources.cc/model"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	err = filepath.Walk(dir, func(infoPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirName := info.Name()
			childDir := filepath.Join(infoPath, dirName)
			Dirs = append(Dirs, childDir)
		} else {
			fileName := info.Name()
			f := model.File{
				Name: fileName,
				Ext:  path.Ext(fileName),
				Dir:  filepath.Dir(infoPath),
				Path: infoPath,
			}

			isMedia := lib.StringsContains(config.GetConfig().Media.Ext, strings.ToLower(f.Ext))
			if isMedia {
				Files = append(Files, f)
			}
		}

		return err
	})

	return
}
