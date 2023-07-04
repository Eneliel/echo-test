package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	e := echo.New()
	e.POST("/save", save)
	e.GET("/save1/:name", save1)
	e.GET("/save2", save2)
	e.POST("user", AddUser)
	e.POST("user2", AddSecondUser)
	e.POST("user3", AddThirdUser)
	e.Logger.Fatal(e.Start(":1323"))
}

func save(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "Its FormValue: name:"+name+", email:"+email)
}

func save1(c echo.Context) error {
	name := c.Param("name")
	return c.String(http.StatusOK, "Its Param: name:"+name)
}

func save2(c echo.Context) error {
	name := c.QueryParam("name")
	datatype := c.QueryParam("type")
	if datatype == "string" {
		return c.String(http.StatusOK, "Its QueryParam: name:"+name)
	}
	if datatype == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": name,
		})
	}
	return c.JSON(http.StatusBadRequest, "Error: Invalid type")
}

func AddUser(c echo.Context) error {
	u := User{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Error, Failed reading the request body for User: %s", err)
	}
	err = json.Unmarshal(b, &u)
	if err != nil {
		log.Printf("Error, Failed unmarshaling in User:%s", err)
	}
	log.Printf("user:%v", u)
	return c.String(http.StatusOK, fmt.Sprintf("Name:%v Age:%v", u.Name, u.Age))
}

func AddSecondUser(c echo.Context) error {
	u := User{}
	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&u)
	if err != nil {
		log.Printf("Failed decode process %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("Done!!")
	return c.String(http.StatusOK, fmt.Sprintf("Name:%v Age:%v", u.Name, u.Age))
}

func AddThirdUser(c echo.Context) error {
	u := User{}
	err := c.Bind(&u)
	if err != nil {
		log.Printf("Failed Bind User %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("Done!")
	return c.String(http.StatusOK, fmt.Sprintf("Name:%v Age:%v", u.Name, u.Age))
}
