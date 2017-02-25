package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/micrypt/go-plivo/plivo"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// Port at which the server starts listening
	Port = ":8080"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	// Simple group: v1
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, X-HTTP-Method-Override,Authorization, Content-Type, Accept")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		if c.Request.Method == "OPTIONS" {
			c.JSON(200, gin.H{"All": "Good"})
		} else {
			c.Next()
		}
	})
	// r.POST("/", upload)
	r.GET("/", Get)
	r.GET("/record", Record)
	r.GET("/plivo/callback", Callback)
	r.POST("/", Create)

	// Start listening
	port := Port
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}
	r.Run(port)
}

// func upload(c *gin.Context) {
// 	// fmt.Println(c.Request.MultipartForm)
// 	file, header, err := c.Request.FormFile("file")
// 	// fmt.Println("file", file)
// 	// fmt.Println("header", header)
// 	// fmt.Println("err", err)
// 	filename := header.Filename
// 	fmt.Println("yoooo", header.Filename)
// 	out, err := os.Create("./" + filename + ".png")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer out.Close()
// 	_, err = io.Copy(out, file)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func Get(c *gin.Context) {
	file, err := os.Open("./plivo.xml")
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	// fmt.Println("data", data)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	c.Data(200, "text/xml", data)
}

func Create(c *gin.Context) {
	client := plivo.NewClient(nil, os.Getenv("PLIVO_AUTH_ID"), os.Getenv("PLIVO_AUTH_TOKEN"))
	req := &plivo.CallMakeParams{From: os.Getenv("PHONE_FROM"), To: os.Getenv("PHONE_TO"), AnswerURL: os.Getenv("ANSWER_URL"), AnswerMethod: "GET"}
	_, _ = client.Call.Make(req)
}

func Record(c *gin.Context) {
	client := plivo.NewClient(nil, os.Getenv("PLIVO_AUTH_ID"), os.Getenv("PLIVO_AUTH_TOKEN"))
	req := &plivo.CallRecordParams{TimeLimit: 60, FileFormat: "mp3", CallbackURL: os.Getenv("CALLBACK_URL"), CallbackMethod: "GET"}
	_, _ = client.Call.Record(c.Query("CallUUID"), req)
}

func Callback(c *gin.Context) {
	str := strings.Split(strings.Split(c.Query("response"), "record_url%22%3A%22")[1], "%22%2C%22recording_duration")[0]
	str = strings.Replace(str, "%3A", ":", -1)
	str = strings.Replace(str, "%5C%2F", "/", -1)
	fmt.Println("record_url", str)
}
