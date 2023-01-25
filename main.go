package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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
	Pos       int          `json:"pos"`
	Data      BookCheckout `json:"data"`
	TimeStamp string       `json:"time_stamp"`
	Hash      string       `json:"hash"`
	PrevHash  string       `json:"prev_hash"`
}

type Blockchain struct {
	Blocks []*Block
}

var blockchain *Blockchain

func (bc *Block) generatehash() {
	bytes, _ := json.Marshal(bc.Data)
	data := string(rune(bc.Pos)) + bc.TimeStamp + bc.PrevHash + string(bytes)
	hash := sha256.New()
	hash.Write([]byte(data))
	bc.Hash = hex.EncodeToString(hash.Sum(nil))
}

func CreateBlock(prevblock *Block, checkoutitem BookCheckout) *Block {
	block := &Block{}
	block.Pos = prevblock.Pos + 1
	block.TimeStamp = time.Now().GoString()
	block.Data = checkoutitem
	block.PrevHash = prevblock.Hash
	block.generatehash()
	return block
}

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckout{
		BookId:       "nil",
		User:         "nil",
		CheckoutDate: "nil",
		IsGenesis:    true,
	})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func validblock(prevblock, block *Block) bool {
	if prevblock.Hash != block.PrevHash {
		return false
	}
	if !block.validatehash(block.Hash) {
		return false
	}
	if prevblock.Pos+1 != block.Pos {
		return false
	}
	return true
}

func (bc *Blockchain) AddBlock(data BookCheckout) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	block := CreateBlock(prevBlock, data)
	if validblock(prevBlock, block) {
		bc.Blocks = append(bc.Blocks, block)
	}
}

func (b *Block) validatehash(hash string) bool {
	b.generatehash()
	return b.Hash == hash
}

func main() {
	blockchain = NewBlockchain()
	r := gin.Default()
	r.GET("/", GetBlockchain)
	r.POST("/", CreateBlocks)
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
	io.WriteString(h, body.PublishDate+body.Isbn)
	body.Id = fmt.Sprintf("%x", h.Sum(nil))
	mapd["error"] = false
	mapd["data"] = body

	c.JSON(http.StatusOK, mapd)
}

func CreateBlocks(c *gin.Context) {

	var checkoutitem BookCheckout
	err := c.Bind(&checkoutitem)
	if err != nil {
		log.Println(err)
	}
	blockchain.AddBlock(checkoutitem)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  checkoutitem,
	})
}

func GetBlockchain(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  blockchain.Blocks,
	})
}
