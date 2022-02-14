package capture

import (
	"a.resources.cc/model"
	"fmt"
	"log"
	"os"
)

func Capture(file *model.File) (err error) {
	if err = reName(file); err != nil {
		return
	}

	defer recoverName(file)

	err = getMediaInfo(file)
	if err != nil {
		return
	}

	err = CreateCaptures(*file)
	if err != nil {
		return
	}

	return
}

func reName(file *model.File) (err error) {
	err = os.Rename(file.Path, file.RePath)
	return
}

func recoverName(file *model.File) {
	err := os.Rename(file.RePath, file.Path)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("恢复名称失败：%v -> %v", file.ReName, file.Name)
	}
}
