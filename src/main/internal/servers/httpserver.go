package servers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-ginApp/src/main/internal/config"
	"go-ginApp/src/main/pkg/common/middleware/gin_logger"
	"go-ginApp/src/main/pkg/common/middleware/self_pprof"
	"go-ginApp/src/main/pkg/utils/array"
	"go-ginApp/src/main/pkg/utils/datetool"
	"go-ginApp/src/main/pkg/utils/prof"
	"log"
	"net/http"
	"time"
)

func HttpInstance() *HTTPServer {
	if instance == nil {
		instance = &HTTPServer{
			reqCnt:  make(map[string]*prof.CountQPS),
			errCnt:  make(map[string]*prof.CountBase),
			reqTime: make(map[string]*prof.AverageTimer),
		}
		//gin 路由初始化
		r := gin.New()
		gin_logger.Global(r) // 全局中间件校验
		gin.SetMode(config.AppConfig.Server.RunMod)
		// 注册pprof
		self_pprof.RegisterPprof(r)
		//注册restful路由
		RegisterRestfulRoute(r, instance)
		//注册graphql路由

		RegisterGraphqlRoute(r)
		//开启监听服务端口
		err := r.Run(":8089")
		if err != nil {
			return nil
		}

		log.Printf("[info] start http server listening %s", config.AppConfig.Server.Addr)
	}
	return instance
}

func (hs *HTTPServer) Start() error {
	if config.AppConfig.Server.Addr == "" {
		return fmt.Errorf("server addr is nil")
	}
	go func() {
		if config.AppConfig.Server.CertFile != "" && config.AppConfig.Server.KeyFile != "" {
			config.AppLog.Info("server.http : start https", "addr", config.AppConfig.Server.Addr)
			err := hs.server.ListenAndServeTLS(config.AppConfig.Server.CertFile, config.AppConfig.Server.KeyFile)
			if err != http.ErrServerClosed {
				config.AppLog.Fatal("server.http : start https fail", "addr", config.AppConfig.Server.Addr, "error", err)
			}
			return
		}
		config.AppLog.Info("server.http : start", "addr", config.AppConfig.Server.Addr)
		err := hs.server.ListenAndServe()
		if err != http.ErrServerClosed {
			config.AppLog.Fatal("server.http : start fail", "addr", config.AppConfig.Server.Addr, "error", err)
		}
	}()
	return nil
}

func (hs *HTTPServer) profilerMiddleware(c *gin.Context) {
	st := time.Now()
	path := c.Request.URL.Path

	// 执行真实路由
	c.Next()
	// 尝试创建计数器
	if _, ok := hs.reqCnt[path]; !ok {
		hs.reqCnt[path] = prof.NewCountQPS("req")
		hs.reqTime[path] = prof.NewAverageTimer("time", 60)
		hs.errCnt[path] = prof.NewCountBase("error")
	}
	// 统计计数
	hs.reqCnt[path].Incr()
	hs.reqTime[path].Set(datetool.Since(st))
	if !array.IntFind([]int{http.StatusOK, http.StatusNoContent}, c.Writer.Status()) {
		hs.errCnt[path].Incr()
	}
}
