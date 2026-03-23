package webserver

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/edgehook/ithings/common/config"
	"github.com/edgehook/ithings/common/dbm/model"
	v1 "github.com/edgehook/ithings/common/types/v1"
	"github.com/edgehook/ithings/common/utils"
	"github.com/edgehook/ithings/webserver/router"
	"github.com/jwzl/beehive/pkg/core"
	"k8s.io/klog"
)

const (
	WebServerName = "webserver"
)

type WebServer struct {
}

// Register this module.
func Register() {
	ws := &WebServer{}
	core.Register(ws)
}

// Name
func (ws *WebServer) Name() string {
	return WebServerName
}

// Group
func (ws *WebServer) Group() string {
	return WebServerName
}

// Enable indicates whether this module is enabled
func (ws *WebServer) Enable() bool {
	//The module is always enabled!
	return true
}

// Start this module.
func (ws *WebServer) Start() {
	initDb()
	var err error

	initRouter := router.InitRouter()
	cfg := config.GetWebServerConfig()
	klog.Infof("Start web server on %s ", cfg.BindAddress)
	s := &http.Server{
		Addr:           cfg.BindAddress,
		Handler:        initRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 创建优雅关闭的信号通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 在 goroutine 中启动服务器
	serverErr := make(chan error, 1)
	go func() {
		if cfg.SSL {
			s.TLSConfig = createServerTLSConfiguration()
			err = s.ListenAndServeTLS(cfg.SSLCert, cfg.SSLKey)
		} else {
			err = s.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			serverErr <- err
		} else {
			serverErr <- nil
		}
	}()

	klog.Info("Web server started successfully. Press Ctrl+C to shut down.")

	// 等待关闭信号
	select {
	case <-quit:
		klog.Info("Shutting down server...")

		// 创建关闭上下文，最多等待 30 秒
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// 优雅关闭服务器
		if err := s.Shutdown(ctx); err != nil {
			klog.Errorf("Server forced to shutdown: %v", err)
		}
		klog.Info("Server exited")

	case err := <-serverErr:
		if err != nil {
			klog.Errorf("Start web server with error: %v", err)
			return
		}
	}
}

// createServerTLSConfiguration creates a basic tls.Config to be used by servers with recommended TLS settings
func createServerTLSConfiguration() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}
}

func initDb() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				klog.Errorf("initDb failed with err: %v", string(debug.Stack()))
			}
		}()
		enPwd := utils.Md5V(v1.DefaultPassword)
		if !model.IsExistUserByName(v1.DefaultUsername) {
			klog.Infof("admin user exist")
			id := utils.NewUUID()
			if err := model.AddUser(&model.User{
				ID:       id,
				Name:     v1.DefaultUsername,
				Password: enPwd,
				Rule:     "",
			}); err != nil {
				klog.Errorf("Add user config error: %v", err.Error())
			}
		}
	}()

}
