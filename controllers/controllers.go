package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samuelIkoli/GODM/entity"
	"github.com/samuelIkoli/GODM/internal/config"
	"github.com/samuelIkoli/GODM/services"
)

type Controller struct {
	logger *config.Log
}

func NewController(log *config.Log) *Controller {
	return &Controller{
		logger: log,
	}
}

func Test(ctx *gin.Context){
	ctx.JSON(200, gin.H{
		"message": "test working with air hot reload and refactoring",
	})
}

func Task(ctx *gin.Context){
	now:= time.Now().UTC()
	result := entity.Response{
		Email: "ayibanimiikoli@gmail.com",
		Current_datetime: now.Format(time.RFC3339),
		Github_url: "https://github.com/samuelIkoli/HNG12_BE_0",
	}
	ctx.JSON(200, result)
}

func Ping(ctx *gin.Context){
	ctx.JSON(200, gin.H{
		"message": "Pong",
	})
}

func GetMessage(ctx *gin.Context){
	ctx.JSON(200, gin.H{
		"message": "This is the message from the DB",
	})
}

func (c *Controller) AIResponse(ctx *gin.Context) {
	startPrompt := "Checkout telex.im and summarize it for me please and also checkout and summarize who ELon Musk is. Does he and telex have similarities?"

	if startPrompt != "" {
		client := services.InitGeminiClient()

		answer, err := services.GetAIResponse(client, startPrompt)
		if err != nil {
			fmt.Println("Failed to process file: ", err)
			config.PrintLog(c.logger, "Failed to process file", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file", "detail": err.Error()})
			return
		}

		var parsedResponse map[string]interface{}
		parseErr := json.Unmarshal([]byte(answer), &parsedResponse)
		if parseErr != nil {
			// If unmarshalling fails, wrap it in a response struct
			parsedResponse = map[string]interface{}{"response": answer}
		}

		formattedAnswer, err := services.FormatResponse(parsedResponse)
		if err != nil {
			fmt.Println("Failed to format response: ", err)
			config.PrintLog(c.logger, "Failed to format response", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to format response", "detail": err.Error()})
			return
		}
		ctx.Data(http.StatusOK, "application/json", formattedAnswer)
	} else {
		config.PrintLog(c.logger, "Invalid data format", nil)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format", "detail": "You have not provided a valid message"})
		return
	}
}
