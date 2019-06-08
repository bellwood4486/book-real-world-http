package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var image []byte

func init() {
	var err error
	image, err = ioutil.ReadFile("./chapter09/9-3_server_push/image.png")
	if err != nil {
		panic(err)
	}
}

// HTML をブラウザに送信
// 画像をプッシュする
func handlerHtml(w http.ResponseWriter, r *http.Request) {
	// Pusher にキャスト可能であれば (HTTP/2で接続していたら) プッシュする
	pusher, ok := w.(http.Pusher)
	if ok {
		_ = pusher.Push("/image", nil)
	}
	w.Header().Add("Content-Type", "text/html")
	_, _ = fmt.Fprintf(w, `<html><body><img src="/image"></body></html>`)
}

// 画像ファイルをブラウザに送信
func handlerImage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	_, _ = w.Write(image)
}

func main() {
	http.HandleFunc("/", handlerHtml)
	http.HandleFunc("/image", handlerImage)
	fmt.Println("start http listening :18443")
	err := http.ListenAndServeTLS(":18443", "./chapter06/server.crt", "./chapter06/server.key", nil)
	fmt.Println(err)
}
