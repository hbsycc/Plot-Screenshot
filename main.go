package main

import (
	"a.resources.cc/capture"
	"a.resources.cc/config"
	"a.resources.cc/demo"
	"log"
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
	demo.Hash()
	return
	_ = capture.Capture(&Files[0])
	return

	for _, file := range Files {
		if err := capture.Capture(&file); err != nil {
			log.Fatalln(err.Error())
		}
	}

	//err := database.InsertFiles()
	//if err != nil {
	//	return
	//}
}
