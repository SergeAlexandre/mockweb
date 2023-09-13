package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/pior/runnable"
	"mockweb/internal/config"
	"mockweb/internal/handlers"
	"mockweb/pkg/skserver"
)

func main() {
	config.Setup()
	fmt.Printf("mockweb version:%s   buildTS:%s\n", config.Version, config.BuildTs)
	sessionManager := scs.New()
	sessionManager.Cookie.Name = "mock_session"
	sessionManager.IdleTimeout = config.Conf.IdleTimeout
	sessionManager.Lifetime = config.Conf.SessionLifetime
	server := skserver.New("Main", &config.Conf.Server, config.Log)

	hdl := &handlers.InfoHandler{
		Log:            config.Log,
		SessionManager: sessionManager,
	}

	server.AddHandler("/info", "GET", sessionManager.LoadAndSave(hdl))

	runnableMgr := runnable.NewManager()
	runnableMgr.Add(server)
	runnable.Run(runnableMgr.Build())
}
