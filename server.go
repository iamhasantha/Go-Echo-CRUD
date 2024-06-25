package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"sync"
)

type User struct {
	ID    string `json:"id" form:"id" query:"id"`
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

var (
	users = make(map[string]User)
	mu    sync.Mutex
)

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	mu.Lock()
	user, exists := users[id]
	mu.Unlock()
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func getUserByName(c echo.Context) error {
	name := c.Param("name")
	mu.Lock()
	//look for all users with the same name
	var user User
	exists := false
	for _, u := range users {
		if u.Name == name {
			user = u
			exists = true
			break
		}
	}
	mu.Unlock()
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}
	return c.JSON(http.StatusOK, user)

}

func createUser(c echo.Context) error {
	id := generateID()
	name := c.FormValue("name")
	email := c.FormValue("email")
	user := User{ID: id, Name: name, Email: email}
	mu.Lock()
	users[id] = user
	mu.Unlock()
	return c.JSON(http.StatusCreated, user)

}

func getAllUsers(c echo.Context) error {
	mu.Lock()
	defer mu.Unlock()
	return c.JSON(http.StatusOK, users)

}

func updateUser(c echo.Context) error {
	id := c.Param("id")
	name := c.FormValue("name")
	email := c.FormValue("email")
	mu.Lock()
	user, exists := users[id]
	if !exists {
		mu.Unlock()
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}
	user.Name = name
	user.Email = email
	users[id] = user
	mu.Unlock()
	return c.JSON(http.StatusOK, user)

}

func updateUserPatch(c echo.Context) error {
	id := c.Param("id")
	name := c.FormValue("name")
	email := c.FormValue("email")
	mu.Lock()
	user, exists := users[id]
	if !exists {
		mu.Unlock()
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}
	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}
	users[id] = user
	mu.Unlock()
	return c.JSON(http.StatusOK, user)


}

func deleteUser(c echo.Context) error {
	id := c.Param("id")
	mu.Lock()
	_, exists := users[id]
	mu.Unlock()
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}
	mu.Lock()
	delete(users, id)
	mu.Unlock()
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted"})


}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/users/:id", getUser)
	e.GET("/users/name/:name", getUserByName)
	e.POST("/users", createUser)
	e.GET("/users/all", getAllUsers)
	e.PUT("/users/update/:id", updateUser)
	e.PATCH("/users/updateOne/:id", updateUserPatch)
	e.DELETE("/users/delete/:id", deleteUser)
	e.Logger.Fatal(e.Start(":1323"))

}
