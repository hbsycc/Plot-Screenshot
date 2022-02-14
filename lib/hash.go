package lib

import (
	"fmt"
	"github.com/codingsince1985/checksum"
	"time"
)

const p = `D:\Workspaces\webDav\8T\test`
const f = `input.avi`

func Hash() {
	startTime := time.Now()
	if md5, err := checksum.MD5sum(p + "/" + f); err != nil {
		panic(err)
	} else {
		fmt.Printf("MD5:%v,耗时:%v\n", md5, time.Since(startTime))
	}

	startTime = time.Now()
	if sha1, err := checksum.SHA1sum(p + "/" + f); err != nil {
		panic(err)
	} else {
		fmt.Printf("SHA1:%v,耗时:%v\n", sha1, time.Since(startTime))
	}

	startTime = time.Now()
	if sha256, err := checksum.SHA256sum(p + "/" + f); err != nil {
		panic(err)
	} else {
		fmt.Printf("SHA256:%v,耗时:%v\n", sha256, time.Since(startTime))
	}
}
