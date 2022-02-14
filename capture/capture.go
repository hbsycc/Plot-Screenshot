package capture

import (
	"a.resources.cc/config"
	"a.resources.cc/lib"
	"a.resources.cc/model"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/fogleman/gg"
	"image/color"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
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

	err = createCaptures(*file)
	if err != nil {
		return
	}

	spew.Dump(file)

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

//
// createCaptures
// @Description: 创建临时目录、生成截图文件
// @param dir
// @param fileName
// @param durationSeconds
// @return err
//
func createCaptures(file model.File) (err error) {
	lib.DebugLog("开始截图\n", "ffmpeg")
	start := time.Now()

	// 创建临时目录
	_, err = os.Stat(config.GetConfig().Capture.Dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(config.GetConfig().Capture.Dir, 777)
		return
	}
	tempDir, err := os.MkdirTemp(config.GetConfig().Capture.Dir, "*")
	lib.DebugLog(fmt.Sprintf("创建临时目录:%v\n", tempDir), "dir")
	if err != nil {
		return
	}

	// 准备命令数组
	commands := make([]model.Capture, config.GetConfig().Capture.Count)
	for i := 0; i < config.GetConfig().Capture.Count; i++ {
		output := fmt.Sprintf("%v\\%v.jpg", tempDir, i+1)

		captureTime := file.MediaInfo.DurationSeconds / int64(config.GetConfig().Capture.Count) * int64(i)
		item := model.Capture{
			Command:   fmt.Sprintf("ffmpeg -ss %v -i %v -f image2 -y -frames:v 1 %v", captureTime, file.RePath, output),
			TimeStamp: captureTime,
			Image:     output,
		}

		stamp := strconv.FormatInt(captureTime, 10)
		if duration, err := time.ParseDuration(stamp + "s"); err != nil {
			return err
		} else {
			item.TimeDuration = fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)
		}

		commands[i] = item
	}

	// 多线程执行任务
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	processCount := 0
	for {
		command := commands[0]
		commands = commands[1:]

		processCount += 1
		wg.Add(1)
		go func() {
			e := ffmpegCaptures(ctx, command.Command)
			// 发生错误，通知所有协程取消
			if e != nil {
				cancel()
				err = e
			}

			// 当前截图写入时间
			err = drawTime(command)
			if err != nil {
				cancel()
				err = e
			}

			wg.Done()
		}()
		// 等待当前线程组的所有任务结束，发生错误的话直接返回
		if processCount == config.GetConfig().Capture.Thread {
			wg.Wait()
			processCount = 0
			if err != nil {
				return
			}
		}

		if len(commands) == 0 {
			break
		}
	}

	lib.DebugLog(fmt.Sprintf("截图完成：%v,耗时:%v\n", tempDir, time.Since(start)), "ffmpeg")

	// 删除临时文件夹
	if !config.GetConfig().Debug {
		err = os.RemoveAll(tempDir)
		if err != nil {
			return
		}
	}

	return
}

//
//  ffmpegCaptures
//  @Description: ffmpeg生成截图
//  @param ctx
//  @param commandStr
//  @return err
//
func ffmpegCaptures(ctx context.Context, commandStr string) (err error) {
	select {
	case <-ctx.Done():
		//fmt.Println("主程通知取消")
		return
	default:
		//fmt.Println("协程继续运行")
	}

	command := exec.Command("cmd", "/C", commandStr)
	lib.DebugLog(fmt.Sprintf("执行命令:%v\n", commandStr), "ffmpeg")
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}

	err = command.Run()
	if err != nil {
		msg := command.Stderr.(*bytes.Buffer).String()
		msg = lib.ConvertCodeToUTF8(msg, lib.GB18030)
		return errors.New(msg)
	}

	return
}

//
//  drawTime
//  @Description: 截图写入截取时间
//  @param imageSource
//  @param drawTime
//  @return err
//
func drawTime(capture model.Capture) (err error) {
	lib.DebugLog(fmt.Sprintf("截图文件:%v,截取时间:%v\n", capture.Image, capture.TimeDuration), "TimeDuration")

	image, err := gg.LoadImage(capture.Image)
	if err != nil {
		return
	}

	width := image.Bounds().Dx()
	height := image.Bounds().Dy()
	fontSize := width / 100 * 3

	dc := gg.NewContextForImage(image)
	err = dc.LoadFontFace("C:\\Windows\\Fonts\\Arial.ttf", float64(fontSize))
	if err != nil {
		return err
	}
	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 168})
	stringWidth, _ := dc.MeasureString(capture.TimeDuration)
	dc.DrawString(capture.TimeDuration, float64(width-int(stringWidth)-15), float64(height-10))
	err = dc.SavePNG(capture.Image)

	return
}

// mergeCaptures
// @Description: 合并截图
func mergeCaptures() {

}
