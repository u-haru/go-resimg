package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

var (
	source    = ""
	dest      = ""
	imagepath []string
	size      = 100
	ps        = 100
	lc        *sync.Mutex
)

func main() {
	flag.IntVar(&size, "s", 100, "-s [image-size]")
	flag.IntVar(&ps, "p", 12, "-p [process-num]")
	flag.Parse()
	source = filepath.Clean(flag.Arg(0))
	dest = filepath.Clean(flag.Arg(1))
	wg := &sync.WaitGroup{}
	lc = new(sync.Mutex)

	rand.Seed(time.Now().UnixNano())
	if flag.NArg() < 2 {
		fmt.Printf("usage:\n %s [Option] dir target\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	dirwalk(source)

	for i := 0; i < ps; i++ {
		wg.Add(1)
		go func() {
			converter()
			wg.Done()
		}()
	}
	wg.Wait()
}

func converter() {
	for {
		srcfile_path := pop()
		if srcfile_path == "" {
			return
		}
		dstfile_path := strings.ReplaceAll(srcfile_path, source, dest)
		extpos := strings.LastIndex(dstfile_path, ".")
		dstfile_path = dstfile_path[:extpos] + ".jpg" //拡張子固定

		if checkFile(dstfile_path) {
			continue
		}

		srcfile, err := os.Open(srcfile_path)
		if err != nil {
			fmt.Println("Error: File could not be opened")
			continue
		}

		os.MkdirAll(filepath.Dir(dstfile_path), 755) //ディレクトリ作成

		dst, err := os.Create(dstfile_path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		img, _, err := image.Decode(srcfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		imgdst := resize(img, 100)
		jpeg.Encode(dst, imgdst, &jpeg.Options{Quality: 100})
		srcfile.Close()
		dst.Close()
	}
}

func pop() string {
	lc.Lock()
	defer lc.Unlock()

	lastindex := len(imagepath)
	if lastindex == 0 {
		return ""
	}
	str := imagepath[lastindex-1]
	imagepath = imagepath[:lastindex-1]
	return str
}

func checkFile(filename string) bool {
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		return true
	}
	return false
}

func resize(img image.Image, s int) image.Image {
	rct := img.Bounds()
	m := float64(rct.Dy()) / float64(rct.Dx()) //アスペクト比
	w := s
	h := s
	if m > 1 { // s * s に収める
		w = int(float64(s) / m)
	} else {
		h = int(float64(s) * m)
	}

	src := image.NewRGBA(image.Rect(0, 0, w, h))
	dst := image.NewRGBA(image.Rect(0, 0, s, s))

	c := color.RGBA{0, 0, 0, 255}                                       //黒ペン
	draw.Draw(dst, dst.Bounds(), &image.Uniform{c}, image.ZP, draw.Src) //黒塗り

	draw.CatmullRom.Scale(src, src.Bounds(), img, rct, draw.Over, nil) //リサイズ

	dstrct := image.Rectangle{image.Point{(s - w) / 2, (s - h) / 2}, image.Point{(s + w) / 2, (s + h) / 2}} //範囲指定
	draw.Draw(dst, dstrct, src, image.Point{0, 0}, draw.Src)                                                //黒塗りと合成
	return dst
}

func dirwalk(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, file_info := range files {
		if file_info.IsDir() {
			path := filepath.Join(dir, file_info.Name())
			path = filepath.Clean(path)
			dirwalk(path)
			continue
		}
		if checkext(file_info.Name()) != true {
			continue
		}
		path := filepath.Join(dir, file_info.Name())
		path = filepath.Clean(path)
		imagepath = append(imagepath, path)
	}
}

var media_ext = []string{"png", "jpg", "jpeg", "gif", "webm", "webp", "bmp"}

func checkext(filename string) bool {
	filename = strings.ToLower(filename)
	for _, s := range media_ext {
		if strings.HasSuffix(filename, s) {
			return true
		}
	}
	return false
}
