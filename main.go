package main

import (
	"a.resources.cc/capture"
	"a.resources.cc/config"
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
	//if err := capture.Capture(&Files[0]); err != nil {
	//	fmt.Println(err)
	//}
	//return

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
