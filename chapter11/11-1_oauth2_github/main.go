package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io"
	"net/http"
	"os"
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
	// ここで様々なこと行う
	fmt.Println(client)
}
