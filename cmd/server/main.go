package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yinyin/myblog/internal/bootstrap"
)

// version 由 -ldflags 注入,默认 dev
var version = "dev"

func main() {
	cfgPath := flag.String("config", "", "path to config file (yaml). empty = auto-discover configs/config.yaml")
	showVer := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVer {
		fmt.Println("myblog", version)
		return
	}

	if err := bootstrap.Run(*cfgPath, version); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}
