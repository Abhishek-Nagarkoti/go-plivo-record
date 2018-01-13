package main

import (
	"encoding/json"
	// "encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/plivo/plivo-go/plivo"
	"github.com/plivo/plivo-go/plivo/xml"
	"os"
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
	r.GET("/play", Play)
	r.GET("/record", Record)
	r.GET("/plivo/callback", Callback)
	r.POST("/", Create)
	r.GET("/hangup", Hangup)

	// Start listening
	port := Port
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}
	r.Run(port)
}

func Get(c *gin.Context) {

	action := "http://url_where_to_redirect_when_recording_starts"
	response := xml.ResponseElement{
		Contents: []interface{}{
			new(xml.GetDigitsElement).SetAction(action).SetMethod("GET").SetRedirect(false).SetRetries(1).SetNumDigits(1).SetValidDigits("1").SetTimeout(15).SetContents([]interface{}{
				new(xml.SpeakElement).SetContents("Press 1 to record a message."),
			}),
			new(xml.WaitElement).SetLength(10),
		},
	}
	c.XML(200, response)
}

func Play(c *gin.Context) {
	type GetDigits struct {
		Redirect  string `xml:"redirect,attr"`
		Retries   string `xml:"retries,attr"`
		Method    string `xml:"method,attr"`
		NumDigits string `xml:"numDigits,attr"`
		Action    string `xml:"action,attr"`
		Timeout   string `xml:"timeout,attr"`
		Speak     string `xml:"Speak"`
	}

	type Wait struct {
		Length string `xml:"length,attr"`
	}

	type Response struct {
		XMLName   xml.Name `xml:"Response"`
		GetDigits GetDigits
		Wait      Wait
	}

	action := "http://url_where_to_redirect_where_mp3_plays"

	response := &Response{}
	response.GetDigits = GetDigits{Redirect: "false", Retries: "1", Method: "GET", NumDigits: "1", Action: action, Timeout: "7", Speak: "Press 1 to hear your message."}
	response.Wait = Wait{Length: "10"}

	c.XML(200, response)
}

func Create(c *gin.Context) {
	client, err := plivo.NewClient(os.Getenv("PLIVO_AUTH_ID"), os.Getenv("PLIVO_AUTH_TOKEN"), &plivo.ClientOptions{})
	if err != nil {
		panic(err)
	}
	client.Calls.Create(plivo.CallCreateParams{
		From:             os.Getenv("PHONE_FROM"),
		To:               os.Getenv("PHONE_TO"),
		AnswerURL:        os.Getenv("ANSWER_URL"),
		AnswerMethod:     "GET",
		HangupURL:        os.Getenv("HANGUP_URL"),
		HangupMethod:     "GET",
		TimeLimit:        60,
		RingTimeout:      60,
		MachineDetection: "hangup",
	})
}

func Record(c *gin.Context) {
	client, err := plivo.NewClient(os.Getenv("PLIVO_AUTH_ID"), os.Getenv("PLIVO_AUTH_TOKEN"), &plivo.ClientOptions{})
	if err != nil {
		panic(err)
	}
	client.Calls.Record(c.Query("CallUUID"), plivo.CallRecordParams{
		FileFormat:     "mp3",
		CallbackURL:    os.Getenv("CALLBACK_URL"),
		CallbackMethod: "GET",
		TimeLimit:      60,
	})
}

func Callback(c *gin.Context) {
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
	response := &Response{}
	err := json.Unmarshal([]byte(c.Query("response")), response)
	if err == nil {
		fmt.Println("response", response.RecordURL)
	} else {
		panic(err)
	}
}

func Hangup(c *gin.Context) {
	fmt.Println("TotalCost", c.Query("TotalCost"))
	fmt.Println("Direction", c.Query("Direction"))
	fmt.Println("HangupCause", c.Query("HangupCause"))
	fmt.Println("From", c.Query("From"))
	fmt.Println("BillDuration", c.Query("BillDuration"))
	fmt.Println("BillRate", c.Query("BillRate"))
	fmt.Println("To", c.Query("To"))
	fmt.Println("RequestUUID", c.Query("RequestUUID"))
	fmt.Println("Duration", c.Query("Duration"))
	fmt.Println("CallUUID", c.Query("CallUUID"))
	fmt.Println("EndTime", c.Query("EndTime"))
	fmt.Println("CallStatus", c.Query("CallStatus"))
	fmt.Println("Event", c.Query("Event"))
}
