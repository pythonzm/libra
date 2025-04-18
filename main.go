package main

import (
	"log"
	"os"

	"libra/internal/cli"
)

func main() {
	// 设置日志格式，方便调试
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 创建并运行 CLI 应用
	app := cli.SetupApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Error running CLI app: %v", err)
	}
}
