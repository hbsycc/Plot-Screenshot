package capture

import (
	"a.resources.cc/config"
	"a.resources.cc/font"
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
	"math"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func Capture(file *model.File) (err error) {
	// 创建监听系统信号 channel
	c := make(chan os.Signal)
	// 监听信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	// 协程监听退出
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
				RecoverName(*file)
				os.Exit(0)
			default:
				fmt.Println("other signal", s)
			}
		}
	}()

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

	// 删除临时文件夹
	if !config.GetConfig().Debug {
		err = os.RemoveAll(file.TempDir)
		if err != nil {
			return
		}
	}

	return
}

func TempName(file model.File) (err error) {
	err = os.Rename(file.Path, file.RePath)
	return
}

func RecoverName(file model.File) {
	err := os.Rename(file.RePath, file.Path)
	if err != nil {
		log.Fatalf("恢复名称失败：%v -> %v", file.RePath, file.Path)
		return
	}

	return
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
		err = os.MkdirAll(file.TempDir, os.ModePerm)
		if err != nil {
			fmt.Println("这里", file.TempDir)
			return
		}
	}

	// 准备命令数组
	captureTotal := config.GetConfig().Capture.Grid.Row * config.GetConfig().Capture.Grid.Column
	commands := make([]model.Capture, captureTotal)
	for i := 0; i < captureTotal; i++ {
		output := fmt.Sprintf("%v%v.jpg", file.TempDir, i+1)
		if output, err = filepath.Rel(file.Dir, file.TempDir); err != nil {
			return err
		} else {
			output = fmt.Sprintf("%v%v%v.jpg", output, string(filepath.Separator), i+1)
		}

		captureTime := file.MediaInfo.DurationSeconds / int64(captureTotal) * int64(i)
		item := model.Capture{
			Command:   fmt.Sprintf("ffmpeg -ss %v -i %v -f image2 -y -frames:v 1 %v", captureTime, file.XxHash+file.Ext, output),
			TimeStamp: captureTime,
			Image:     filepath.Join(file.Dir, output),
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
			e := ffmpegCaptures(ctx, file.Dir, command.Command)
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

	return
}

// ffmpegCaptures
// @Description: ffmpeg生成截图
// @param ctx
// @param mediaDir
// @param commandStr
// @return err
func ffmpegCaptures(ctx context.Context, mediaDir string, commandStr string) (err error) {
	select {
	case <-ctx.Done():
		fmt.Println("主程通知取消")
		return
	default:
		//fmt.Println("协程继续运行")
	}

	command := exec.Command("cmd", "/C", commandStr)
	lib.DebugLog(fmt.Sprintf("执行命令:%v", commandStr), "ffmpeg")
	command.Dir = mediaDir
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

// drawTime
//  @Description: 截图写入截取时间、描边
// @param capture
// @return err
func drawTime(capture model.Capture) (err error) {
	lib.DebugLog(fmt.Sprintf("截图文件:%v,写入截取时间:%v", capture.Image, capture.TimeDuration), "TimeDuration")

	img, err := gg.LoadImage(capture.Image)
	if err != nil {
		return
	}

	width := float64(img.Bounds().Dx())
	height := float64(img.Bounds().Dy())
	fontSize := float64(width) * 0.04
	dc := gg.NewContextForImage(img)
	if face, err := font.GetFontFace(fontSize); err != nil {
		return err
	} else {
		dc.SetFontFace(face)
		dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 168})
		stringWidth, _ := dc.MeasureString(capture.TimeDuration)
		dc.DrawString(capture.TimeDuration, width-stringWidth, height-5)
	}

	if config.GetConfig().Capture.Grid.BorderWidth > 0 {
		dc.DrawRectangle(0, 0, width, height)
		dc.SetColor(color.Black)
		dc.SetLineWidth(config.GetConfig().Capture.Grid.BorderWidth)
		dc.Stroke()
	}
	err = gg.SaveJPG(capture.Image, dc.Image(), 100)
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

	// 截图区域
	var (
		column    = config.GetConfig().Capture.Grid.Column
		columnGap = config.GetConfig().Capture.Grid.ColumnGap
		row       = config.GetConfig().Capture.Grid.Row
		rowGap    = config.GetConfig().Capture.Grid.RowGap
		bgWidth   = file.MediaInfo.Video.Width*column + (column-1)*columnGap
		bgHeight  = file.MediaInfo.Video.Height*row + (row-1)*rowGap
	)
	rect := image.Rect(0, 0, bgWidth, bgHeight)
	bg := image.NewRGBA(rect)
	draw.Draw(bg, rect.Bounds(), &image.Uniform{C: color.White}, image.Pt(0, 0), draw.Src)
	dc := gg.NewContextForRGBA(bg)
	for i, cn := range capturesName {
		var c image.Image
		path := fmt.Sprintf("%v//%v", file.TempDir, cn)
		if c, err = gg.LoadImage(path); err != nil {
			return
		}

		xIndex := i % column
		x := (xIndex * columnGap) + (xIndex * file.MediaInfo.Video.Width)
		yIndex := i / column
		y := (yIndex * rowGap) + (yIndex * file.MediaInfo.Video.Height)
		dc.DrawImage(c, x, y)
	}
	capturesImage := dc.Image()

	// 信息区域
	metas := []string{
		fmt.Sprint(file.MediaInfo.Video.Width, "*", file.MediaInfo.Video.Height),
		file.MediaInfo.Video.DisplayAspectRatio,
		file.MediaInfo.Video.CodecName,
		file.MediaInfo.Video.PixFmt,
	}
	fileSize, err := strconv.ParseInt(file.MediaInfo.Format.Size, 10, 64)
	if err != nil {
		return
	}
	drawStrings := []string{
		fmt.Sprintf("文件名称 ：%v", file.Name),
		fmt.Sprintf("文件大小 ：%v", lib.FormatFileSize(fileSize)),
		fmt.Sprintf("播放时长 ：%v", file.MediaInfo.DurationFormat),
		fmt.Sprintf("编码信息 ：%v", strings.Join(metas, " / ")),
		fmt.Sprintf("文件Hash：%v", file.XxHash),
	}
	var fontSize = math.Ceil(float64(bgWidth) * 0.013)
	lineHeight := math.Ceil(fontSize * 1.5)
	metaHeight := int(lineHeight)*len(drawStrings) + int(lineHeight*0.5)
	rect = image.Rect(0, 0, bgWidth, bgHeight+metaHeight)
	bg = image.NewRGBA(rect)
	dc = gg.NewContextForRGBA(bg)
	draw.Draw(bg, rect.Bounds(), &image.Uniform{C: color.White}, image.Pt(0, 0), draw.Src)
	dc.DrawImage(capturesImage, 0, metaHeight)
	if fontFace, err := font.GetFontFace(fontSize); err != nil {
		return err
	} else {
		dc.SetFontFace(fontFace)
		dc.SetRGB(0, 0, 0)
		textIndent, _ := dc.MeasureString("A")
		for i, s := range drawStrings {
			dc.DrawString(s, textIndent, lineHeight*float64(i+1))
		}
	}

	// 缩放、保存图片
	o := strings.Split(file.TempDir, string(filepath.Separator))
	out := strings.Join(o[0:len(o)-1], string(filepath.Separator))
	outFile := strings.ReplaceAll(file.Name, file.Ext, ".jpg")
	out = filepath.Join(out, outFile)
	outImage := dc.Image()
	resizeWidth := config.GetConfig().Capture.ResizeWidth
	if resizeWidth > 0 {
		outImage = imaging.Resize(outImage, resizeWidth, 0, imaging.Lanczos)
	}
	err = gg.SaveJPG(out, outImage, config.GetConfig().Capture.Quality)
	lib.DebugLog(fmt.Sprintf("合成完成：%v,耗时:%v", out, time.Since(startTime)), "merge")

	return
}
