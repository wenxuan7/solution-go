package main

import (
	settingsproto "github.com/wenxuan7/oms-grpc/settings"
	"github.com/wenxuan7/solution/external"
	"google.golang.org/grpc"
	"log"
	"net"
)

func setup() {
	external.Mysql()
	external.Redis()
}

func main() {
	setup()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// 注册服务
	settingsproto.RegisterReaderServer(s, NewSettingsReaderServer())
	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
