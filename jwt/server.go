package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"time"
)

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "admin" && password == "admin" {
		// create TOKEN
		token := jwt.New(jwt.SigningMethodHS256)

		// set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = username
		claims["admnin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{"token": t})
	}
	return echo.ErrUnauthorized
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func unaccessible(c echo.Context) error {
	return c.String(http.StatusOK, "Unaccessible")
}

func restricted(c echo.Context) error {

	// TOKENを復号し値を取り出す
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", login)

	// Resticred group
	ex := e.Group("")
	ex.Use(middleware.JWT([]byte("secret")))

	// Authenticated route
	ex.GET("/restricted", restricted)

	// Unauthenticated route
	e.GET("/accessible", accessible)

	// Authenticated route
	ex.GET("/unaccessible", unaccessible)

	e.Logger.Fatal(e.Start(":1323"))
}

// # login
// curl -X POST -d 'username=admin' -d 'password=admin' localhost:1323/login
//
// curl localhost:1323/restricted
// => {"message":"Missing or malformed jwt"}
//
// curl localhost:1323/restricted -H "Authorization: Bearer HOGEFUGATOKEN"
// => Welcome amdmin!
//
// curl localhost:1323/accessible
// => Accessible
//
// curl localhost:1323/unaccessible
// => {"message":"Missing or malformed jwt"}
//
// curl localhost:1323/unaccessible -H "Authorization: Bearer HOGEFUGATOKEN"
// => Uncessible
