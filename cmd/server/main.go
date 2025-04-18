package main

import (
	"article-service/cmd"
	"article-service/config"
)

func init() {
	config.InitConfig()
}

func main() {
	cmd.Execute()
}
