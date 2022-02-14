package lib

import (
	"a.resources.cc/config"
	"fmt"
	"strings"
)

func DebugLog(msg string, flags ...string) {
	if !config.GetConfig().Debug {
		return
	}

	flagStr := ""
	for _, flag := range flags {
		strings.Join(flags, "")
		flagStr += fmt.Sprintf("[%v]", flag)
	}

	fmt.Printf("%v %v", flagStr, msg)
}
