package main

import (
	"dnsutils"
	"fmt"
	"os"
	"restapi/controllers"
	"restapi/middlewares"
	"runtime"

	"github.com/go-redis/redis"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/basicauth"
	"github.com/kataras/iris/middleware/logger"
)

var basicAuth = basicauth.New(basicauth.Config{
	Users: map[string]string{
		"admin": "qwe123!Q",
		"test":  "123456Qw!",
	},
})

func main() {
	NCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(NCPU)
	// initialize redis connect
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "Cyber123456", // no password set
		DB:       0,             // use default DB
	})

	if _, err := client.Ping().Result(); err != nil {
		fmt.Printf("Redis Connect Failed!")
		os.Exit(-1)
	} else {
		fmt.Println("Redis Connect Succeed!")
	}

	// initialize tencent lib
	if err := dnsutils.UrlLibInit(263785, 150000000, "/home/privateCloud/conf/licence.conf"); err != nil {
		fmt.Printf("UrlLibInit Failed\tErrcode:%#x", err)
		os.Exit(-2)
	} else {
		fmt.Println("UrlLibInit Succeed!")
	}

	// initialize logger
	c := &logger.Config{
		Status:            true,
		IP:                true,
		Method:            true,
		Path:              true,
		Columns:           false,
		MessageContextKey: "logger_message",
	}
	cp := &middlewares.ConfigPlus{
		LogPath: "./log",
	}
	r, close := middlewares.NewRequestLogger(c, cp)
	defer close()

	// instantiation
	app := iris.New()

	// use logger
	app.Use(r)
	app.OnAnyErrorCode(r, func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]interface{}{
			"success": false,
		})
	})

	// user auth
	app.Use(basicAuth)

	// bind controller to route
	app.Controller("/icp", new(controllers.IcpController), client)
	app.Controller("/urlsafe", new(controllers.UrlsafeController))

	// run
	app.Run(iris.Addr(":3000"), iris.WithoutServerError(iris.ErrServerClosed))
}
