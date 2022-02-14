package main

import (
	"a.resources.cc/capture"
	"a.resources.cc/config"
)

func init() {
	if err := config.SetConfig(); err != nil {
		panic(err)
	}

	if err := MediaDirWalk(); err != nil {
		panic(err)
	}

	//fmt.Printf("文件夹总数：%v,文件总数：%v\n", len(Dirs), len(Files))
}

func main() {
	for _, file := range Files {
		_ = capture.Capture(&file)
	}

	//err := database.InsertFiles()
	//if err != nil {
	//	return
	//}
}
