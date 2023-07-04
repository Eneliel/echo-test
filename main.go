package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.POST("/save", save)
	e.GET("/save1/:name", save1)
	e.GET("/save2", save2)
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
