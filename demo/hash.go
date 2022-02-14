package demo

import (
	"fmt"
	"github.com/OneOfOne/xxhash"
	xxhashv2 "github.com/cespare/xxhash/v2"
	"github.com/codingsince1985/checksum"
	"hash/crc32"
	"hash/crc64"
	"io"
	"os"
	"time"
)

func Hash() {
	var startTime time.Time
	file := `E:\windows\cn_windows_10_enterprise_ltsc_2019_x64_dvd_9c09ff24.iso`
	//file := `E:\第五届索尔维会议\Solvay_conference_1927.jpg`

	//startTime := time.Now()
	//if md5, err := checksum.MD5sum(file); err != nil {
	//	panic(err)
	//} else {
	//	fmt.Printf("MD5:%v,耗时:%v\n", md5, time.Since(startTime))
	//}

	startTime = time.Now()
	if sha1, err := checksum.SHA1sum(file); err != nil {
		panic(err)
	} else {
		fmt.Printf("SHA1:%v,耗时:%v\n", sha1, time.Since(startTime))
	}

	//startTime = time.Now()
	//if sha256, err := checksum.SHA256sum(file); err != nil {
	//	panic(err)
	//} else {
	//	fmt.Printf("SHA256:%v,耗时:%v\n", sha256, time.Since(startTime))
	//}
	//
	startTime = time.Now()
	if crc32, err := checksum.CRC32(file); err != nil {
		panic(err)
	} else {
		fmt.Printf("CRC32:%v,耗时:%v\n", crc32, time.Since(startTime))
	}

	f1, _ := os.Open(file)
	defer f1.Close()
	startTime = time.Now()
	h := xxhash.New64()
	io.Copy(h, f1)
	fmt.Printf("xxhash:%v,耗时:%v\n", h.Sum64(), time.Since(startTime))

	f12, _ := os.Open(file)
	defer f12.Close()
	h2 := xxhashv2.New()
	io.Copy(h2, f12)
	fmt.Printf("xxhash:%v,耗时:%v\n", h2.Sum64(), time.Since(startTime))

	f11, _ := os.Open(file)
	defer f11.Close()
	startTime = time.Now()
	crc32Hash := crc32.New(crc32.MakeTable(crc32.IEEE))
	io.Copy(crc32Hash, f11)
	fmt.Printf("crc32:%v,耗时:%v\n", crc32Hash.Sum32(), time.Since(startTime))

	f2, _ := os.Open(file)
	defer f2.Close()
	startTime = time.Now()
	t := crc64.New(crc64.MakeTable(crc64.ECMA))
	_, _ = io.Copy(t, f2)
	fmt.Printf("crc64:%v,耗时:%v\n", t.Sum64(), time.Since(startTime))

}
