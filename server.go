package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/micrypt/go-plivo/plivo"
	"io/ioutil"
	"os"
)

type Response struct {
	RecordingStartMs    string `json:"recording_start_ms"`
	RecordingEndMs      string `json:"recording_end_ms"`
	CallUUID            string `json:"call_uuid"`
	APIID               string `json:"api_id"`
	RecordURL           string `json:"record_url"`
	RecordingDurationMs string `json:"recording_duration_ms"`
	RecordingID         string `json:"recording_id"`
	Message             string `json:"message"`
	RecordingDuration   string `json:"recording_duration"`
}

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

func Get(c *gin.Context) {
	file, err := os.Open("./plivo.xml")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	c.Data(200, "text/xml", data)
}

func Create(c *gin.Context) {
	client := plivo.NewClient(nil, os.Getenv("PLIVO_AUTH_ID"), os.Getenv("PLIVO_AUTH_TOKEN"))
	req := &plivo.CallMakeParams{From: os.Getenv("PHONE_FROM"), To: os.Getenv("PHONE_TO"), AnswerURL: os.Getenv("ANSWER_URL"), AnswerMethod: "GET"}
	_, err := client.Call.Make(req)
	if err != nil {
		panic(err)
	}
}

func Record(c *gin.Context) {
	client := plivo.NewClient(nil, os.Getenv("PLIVO_AUTH_ID"), os.Getenv("PLIVO_AUTH_TOKEN"))
	req := &plivo.CallRecordParams{TimeLimit: 60, FileFormat: "mp3", CallbackURL: os.Getenv("CALLBACK_URL"), CallbackMethod: "GET"}
	_, err := client.Call.Record(c.Query("CallUUID"), req)
	if err != nil {
		panic(err)
	}
}

func Callback(c *gin.Context) {
	response := &Response{}
	err := json.Unmarshal([]byte(c.Query("response")), response)
	if err == nil {
		fmt.Println("response", response.RecordURL)
	} else {
		panic(err)
	}
}
