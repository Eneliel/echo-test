package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"pass"`
}

func main() {

	e := echo.New()
	e.Use(ServerHeader)

	g := e.Group("/main")

	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_unix} ${host} ${method} ${path} ${route} ${protocol}",
	}))

	g.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {
		if GetUser(username, password) {
			cookie := &http.Cookie{} // cookie := new(http.Cookie)

			cookie.Name = "sessionID"
			cookie.Value = "some_string_value"
			cookie.Expires = time.Now().Add(48 * time.Hour)

			ctx.SetCookie(cookie)
			return true, nil
		}
		return false, nil
	}))

	g.GET("/maincookie", MainCookie)
	g.GET("/newget", NewGet)
	g.Use(checkCookie)
	e.Logger.Fatal(e.Start(":8000"))
}

// Cookie

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "any_server/v1.1")
		c.Response().Header().Set("AnyKeyHeader", "AnyValueHeader")
		return next(c)
	}
}

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionID")
		if err != nil {
			log.Println(err)
			return err
		}
		if cookie.Value == "some_string_value" {
			return next(c)
		}
		return c.String(http.StatusOK, "you dont have cookie")
	}
}

func NewGet(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, its main group!!")
}

func MainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, its your cookie!!")
}

func GetUser(u_name, pass string) bool {
	u := User{}
	db, err := sql.Open("mysql", "root:@/test_auth")
	if err != nil {
		log.Printf("Проблема со входом в БД:%s", err)
	}
	defer db.Close()
	sel := "SELECT `username`, `password` FROM `user`"
	sel_db, err := db.Query(sel)
	if err != nil {
		log.Printf("Ошибка с запросом SELECT:%s", err)
	}
	for sel_db.Next() {
		err = sel_db.Scan(&u.Name, &u.Password)
		if err != nil {
			log.Printf("Ошибка при чтении с БД:%s", err)
		}
		if u.Name == u_name && u.Password == pass {
			return true
		}
	}
	return false
}
