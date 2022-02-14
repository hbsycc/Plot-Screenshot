package demo

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"image/draw"
	"os/exec"
	"syscall"
)

func demo1() {
	inBuffer := bytes.NewBuffer(nil)

	//cmd := exec.Command("sh") //linux
	cmd := exec.Command("cmd") //windows
	cmd.Stdin = inBuffer
	go func() {
		inBuffer.WriteString(`avifenc --jobs 8 --depth 12 "D:\input.png" "D:\output.avif" >> test.txt`)
		inBuffer.WriteString("\n")
		inBuffer.WriteString("exit")
		inBuffer.WriteString("\n")
	}()
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}
}

func demo2() {
	cmd := exec.Command("cmd")
	inBuffer := bytes.Buffer{}
	outBuffer := bytes.Buffer{}
	cmd.Stdout = &outBuffer
	cmd.Stdin = &inBuffer
	inBuffer.WriteString("CHCP 65001" + "\n")
	inBuffer.WriteString(`avifenc --jobs 8 --depth 12 "D:\input.png" "D:\output.avif" >> test.txt` + "\n")
	inBuffer.WriteString("exit" + "\n")

	err := cmd.Start()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err.Error())
	}

	// 程序Pid
	fmt.Println(cmd.ProcessState.Pid())
	// 程序退出code
	fmt.Println(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitCode)
	// 输出结果
	fmt.Println(outBuffer.String())
}

func demo3() {
	var err error

	im1, err := gg.LoadImage("innput.jpg")
	if err != nil {
		panic(err)
	}

	im2, err := gg.LoadImage("innput.jpg")
	if err != nil {
		panic(err)
	}

	s1 := im1.Bounds().Size()

	const (
		width  = 1024
		height = 1024
	)
	rect := image.Rect(0, 0, width, height)
	bg := image.NewRGBA(rect)
	draw.Draw(bg, rect.Bounds(), &image.Uniform{C: color.White}, image.Pt(500, 500), draw.Src)
	dc := gg.NewContextForRGBA(bg)
	dc.DrawImage(im1, 0, 0)
	dc.DrawImage(im2, 0, s1.Y)
	err = dc.SavePNG("output.png")
	if err != nil {
		panic(err)
	}
}
