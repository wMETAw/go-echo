package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func MyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("MyMiddleware")
		defer fmt.Println("defer MyMiddleware")
		return next(c)
	}
}

func main() {
	fmt.Println("111111111")
	e := echo.New()
	fmt.Println("222222222")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	fmt.Println("333333333")

	// 自作ミドルウェアをグループ化
	m := e.Group("", MyMiddleware)
	fmt.Println("44444444")

	// グループに適応させるルーティング
	m.GET("/", func(c echo.Context) error {
		fmt.Println("55555555")
		return c.String(http.StatusOK, "Hello")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

// 111111111
// 222222222
// 333333333
// 44444444
//
// ⇨ http://localhost:1323/ にアクセス
//
// MyMiddleware
// 55555555
// defer MyMiddleware

// https://medium.com/veltra-engineering/echo-middleware-in-golang-90e1d301eb27
