package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/BenLubar/vpk"
)

var rules = map[string]func(map[string]string) string{
	"abilities_japanese":       defaultBuilder,
	"broadcastfacts_japanese":  defaultBuilder,
	"chat_japanese":            defaultBuilder,
	"dota_japanese":            defaultBuilder,
	"gameui_japanese":          defaultBuilder,
	"hero_chat_wheel_japanese": simpleBuilder("hero_chat_wheel"),
	"hero_lore_japanese":       defaultBuilder,
	"leagues_japanese":         simpleBuilder("leagues"),
	"richpresence_japanese":    defaultBuilder,
}

func main() {
	var force bool
	flag.BoolVar(&force, "force", false, "ignore check current directory")
	flag.Parse()

	wd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	if filepath.Base(wd) != "dota 2 beta" && !force {
		log.Println("僕の配置場所が間違ってるっぽいよ?")
		log.Println("ヒント: 'dota 2 beta'フォルダの真下においてね!")
		log.Println("--force で無視できるよ")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}

	log.Println("ダウンロード中…")
	resp, err := http.Get("https://nihongoka.games/download/dota2/?format=zip")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	log.Println("ダウンロード完了!")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	os.MkdirAll("game/dota_japanese", os.ModePerm)

	vpkc := vpk.SingleVPKCreator("game/dota_japanese/pak01_dir.vpk")

	var contents []vpk.Entry

	zip, _ := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	for _, file := range zip.File {
		name := path.Base(file.Name)
		woeName := name[:len(name)-len(path.Ext(name))]
		log.Println(woeName)
		builder, ok := rules[woeName]
		if !ok {
			continue
		}
		f, err := file.Open()
		if err != nil {
			panic(err)
		}
		defer f.Close()
		var data map[string]string
		if err := json.NewDecoder(f).Decode(&data); err != nil {
			panic(err)
		}
		contents = append(contents, entry{
			rel:  "resource/localization/" + woeName + ".txt",
			data: builder(data),
		})
	}
	log.Println("書き出し中…")
	if err := vpk.Create(vpkc, contents, -1); err != nil {
		panic(err)
	}
	log.Println("完了!!")
	time.Sleep(time.Second * 3)
}

func defaultBuilder(data map[string]string) string {
	b := new(strings.Builder)
	fmt.Fprintln(b, `"lang"`)
	fmt.Fprintln(b, `{`)
	fmt.Fprintln(b, `	"Language" "japanese"`)
	fmt.Fprintln(b, `	"Tokens"`)
	fmt.Fprintln(b, `	{`)
	for k, v := range data {
		fmt.Fprintf(b, `		"%v" "%v"`, k, escape(v))
		fmt.Fprintln(b)
	}
	fmt.Fprintln(b, `	}`)
	fmt.Fprintln(b, `}`)
	return b.String()
}

func simpleBuilder(key string) func(map[string]string) string {
	return func(data map[string]string) string {
		b := new(strings.Builder)
		fmt.Fprintf(b, `"%v"`, key)
		fmt.Fprintln(b)
		fmt.Fprintln(b, `{`)
		for k, v := range data {
			fmt.Fprintf(b, `	"%v" "%v"`, k, escape(v))
			fmt.Fprintln(b)
		}
		fmt.Fprintln(b, `}`)
		return b.String()
	}
}

func escape(s string) string {
	vv := [][2]string{
		{"\\", "\\\\"},
		{"?", "\\?"},
		{"\n", "\\n"},
		{"\t", "\\t"},
		{"\"", "\\\""},
	}
	for _, v := range vv {
		s = strings.ReplaceAll(s, v[0], v[1])
	}
	return s
}

type entry struct {
	rel  string
	data string
}

func (e entry) Rel() string {
	return e.rel
}

func (e entry) Open() (io.ReadCloser, error) {
	return ioutil.NopCloser(strings.NewReader(e.data)), nil
}
