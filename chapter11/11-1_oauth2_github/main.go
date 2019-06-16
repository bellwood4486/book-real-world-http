package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var clientID = os.Getenv("GITHUB_CLIENT_ID")
var clientSecret = os.Getenv("GITHUB_CLIENT_SECRET")

//var redirectURL = "http://localhost:18888"
var state = "your state"

const AccessTokenCacheFile = "./chapter11/11-1_oauth2_github/access_token.json"

func main() {
	// OAuth2の接続先などの情報
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user:email", "gist"},
		Endpoint:     github.Endpoint,
	}
	// これをこれから初期化する
	var token *oauth2.Token

	// ローカルにすでに保存済み？
	file, err := os.Open(AccessTokenCacheFile)
	if os.IsNotExist(err) {
		// 初回アクセス
		// まず認可画面のURLを取得
		url := conf.AuthCodeURL(state, oauth2.AccessTypeOnline)

		// コールバックを受け取るウェブサーバーをセットアップ
		code := make(chan string)
		var server *http.Server
		server = &http.Server{
			Addr: ":18888",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// クエリーパラメータからcodeを取得し、ブラウザを閉じる
				w.Header().Set("Content-Type", "text/html")
				_, _ = io.WriteString(w, "<html><script>window.open('about:blank', '_self').close()</script></html>")
				w.(http.Flusher).Flush()
				code <- r.URL.Query().Get("code")
				// サーバーも閉じる
				_ = server.Shutdown(context.Background())
			}),
		}
		go server.ListenAndServe()

		// ブラウザで認可画面を開く
		// GitHubの認可が完了すれば上記のサーバーにリダイレクト
		// されて、Handlerが実行される
		_ = open.Start(url)

		// 取得したコードをアクセストークンに交換
		token, err = conf.Exchange(context.TODO(), <-code)
		if err != nil {
			panic(err)
		}

		// アクセストークンをファイルに保存
		file, err := os.Create(AccessTokenCacheFile)
		if err != nil {
			panic(err)
		}
		_ = json.NewEncoder(file).Encode(token)
	} else if err == nil {
		// 一度認可をしてローカルに保存済み
		token = &oauth2.Token{}
		_ = json.NewDecoder(file).Decode(token)
	} else {
		panic(err)
	}
	client := oauth2.NewClient(context.TODO(), conf.TokenSource(context.TODO(), token))

	// Email 取得(11-3)
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	emails, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(emails))

	// gistに投稿(11-4)
	gist := `{
      "description": "API example",
      "public": true,
      "files": {
        "hello_from_rest_api.txt": {
          "content": "Hello World"
        }
      }
    }`
	// 投稿
	resp2, err := client.Post("https://api.github.com/gists", "application/json", strings.NewReader(gist))
	if err != nil {
		panic(err)
	}
	fmt.Println(resp2.Status)
	defer resp2.Body.Close()
	// 結果をパースする
	type GistResult struct {
		Url string `json:"html_url"`
	}
	gistResult := &GistResult{}
	err = json.NewDecoder(resp2.Body).Decode(gistResult)
	if err != nil {
		panic(err)
	}
	if gistResult.Url != "" {
		// ブラウザで開く
		_ = open.Start(gistResult.Url)
	}
}
