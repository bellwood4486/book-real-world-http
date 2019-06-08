package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// メソッドが属す構造体
type Calculator int

// RPC で外部から呼ばれるメソッド
func (c *Calculator) Multiply(args Args, result *int) error {
	log.Printf("Multiply called: %d, %d\n", args.A, args.B)
	*result = args.A * args.B
	return nil
}

// 外部から呼ばれる時の引数
type Args struct {
	A, B int
}

func main() {
	calculator := new(Calculator)
	server := rpc.NewServer()
	_ = server.Register(calculator)
	http.Handle(rpc.DefaultRPCPath, server)
	log.Println("start http listening :18888")
	listener, err := net.Listen("tcp", ":18888")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
