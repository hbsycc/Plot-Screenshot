package capture

import (
	"a.resources.cc/config"
	"a.resources.cc/lib"
	"a.resources.cc/model"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
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

	err = createCaptures(file)
	if err != nil {
		return
	}

	err = mergeCaptures(file)
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
		log.Fatalf("恢复名称失败：%v -> %v", file.XxHash, file.Name)
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
func createCaptures(file *model.File) (err error) {
	lib.DebugLog("开始截图", "ffmpeg")
	start := time.Now()

	// 创建临时目录
	lib.DebugLog(fmt.Sprintf("创建临时目录(文件xxhash):%v", file.TempDir), "dir")
	_, err = os.Stat(file.TempDir)
	if err != nil && os.IsNotExist(err) {
		if os.IsNotExist(err) {
			err = os.MkdirAll(file.TempDir, 777)
		} else {
			return
		}
	}

	// 准备命令数组
	captureTotal := config.GetConfig().Capture.Grid.Row * config.GetConfig().Capture.Grid.Column
	commands := make([]model.Capture, captureTotal)
	for i := 0; i < captureTotal; i++ {
		output := fmt.Sprintf("%v\\%v.jpg", file.TempDir, i+1)

		captureTime := file.MediaInfo.DurationSeconds / int64(captureTotal) * int64(i)
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
		go func(file *model.File) {
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
		}(file)
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

	lib.DebugLog(fmt.Sprintf("截图完成：%v,耗时:%v", file.TempDir, time.Since(start)), "ffmpeg")

	// 删除临时文件夹
	if !config.GetConfig().Debug {
		err = os.RemoveAll(file.TempDir)
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
	lib.DebugLog(fmt.Sprintf("执行命令:%v", commandStr), "ffmpeg")
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
	lib.DebugLog(fmt.Sprintf("截图文件:%v,截取时间:%v", capture.Image, capture.TimeDuration), "TimeDuration")

	img, err := gg.LoadImage(capture.Image)
	if err != nil {
		return
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	fontSize := width / 100 * 3

	dc := gg.NewContextForImage(img)
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
func mergeCaptures(file *model.File) (err error) {
	lib.DebugLog("开始合成", "merge")
	startTime := time.Now()

	// 读取临时目录下截图文件、按文件名排序
	var capturesName []string
	read, err := os.ReadDir(file.TempDir)
	if err != nil {
		return
	}
	for _, entry := range read {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			capturesName = append(capturesName, info.Name())
		}
	}
	sort.SliceStable(capturesName, func(i, j int) bool {
		iName, _ := strconv.Atoi(strings.Split(capturesName[i], ".")[0])
		jName, _ := strconv.Atoi(strings.Split(capturesName[j], ".")[0])
		return iName < jName
	})

	// 合成大图
	var (
		column    = config.GetConfig().Capture.Grid.Column
		columnGap = config.GetConfig().Capture.Grid.ColumnGap
		row       = config.GetConfig().Capture.Grid.Row
		rowGap    = config.GetConfig().Capture.Grid.RowGap
		bgWidth   = file.MediaInfo.Width*column + (column-1)*columnGap
		bgHeight  = file.MediaInfo.Height*row + (row-1)*rowGap
	)
	rect := image.Rect(0, 0, bgWidth, bgHeight)
	bg := image.NewRGBA(rect)
	draw.Draw(bg, rect.Bounds(), &image.Uniform{C: color.White}, image.Pt(500, 500), draw.Src)
	dc := gg.NewContextForRGBA(bg)
	for i, cn := range capturesName {
		path := fmt.Sprintf("%v//%v", file.TempDir, cn)
		c, err := gg.LoadImage(path)
		if err != nil {
			return err
		}

		xIndex := i % column
		x := (xIndex * columnGap) + (xIndex * file.MediaInfo.Width)
		yIndex := i / column
		y := (yIndex * rowGap) + (yIndex * file.MediaInfo.Height)
		dc.DrawImage(c, x, y)
	}

	out := fmt.Sprintf("%v\\%v.jpg", config.GetConfig().Capture.Dir, file.XxHash)
	zoom := imaging.Resize(dc.Image(), 4096, 0, imaging.Lanczos)
	err = gg.SaveJPG(out, zoom, config.GetConfig().Capture.Quality)

	lib.DebugLog(fmt.Sprintf("合成完成：%v,耗时:%v", out, time.Since(startTime)), "merge")

	return
}
