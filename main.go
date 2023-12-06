package main

import (
	"fmt"
	"net/http"

	//"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Note struct {
	User string `json:"user,omitempty"`
	Id   uint   `json:"id" gorm:"primaryKey"`
	Note string `json:"note"`
}

type NoteDel struct {
	User string `json:"sid"`
	Id   uint   `json:"id"`
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/signup", signingup)
	r.POST("/login", loggingin)
	// Middleware to check for a valid session token
	r.Use(AuthMiddleware())
	r.GET("/notes", getNotes)
	r.POST("/notes", createNote)
	r.DELETE("/notes", deleteNote)

	InitDB()
	r.Run() // listen and serve on 0.0.0.0:8080

	//http.ListenAndServe(":8080", nil)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")

		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func signingup(c *gin.Context) {

	var user1 User
	if err := c.BindJSON(&user1); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Request Format Invalid"})
		return
	}
	fmt.Println(user1)
	GetDB().Create(&user1)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success"})
}

func loggingin(c *gin.Context) {
	var login Login
	if err := c.BindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	var count int64
	GetDB().Debug().Model(&User{}).Where("email = ? and password = ?", login.Email, login.Password).Count(&count)
	if count == 1 {
		// Create or retrieve the session
		session := sessions.Default(c)
		session.Set("user", login.Email)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			fmt.Println(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

func getNotes(c *gin.Context) {
	userId := c.Query("sid")
	var notes []Note

	GetDB().Debug().Select("id, note").Where("user = ?", userId).Find(&notes)

	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func createNote(c *gin.Context) {
	var input Note
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	GetDB().Create(&input)
	c.JSON(http.StatusOK, input.Id)
}

func deleteNote(c *gin.Context) {
	var input NoteDel
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var note Note
	if err := GetDB().Select("id").Where("id = ? and user = ?", input.Id, input.User).First(&note).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	GetDB().Delete(&note)
	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
