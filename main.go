package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/server"
	"net"
	"sync"
	"sync/atomic"
	gopush "test/kitex-issue/kitex_gen"
	"test/kitex-issue/kitex_gen/gopushservice"
	"time"
)

type (
	app struct {
		count    int64 //连接数统计
		streamId int64 //grpc stream id
		addr     string
	}

	streamer struct {
		stream gopush.GoPushService_OnPushMessageServer
		id     int64
	}

	Req struct {
		data []byte
	}
)

func main() {
	impl := &app{
		count:    0,
		streamId: 0,
	}

	addr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:8888")
	srv := gopushservice.NewServer(impl, server.WithServiceAddr(addr))

	go func() {
		err := srv.Run()
		fmt.Println(err)
	}()

	//服务端的方法是否感知到这个事件
	for i := 0; i < 5; i++ {

		cli, _ := gopushservice.NewClient("gopushservice", client.WithHostPorts("0.0.0.0:8888"))

		go func() {
			_, err := cli.OnPushMessage(context.Background())
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	time.Sleep(time.Second)
	fmt.Println("当前连接的stream: ", impl.count)
	if impl.count != 5 {
		fmt.Println("连接的客户端数量错误 ...")
	}

	err := srv.Stop()
	if err != nil {
		fmt.Println(err)
	}

	//让服务端进行感知
	time.Sleep(time.Second * 5)
	fmt.Println("当前连接的stream: ", impl.count)
	if impl.count != 0 {
		fmt.Println("关闭的客户端数量错误 ...")
	}
}

// ignored
func (a *app) Push(context.Context, *gopush.GoPushRequest) (*gopush.GoPushResponse, error) {
	return &gopush.GoPushResponse{}, nil
}

func (a *app) OnPushMessage(stream gopush.GoPushService_OnPushMessageServer) error {

	atomic.AddInt64(&a.count, 1)
	streamId := atomic.AddInt64(&a.streamId, 1)

	var wait sync.WaitGroup
	wait.Add(1)

	a.recv(&wait, stream, streamId)

	wait.Wait()

	//某个连接关闭了
	atomic.AddInt64(&a.count, -1)
	return nil
}

func (a *app) recv(wait *sync.WaitGroup, s gopush.GoPushService_OnPushMessageServer, streamId int64) {

	defer wait.Done()

	var (
		msg = &Req{data: make([]byte, 1024)}
		_   = &streamer{
			stream: s,
			id:     streamId,
		}
	)

	for true {
		err := s.RecvMsg(&msg)
		if err != nil {
			println("receive err: %v", err)
			return
		}
		// mock biz code
		time.Sleep(100 * time.Millisecond)
	}

}
