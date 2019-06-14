package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

var html []byte

// HTML をブラウザに送信
func handlerHtml(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	_, _ = w.Write(html)
}

func handlerPrimeSSE(w http.ResponseWriter, _ *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	closeNotify := w.(http.CloseNotifier).CloseNotify()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var num int64 = 1
	for id := 1; id <= 100; id++ {
		// 通信が切れても終了
		select {
		case <-closeNotify:
			fmt.Println("Connection closed from client")
			return
		default:
			// do nothing
		}
		for {
			num++
			// 確率論的に素数を求める
			if big.NewInt(num).ProbablyPrime(20) {
				fmt.Println(num)
				_, _ = fmt.Fprintf(w, "data: {\"id\": %d, \"number\": %d}\n\n", id, num)
				flusher.Flush()
				time.Sleep(time.Second)
				break
			}
		}
		time.Sleep(time.Second)
	}
	// 100 個超えたら送信終了
	fmt.Println("Connection closed from server")
}

func main() {
	var err error
	html, err = ioutil.ReadFile("./chapter09/9-5_sse_view/index.html")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", handlerHtml)
	http.HandleFunc("/prime", handlerPrimeSSE)
	fmt.Println("start http listening :18888")
	err = http.ListenAndServe(":18888", nil)
	fmt.Println(err)
}
