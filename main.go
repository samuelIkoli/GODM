package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// genai "github.com/google/generative-ai-go/genai"
	controller "github.com/samuelIkoli/GODM/controllers"
	"github.com/samuelIkoli/GODM/internal/config"
	"github.com/samuelIkoli/GODM/routes"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/huggingface"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

// var client *mongo.Client

// Connect to MongoDB
func ConnectMongoDB() {
  // Create a new client and connect to the server
  	// uri := "mongodb+srv://Samuel:Layefanimi07@samcluster0.ezatj.mongodb.net/?retryWrites=true&w=majority&appName=SamCluster0"
  	// client, err := mongo.Connect( context.TODO(),options.Client().ApplyURI(uri))
  	
	// if err != nil {
    // 	panic(err)
  	// }
  	// defer func() {
    // 	if err = client.Disconnect(context.TODO()); err != nil {
    //   	panic(err)
    // 	}
  	// }()
		
	// collection := client.Database("sample_mflix").Collection("movies")
	// title := "Back to the Future"
	// var result bson.M
	// err = collection.FindOne(context.TODO(), bson.D{{Key:"title", Value: title}}).
	// 	Decode(&result)
	// if err == mongo.ErrNoDocuments {
	// 	fmt.Printf("No document was found with the title %s\n", title)
	// 	return
	// }
	// if err != nil {
	// 	panic(err)
	// }
	// jsonData, err := json.MarshalIndent(result, "", "    ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", jsonData)
}


func GeminiLangChain(c *gin.Context) {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	llm, err := googleai.New(ctx, googleai.WithAPIKey(apiKey))


	if err != nil {
		log.Fatal(err)
	}

  	prompt := "What is the L2 Lagrange point?"

	parts := []llms.ContentPart{
		llms.TextPart(prompt),
	}
	content := []llms.MessageContent{
		{
		Role:  "human",
		Parts: parts,
		},
	}

  answer, err := llm.GenerateContent(ctx, content, llms.WithModel("gemini-2.0-flash-lite"))
//   answer, err := llm.GenerateContent(ctx, []llms.MessageContent{})
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(answer)
  c.JSON(200,gin.H{"message": answer})
}
func HuggingFace(c *gin.Context) {
	ctx := context.Background()
  	// apiKey := os.Getenv("HUGGINGFACEHUB_API_TOKEN")
	fmt.Println()
  	llm, err := huggingface.New(huggingface.WithModel("google/t5-small"))
  	if err != nil {
    	log.Fatal(err)
  	}

  prompt := "What is the L2 Lagrange point?"
  answer, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}

  fmt.Println(answer)
  c.JSON(200,gin.H{"message": answer})
//   c.HTML(200, answer, "")
}



func main(){
	err := godotenv.Load(".env")
	if err != nil{
  		log.Fatalf("Error loading .env file: %s", err)
 	}
 	logger := config.NewAppLogger()
	// ConnectMongoDB()
	routes := routes.NewRoute(controller.NewController(logger))

	server:= gin.Default()

	server.GET("/gemini-langchain", GeminiLangChain)
	server.GET("/hugging-face", HuggingFace)

	routes.RegisterRoutes(server)

	server.Run(":8080")

}