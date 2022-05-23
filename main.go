package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Todo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type User struct {
	Id   int    `orm:"auto"`
	Name string `orm:"column(name)"`
}

func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Server working fine, Active !")
}

func usersQuery(c echo.Context) error {
	name := c.QueryParam("name")
	pwd := c.QueryParam("pwd")

	return c.String(http.StatusOK, fmt.Sprintf("Query Check for user name of %s and password is %s", name, pwd))
}

func usersParam(c echo.Context) error {
	id := c.Param("id")

	response := map[string]string{
		"status": "Success",
		"id":     id,
	}

	return c.JSON(http.StatusOK, response)
}

func addTodo(c echo.Context) error {
	todo := Todo{}
	resp := &Response{}

	err := c.Bind(&todo)
	if err != nil {
		log.Printf("Failed the request")
		return c.String(http.StatusInternalServerError, "")
	}

	resp.Status = "Success"
	resp.Data = todo

	return c.JSON(http.StatusOK, resp)
}

func admin(c echo.Context) error {
	return c.String(http.StatusOK, "Admin portal !")
}

func getTodos(c echo.Context) error {

	todo1 := Todo{
		Name:        "Learn DS",
		Description: "Focus more on this part !",
		ID:          1,
	}

	todo2 := Todo{
		Name:        "Revise Node Mongo Client",
		Description: "Learn on how to write connections and queries !",
		ID:          2,
	}

	todos := []Todo{todo1, todo2}
	resp := &Response{}

	resp.Data = todos
	resp.Status = "Success"

	return c.JSON(http.StatusOK, resp)
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(localhost:3306)/goCheck")
	orm.RegisterModel(new(User))
}

func main() {

	fmt.Println("Welocme to the server !")

	// o := orm.NewOrm()

	// user := User{Name: "Amma", Id: 4}

	// // insert
	// id, err := o.Insert(&user)
	// fmt.Printf("ID: %d, ERR: %v\n", id, err)

	e := echo.New()

	g := e.Group("/admin")

	// First middleware to logs ...
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}]  Status: ${status}  Method: ${method} Api: ${host}/${path}` + "\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://labstack.com", "https://labstack.net", "http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "shashank" && password == "Greesh123" {
			return true, nil
		}
		return false, nil
	}))

	g.GET("/", admin)

	// Routes ...
	e.GET("/", healthCheck)
	e.GET("/users", usersQuery)
	e.GET("/users/:id", usersParam)
	e.GET("/todos", getTodos)

	e.POST("/todo", addTodo)

	// Logger ...
	e.Logger.Fatal(e.Start(":1323"))
}
