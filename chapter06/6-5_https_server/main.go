package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	}
	fmt.Println(string(dump))
	_, _ = fmt.Fprintf(w, "<html><body>hello</body></html>")
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("start http listening :18443")
	err := http.ListenAndServeTLS(":18443",
		"./chapter06/6-5_https_server/server.crt",
		"./chapter06/6-5_https_server/server.key",
		nil)
	log.Println(err)
}
