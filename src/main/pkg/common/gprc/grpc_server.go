package gprc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

const DefaultUserId = "11111"

var userMap = map[string]User{
	DefaultUserId: {
		UserId:   DefaultUserId,
		UserName: "user_name_1",
		Score:    101,
		Age:      10,
	},
}

type UserInfoService struct{}

type User struct {
	UserId   string
	UserName string
	Score    int32
	Age      int32
}

func (c *UserInfoService) GetUserInfo(ctx context.Context, in *UserRequest) (*pb.UserResponse, error) {
	response := new(pb.UserResponse)
	userId := in.GetUserId()
	fmt.Println("user_id: ", userId)
	if userId != "" {
		user, ok := userMap[userId]
		if ok {
			response.UserId = user.UserId
			response.UserName = user.UserName
			response.Age = user.Age
			response.Score = user.Score
		}
	}
	return response, nil
}

func main() {
	// 监听本地端口// Address 监听地址// Network 网络通信协议
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	pb.RegisterUserServiceServer(server, &UserInfoService{})

	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("server err: %v", err)
	}
}
