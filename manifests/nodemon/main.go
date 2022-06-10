package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"log"
	"nodemon/pkg/k8s"
	"nodemon/pkg/server"
)

type Env struct {
	NodeLabel string `default:"app=renchain-node" split_words:"true"`
	NodePort  string `default:"8545" split_words:"true"`
	MaxDiff   uint64 `default:"5" split_words:"true"`
	Namespace string `required:"true" split_words:"true"`
}

func main() {
	var env Env
	envconfig.MustProcess("", &env)

	client, err := k8s.GetK8sClient()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	srv := server.Handler{
		RenNodeLabel: env.NodeLabel,
		Ns:           env.Namespace,
		Client:       client,
		Port:         env.NodePort,
		MaxBlockDiff: env.MaxDiff,
	}

	router.GET("/addresses", srv.GetAddressHandler)
	router.GET("/sync-status", srv.CheckSync)
	router.POST("/generic", srv.GenericRpc)

	router.Run()
}
