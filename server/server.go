package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"stlviewer/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Server() {
	app := fiber.New()

	if utils.Env.DevMode {
		//logger
		file, err := os.OpenFile(fmt.Sprintf("%v/server.log", utils.DirPath), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer file.Close()
		app.Use(logger.New(logger.Config{
			Next:         nil,
			Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
			TimeFormat:   "15:04:05",
			TimeZone:     "Local",
			TimeInterval: 500 * time.Millisecond,
			Output:       file,
		}))
	}

	// monitor
	app.Get("/dashboard", monitor.New())

	// recover
	app.Use(recover.New())

	// pprof - localhost:port/debug/pprof/
	app.Use(pprof.New())

	app.Use(filesystem.New(filesystem.Config{
		Root:         http.Dir(fmt.Sprintf("%v/dist", utils.DirPath)),
		Browse:       true,
		Index:        "index.html",
		NotFoundFile: "404.html",
		MaxAge:       3600,
	}))

	app.Listen(fmt.Sprintf("%v:%v", utils.Env.IP, utils.Env.Port))
}
