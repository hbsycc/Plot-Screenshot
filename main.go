package main

import (
	"a.resources.cc/capture"
	"a.resources.cc/config"
	"a.resources.cc/dirFile"
	"log"
)

func init() {
	if err := config.SetConfig(); err != nil {
		panic(err)
	}
}

func main() {
	//c := make(chan model.File)

	if err := dirFile.MediaDirWalk(); err != nil {
		panic(err)
	}

	for _, file := range dirFile.Files {
		if err := capture.Capture(&file); err != nil {
			log.Fatalln(err.Error())
		}
	}

	//err := database.InsertFiles()
	//if err != nil {
	//	return
	//}
}
