package capture

import (
	"a.resources.cc/lib"
	"a.resources.cc/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"time"
)

//
//  getMediaInfo
//  @Description: 获取媒体元信息
//  @param file
//  @return err
//
func getMediaInfo(file *model.File) (err error) {
	commandStr := fmt.Sprintf("ffprobe -print_format json  -show_streams %v", file.RePath)
	fmt.Println(file.RePath)
	fmt.Println(commandStr)
	lib.DebugLog(fmt.Sprintf("执行命令:%v\n", commandStr), "ffprobe")

	command := exec.Command("cmd", "/C", commandStr)
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}
	err = command.Run()
	if err != nil {
		msg := command.Stderr.(*bytes.Buffer).String()
		msg = lib.ConvertCodeToUTF8(msg, lib.GB18030)
		return errors.New(msg)
	}

	ffProbe := &model.FFProbe{}
	if err = json.Unmarshal(command.Stdout.(*bytes.Buffer).Bytes(), ffProbe); err != nil {
		return
	}
	var float float64
	if float, err = strconv.ParseFloat(ffProbe.Streams[0].Duration, 64); err != nil {
		return
	} else {
		file.MediaInfo.DurationSeconds = int64(math.Floor(float))
		file.MediaInfo.Width = ffProbe.Streams[0].Width
		file.MediaInfo.Height = ffProbe.Streams[0].Height
		file.MediaInfo.DisplayAspectRatio = ffProbe.Streams[0].DisplayAspectRatio
		file.MediaInfo.CodecName = ffProbe.Streams[0].CodecName
		file.MediaInfo.PixFmt = ffProbe.Streams[0].PixFmt
		file.MediaInfo.RFrameRate = ffProbe.Streams[0].RFrameRate

		if duration, err := time.ParseDuration(strconv.FormatInt(file.MediaInfo.DurationSeconds, 10) + "s"); err != nil {
			return err
		} else {
			df := fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
			file.MediaInfo.DurationFormat = df
		}

	}

	return
}
