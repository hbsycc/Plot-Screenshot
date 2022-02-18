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

// getMediaInfo
//  @Description: 获取媒体元信息
// @param file
// @return err
func getMediaInfo(file *model.File) (err error) {
	// ffprobe 命令行详解	https://blog.csdn.net/ssehs/article/details/106625342
	commandStr := fmt.Sprintf("ffprobe -show_format -show_streams -select_streams v -print_format json -show_entries format=filename,format_name,duration,format_long_name,size:stream=width,height,display_aspect_ratio,r_frame_rate,bit_rate,codec_name,pix_fmt,codec_long_name,avg_frame_rate -loglevel error %v", file.XxHash+file.Ext)
	lib.DebugLog(fmt.Sprintf("执行命令:%v", commandStr), "ffprobe")

	command := exec.Command("cmd", "/C", commandStr)
	command.Dir = file.Dir // 指定command的工作目录
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
	if float, err = strconv.ParseFloat(ffProbe.Format.Duration, 64); err != nil {
		return
	} else {
		file.MediaInfo.Video = *ffProbe.Streams[0]
		file.MediaInfo.Format = *ffProbe.Format
		file.MediaInfo.DurationSeconds = int64(math.Floor(float))
		if duration, err := time.ParseDuration(strconv.FormatInt(file.MediaInfo.DurationSeconds, 10) + "s"); err != nil {
			return err
		} else {
			df := fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
			file.MediaInfo.DurationFormat = df
		}
	}

	return
}
