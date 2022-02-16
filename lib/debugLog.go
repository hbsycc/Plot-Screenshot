package lib

import (
	"a.resources.cc/config"
	"fmt"
	"log"
)

func DebugLog(msg string, prefix ...string) {
	if !config.GetConfig().Debug {
		return
	}

	for _, flag := range prefix {
		log.SetPrefix(fmt.Sprintf("[%v] ", flag))
	}

	log.SetFlags(log.Ltime)
	log.Println(msg)
}
