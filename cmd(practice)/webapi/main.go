package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//"errors"
)

type photo struct {
	ID      string `json:"id"`
	User    string `json:"user"`
	Date    string `json:"date"`
	Likes   int    `json:"likes"`
	Comment int    `json:"comment"`
}

var photos = []photo{
	{ID: "https://media.wired.com/photos/598e35fb99d76447c4eb1f28/master/pass/phonepicutres-TA.jpg", User: "Ilyas", Date: "2024-02-11", Likes: 10, Comment: 1},
}

func getPhotos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, photos)
}

func createPhoto(c *gin.Context) {

}

func main() {
	router := gin.Default()
	router.GET("/photos", getPhotos)
	router.Run("localhost:8080")
}
