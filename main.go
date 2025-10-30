package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sqweek/dialog"
	"github.com/tidwall/gjson"
)

func main() {
	path, err := dialog.File().Title("파일 선택").Filter("All Files", "*").Load()
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !isZipData(data) {
		fmt.Println("ZIP 형식 아님")
		return
	}
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		fmt.Println(err)
		return
	}

	var content []byte
	for _, f := range r.File {
		if f.Name == "card.json" {
			rc, _ := f.Open()
			content, _ = io.ReadAll(rc)
			break
		}
	}
	result := gjson.GetBytes(content, "data.assets")
	arr := result.Array()

	nameMap := make(map[string]string)
	for _, item := range arr {
		name := item.Get("name").String()
		splitUri := strings.Split(item.Get("uri").String(), "/")
		uri := splitUri[len(splitUri)-1]
		nameMap[uri] = name
	}
	for _, f := range r.File {
		splitFilename := strings.Split(f.Name, "/")
		filename := splitFilename[len(splitFilename)-1]
		if name, ok := nameMap[filename]; ok {
			fmt.Printf("✅ %s -> %s\n", filename, name)
		}
	}
}

func isZipData(data []byte) bool {
	_, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	return err == nil
}
