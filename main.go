package main

import (
	"crypto/md5"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookCheckout struct {
	BookId       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

type Book struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishDate string `json:"publish_date"`
	Isbn        string `json:"isbn"`
}

type Block struct {
	Pos       string       `json:"pos"`
	Data      BookCheckout `json:"data"`
	TimeStamp string       `json:"time_stamp"`
	Hash      string       `json:"hash"`
	PrevHash  string       `json:"prev_hash"`
}

type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AddBlock() {
	prebBlock := bc.Blocks[len(bc.Blocks)-1]
	block := CreateBlock()
	
}

func main() {
	r := gin.Default()
	r.GET("/", GetBlockchain)
	r.POST("/", CreateBlock)
	r.POST("/new", NewBook)
	r.Run()
}

func NewBook(c *gin.Context) {
	mapd := make(map[string]interface{})
	var body Book
	err := c.Bind(&body)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error":   true,
			"message": "not created",
		})
		return
	}
	h := md5.New()
	io.WriteString(h, body.Isbn+body.Isbn)
	body.Id = string(h.Sum(nil))
	log.Println(body.Id)

	data := Book{
		Id:          body.Id,
		Title:       body.Title,
		Author:      body.Author,
		PublishDate: body.PublishDate,
		Isbn:        body.Isbn,
	}
	log.Println(data)
	mapd["error"] = false
	mapd["data"] = data

	c.JSON(http.StatusOK, mapd)
}

func CreateBlock(c *gin.Context) {

}

func GetBlockchain(c *gin.Context) {

}
