package cli

import (
	"flag"
	"fmt"

	"fiber_curd_api/app"
	"fiber_curd_api/tools/db"
)

type CliOptions struct {
	Migrate bool
	Seed    bool
	Version bool
}

var cliOpts CliOptions

func Run() bool {
	// 1. 定义参数
	// flag 包原生支持 -h，会自动输出用法
	flag.BoolVar(&cliOpts.Migrate, "migrate", false, "Run database migrations")
	flag.BoolVar(&cliOpts.Seed, "seed", false, "Seed the database with initial data")
	flag.BoolVar(&cliOpts.Version, "version", false, "Print version information")

	// 解析参数
	flag.Parse()

	if cliOpts.Version {
		// 2. 处理 --version
		fmt.Println("Version: ", app.VERSION)
		return true
	}

	if cliOpts.Migrate {
		// 3. 处理 --migrate
		db.Migrate()
		return true
	}

	if cliOpts.Seed {
		// 4. 处理 --seed
		db.Seed()
		return true
	}

	return false
}
