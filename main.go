// main.go
package main

import (
	"sync-tool/cmd"
	"sync-tool/internal/logger"
)

func main() {
	logger.InitLogger()
	cmd.Execute()
}

