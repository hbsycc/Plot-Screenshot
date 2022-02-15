package demo

import (
	"encoding/hex"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"github.com/minio/highwayhash"
	"hash/crc64"
	"io"
	"os"
	"time"
)

const count = 10

func Hash() {
	var startTime time.Time
	file := `E:\windows\cn_windows_10_enterprise_ltsc_2019_x64_dvd_9c09ff24.iso`

	// highwayHash
	startTime = time.Now()
	for i := 0; i < count; i++ {
		key, err := hex.DecodeString("000102030405060708090A0B0C0D0E0FF0E0D0C0B0A090807060504030201000") // use your own key here
		if err != nil {
			fmt.Printf("Cannot decode hex key: %v", err) // add error handling
			return
		}
		hh, err := highwayhash.New64(key)
		if err != nil {
			fmt.Println(err)
		}
		f11, _ := os.Open(file)
		_, _ = io.Copy(hh, f11)

		if i+1 == count {
			fmt.Printf("highwayhash:%v,平均耗时:%v\n", hh.Sum64(), time.Since(startTime)/count)
		}
		_ = f11.Close()
	}

	// xxHash
	startTime = time.Now()
	for i := 0; i < count; i++ {
		f12, _ := os.Open(file)
		h := xxhash.New()
		_, _ = io.Copy(h, f12)
		if i+1 == count {
			fmt.Printf("xxhash:%v,平均耗时:%v\n", h.Sum64(), time.Since(startTime)/count)
		}
		_ = f12.Close()
	}

	// CRC64
	startTime = time.Now()
	for i := 0; i < count; i++ {
		f2, _ := os.Open(file)
		h2 := crc64.New(crc64.MakeTable(crc64.ECMA))
		_, _ = io.Copy(h2, f2)
		if i+1 == count {
			fmt.Printf("crc64:%v,平均耗时:%v\n", h2.Sum64(), time.Since(startTime)/count)
		}
		_ = f2.Close()
	}
}
