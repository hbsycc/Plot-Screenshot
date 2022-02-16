package main

import (
	"a.resources.cc/capture"
	"a.resources.cc/config"
	"a.resources.cc/dirFile"
	"a.resources.cc/lib"
	"a.resources.cc/model"
	"fmt"
	"log"
)

func init() {
	if err := config.SetConfig(); err != nil {
		panic(err)
	}

	if err := dirFile.MediaDirWalk(); err != nil {
		panic(err)
	}

	msg := fmt.Sprintf("文件夹总数：%v,文件总数：%v", len(dirFile.Dirs), len(dirFile.Files))
	lib.DebugLog(msg, "Init")
}

func main() {
	c := make(chan model.File, 10)

	go func() {
		for _, file := range dirFile.Files {
			err := dirFile.HashUnique(&file)
			if err != nil {
				log.Fatalln(err)
			}
			c <- file
		}
		close(c)
	}()

	for file := range c {
		if err := capture.Capture(&file); err != nil {
			log.Fatalln(err.Error())
		}
		lib.DebugLog("主线程处理完了", "main")
	}

	//for _, file := range dirFile.Files {
	//	if err := capture.Capture(&file); err != nil {
	//		log.Fatalln(err.Error())
	//	}
	//}

	//err := database.InsertFiles()
	//if err != nil {
	//	return
	//}
}
