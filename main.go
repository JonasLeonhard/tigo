package main

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/jonasleonhard/go-htmx-time/templates/components"
	"github.com/jonasleonhard/go-htmx-time/templates/pages"
)

func main() {
	e := echo.New()

	// Routes
	e.GET("/", index)
	e.GET("/name/:name/edit", editName)

	e.Static("/static", "static")
	e.Logger.Fatal(e.Start(":3000"))
}

func asHtml(c echo.Context, component templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func index(c echo.Context) error {
	component := pages.IndexPage()
	return asHtml(c, component)
}

func editName(c echo.Context) error {
	component := components.EditName("editedName")
	return asHtml(c, component)
}
