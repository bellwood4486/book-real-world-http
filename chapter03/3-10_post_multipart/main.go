package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	_ = writer.WriteField("name", "Michael Jackson")

	fileWriter, err := writer.CreateFormFile("thumbnail", "photo.png")
	if err != nil {
		panic(err)
	}
	readFile, err := os.Open("./chapter03/3-10_post_multipart/photo.png")
	if err != nil {
		panic(err)
	}
	defer readFile.Close()
	_, _ = io.Copy(fileWriter, readFile)
	_ = writer.Close()

	resp, err := http.Post("http://localhost:18888", writer.FormDataContentType(), &buffer)
	if err != nil {
		panic(err)
	}
	log.Print("Status:", resp.Status)
}
