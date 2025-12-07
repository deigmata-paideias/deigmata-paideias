package main

import (
	"encoding/json"
	"os"

	"github.com/caddyserver/caddy/v2"
	// 加载 caddy 模块必需
	_ "github.com/caddyserver/caddy/v2/modules/standard"
)

// go run main.go
func main() {

	cfgFile, err := os.ReadFile("caddyfile.json")
	if err != nil {
		panic(err)
	}

	var cfg caddy.Config
	if err := json.Unmarshal(cfgFile, &cfg); err != nil {
		panic(err)
	}

	if err := caddy.Run(&cfg); err != nil {
		panic(err)
	}

	select {}
}
