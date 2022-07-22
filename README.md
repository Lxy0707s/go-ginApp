# go-ginApp

golang 整合Gin框架、graphql

本项目支持：

 1、Gin restful路由、graphql路由

 2、http/tcp请求

 3、基本探测组件

 4、基本鉴权组件等

 目录结构

 ├─.idea
└─src
    └─main
        ├─app
        ├─internal
        │  ├─config
        │  ├─dao
        │  │  └─db_models
        │  ├─graphql
        │  │  ├─gqlgen
        │  │  ├─schema
        │  │  └─sdl
        │  ├─middleware
        │  └─servers
        ├─moudule
        │  ├─book
        │  │  ├─book_restful
        │  │  │  └─book
        │  │  └─book_service
        │  │      └─book
        │  ├─graphql_demo
        │  │  └─graphql_query
        │  ├─test_task
        │  │  ├─task_restful
        │  │  │  └─task
        │  │  └─task_service
        │  │      └─task
        │  └─user_manager
        │      ├─user_restful
        │      │  └─user
        │      └─user_service
        │          └─user
        └─pkg
            ├─common
            │  ├─common_type
            │  ├─detect
            │  │  ├─curl
            │  │  ├─dns
            │  │  ├─grpc_check
            │  │  ├─http_https
            │  │  ├─myicmp
            │  │  ├─tcp
            │  │  └─udp
            │  ├─gprc
            │  ├─graphql
            │  ├─httpserver
            │  │  ├─httpClient
            │  │  └─httpServer
            │  ├─middleware
            │  │  ├─gin_logger
            │  │  ├─self_pprof
            │  │  └─sys_jwt
            │  └─trace
            │      ├─models
            │      └─trace_service
            └─utils
                ├─array
                ├─base_struct
                ├─cmd
                ├─configtool
                ├─datetool
                ├─dbtool
                ├─gziptool
                ├─httptool
                ├─logtool
                ├─myfile
                ├─prof
                ├─protocoltool
                └─queue
