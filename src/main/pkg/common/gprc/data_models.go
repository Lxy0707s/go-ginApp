syntax = "proto3";
package gprc

option go_package = "./protos;protos";

// 定义发送请求信息
message UserRequest{
// 定义发送的参数
// 参数类型 参数名 标识号(不可重复) 1  表示确定顺序
string  user_id = 1;
}

// 定义响应信息返回信息
message UserResponse{
string user_id = 1;
int32 score = 2;
int32 age = 3;
string user_name = 4;
}

service UserService{
rpc GetUserInfo (UserRequest) returns (UserResponse){};
}
