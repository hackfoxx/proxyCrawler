package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"proxyCrawler/config"
	"proxyCrawler/web/handler"
	"proxyCrawler/web/middleware"
	"time"
)

type Server struct {
	stopChan chan bool
}

func NewGinServer() *Server {
	return &Server{
		stopChan: make(chan bool),
	}
}

func (f *Server) Start() {
	go func() {
		for {
			select {
			case <-f.stopChan:
				fmt.Println("WebServer stopped")
				return
			default:
				{
					gin.DisableConsoleColor()
					file, _ := os.Create("log/access.log")
					gin.DefaultWriter = file
					r := gin.New()
					r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
						//定制日志格式
						return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
							param.ClientIP,
							param.TimeStamp.Format(time.RFC1123),
							param.Method,
							param.Path,
							param.Request.Proto,
							param.StatusCode,
							param.Latency,
							param.Request.UserAgent(),
							param.ErrorMessage,
						)
					}))
					r.Use(middleware.AuthMiddleware())
					//r.Use(middleware.Logger())
					err := r.SetTrustedProxies([]string{"127.0.0.1"})
					if err != nil {
						panic(err)
					}
					r = handler.SetHandler(r)
					err = r.Run(config.Cfg.Web.Addr)
					if err != nil {
						return
					}
				}
			}
		}
	}()
}
func (f *Server) Stop() {
	f.stopChan <- true
}
