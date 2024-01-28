package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/logger/cclog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var DefaultAddr = ":28883"

func StartServer(router *gin.Engine) {
	server := &http.Server{
		Addr:    DefaultAddr,
		Handler: router,
	}

	quit := make(chan os.Signal)
	// kill -2
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		cclog.Info("receive interrupt signal")
		if err := server.Close(); err != nil {
			cclog.Error("Server Close:", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			cclog.Error("Server closed under request")
		} else {
			cclog.Error("Server closed unexpect")
		}
	}

	cclog.Info("Server exiting")
}

// 用ctx初始化资源，mysql，redis等,内部服务也可以通过mainCtx初始化
var MainCtx, mainStop = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

// StartWithContextNotify
//
//	@param router
func StartWithContextNotify(router *gin.Engine) {
	// Create context that listens for the interrupt signal from the OS.
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it

	defer mainStop()

	srv := &http.Server{
		Addr:    DefaultAddr,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cclog.Error("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-MainCtx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	mainStop()
	cclog.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 服务关闭，不再接受请求
	if err := srv.Shutdown(ctx); err != nil {
		cclog.Error("Server forced to shutdown: ", err)
	}
	// 业务关闭
	HandleClose(ctx)

	cclog.Info("Server exiting")
}

func HandleClose(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			cclog.Info("timeout done")
		}
		break
	}
}

func StartWithSignalNotify(router *gin.Engine) {
	srv := &http.Server{
		Addr:    DefaultAddr,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cclog.Error("listen: %s\n", err)

		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2,ctrl+C is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	cclog.Info("Shutting down server...", <-quit)

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		cclog.Error("Server forced to shutdown: ", err)
	}

	cclog.Info("Server exiting")
}
