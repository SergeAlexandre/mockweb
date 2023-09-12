package main

import (
	"github.com/pior/runnable"
	"mockweb/internal/config"
	"mockweb/internal/handlers"
	"mockweb/pkg/skserver"
)

func main() {
	config.Setup()
	server := skserver.New("Main", &config.Conf.Server, config.Log)

	hdl := &handlers.InfoHandler{}

	server.AddHandler("/info", "GET", hdl)

	runnableMgr := runnable.NewManager()
	runnableMgr.Add(server)
	runnable.Run(runnableMgr.Build())
}
