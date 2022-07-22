package gprc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpcserver/protos"
	"log"
)

const (
	Address = ":8000"
)

func main() {
	// 连接服务器
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc dail err: %v", err)
		return
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	// 创建发送结构体
	req := pb.UserRequest{
		UserId: "11111",
	}
	// 调用我们的服务(GetUserInfo方法)
	// 同时传入了一个 context.Context ，在有需要时可以让我们改变GRPC的行为，比如超时/取消一个正在运行的RPC
	res, err := client.GetUserInfo(context.Background(), &req)
	if err != nil {
		log.Fatalf("call getUserInfo err: %v", err)
		return
	}

	fmt.Println(res)
}
