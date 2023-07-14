package main

import (
	"database/sql"
	"log"
	"net/http"

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

	g := e.Group("/main")

	u_name, pass := GetUser()
	g.Use(middleware.BasicAuth(func(username, password string, ctx echo.Context) (bool, error) {

		if username == u_name && password == pass {
			return true, nil
		}
		return false, nil
	}))

	g.GET("/newget", NewGet)

	e.Logger.Fatal(e.Start(":8000"))
}

func NewGet(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, its main group!!")
}

func GetUser() (username, pass string) {
	u := User{}
	db, err := sql.Open("mysql", "root:@/Test_Auth")
	if err != nil {
		log.Printf("Ошибка при входе в DB:%s", err)
	}
	defer db.Close()
	sel := "SELECT `username`, `password` FROM `user`"
	sel_db, err := db.Query(sel)
	if err != nil {
		log.Printf("Ошибка с запросов SELECT:%s", err)
	}
	for sel_db.Next() {
		err = sel_db.Scan(&u.Name, &u.Password)
		if err != nil {
			log.Printf("Ошибка при записи с DB: %s", err)
		}
	}

	return u.Name, u.Password
}
