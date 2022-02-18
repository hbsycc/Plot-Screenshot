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

	var count uint64 = 1
	log.SetPrefix(fmt.Sprintf("[%v] ", count))
	for file := range c {
		log.Printf("处理：%v", file.Path)

		if err := capture.TempName(file); err != nil {
			log.Fatalln(err.Error())
		}
		if err := capture.Capture(&file); err != nil {
			capture.RecoverName(file)
			log.Fatalln(err.Error())
		}
		capture.RecoverName(file)
	}

	fmt.Printf("已完成！文件夹总数：%v,文件总数：%v", len(dirFile.Dirs), len(dirFile.Files))

	//err := database.InsertFiles()
	//if err != nil {
	//	return
	//}
}
