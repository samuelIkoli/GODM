package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samuelIkoli/GODM/entity"
	"github.com/samuelIkoli/GODM/internal/config"
	"github.com/samuelIkoli/GODM/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri string = "mongodb+srv://Samuel:Layefanimi07@samcluster0.ezatj.mongodb.net/?retryWrites=true&w=majority&appName=SamCluster0"


type Controller struct {
	logger *config.Log
}

func NewController(log *config.Log) *Controller {
	return &Controller{
		logger: log,
	}
}

type Movie struct {
	ID      string   `bson:"_id"`
	Title   string   `bson:"title"`
	Plot    string   `bson:"plot"`
	EmbVec  []float32 `bson:"plot_embedding,omitempty"`
}

type User struct {
	ID      string   `bson:"_id"`
	Email   string   `bson:"email"`
	Username    string   `bson:"username"`
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

func (h *Controller) GenerateEmbeddings(c *gin.Context) {

  	client, err := mongo.Connect( context.TODO(),options.Client().ApplyURI(uri))
  	
	if err != nil {
    	panic(err)
  	}
  	defer func() {
    	if err = client.Disconnect(context.TODO()); err != nil {
      	panic(err)
    	}
  	}()
		
	collection := client.Database("sample_mflix").Collection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Fetch first 50 movies
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetLimit(50))
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var movies []Movie
	if err := cursor.All(ctx, &movies); err != nil {
		log.Fatal(err)
	}

	model:= services.InitGeminiClient()

	// Iterate through movies, generate embeddings, and update MongoDB
	for _, movie := range movies {
		if movie.Plot == "" {
			continue // Skip if no plot
		}

		embedding, err := services.GetGeminiEmbedding(model, movie.Plot)
		if err != nil {
			log.Printf("Error generating embedding for %s: %v\n", movie.Title, err)
			continue
		}

		objectID, err := primitive.ObjectIDFromHex(movie.ID)
		if err != nil {
    		log.Printf("Invalid ObjectID: %v\n", err)
    	return
		}

		// Update MongoDB with embedding
		_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"plot_embedding_hf": embedding}})
		if err != nil {
			log.Printf("Error updating movie %s: %v\n", movie.Title, err)
		} else {
			fmt.Printf("Updated: %s\n", movie.Title)
			// fmt.Printf("Updated: %s\n", embedding)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Embeddings generated and updated successfully"})
}

func (h *Controller) UpdateUsers(c *gin.Context) {

	uri := "mongodb+srv://Samuel:Layefanimi07@samcluster0.ezatj.mongodb.net/?retryWrites=true&w=majority&appName=SamCluster0"
  	client, err := mongo.Connect( context.TODO(),options.Client().ApplyURI(uri))
  	
	if err != nil {
    	panic(err)
  	}
  	defer func() {
    	if err = client.Disconnect(context.TODO()); err != nil {
      	panic(err)
    	}
  	}()
		
	collection := client.Database("yelp-camp").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Fetch first 50 users
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetLimit(10))
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var users []User
	if err := cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
	}

	// model:= services.InitGeminiClient()

	// Iterate through users, generate embeddings, and update MongoDB
	for _, user := range users {
		
		objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
    	log.Printf("Invalid ObjectID: %v\n", err)
    	return
	}

		// Update MongoDB with new field
		_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"gender": "male"}})
		if err != nil {
			log.Printf("Error updating user %s: %v\n", user.Username, err)
		} else {
			fmt.Printf("Updated: %s\n", user.Username)
			// fmt.Printf("Updated: %s\n", embedding)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Users updated successfully"})
}

func (h *Controller) VectorSearch(c *gin.Context){

	query := "District Attorney"
	model := services.InitGeminiClient()
	client, err := mongo.Connect( context.TODO(),options.Client().ApplyURI(uri))
	if err != nil {
		log.Printf("Error:%v\n", err)
	}
	collection := client.Database("sample_mflix").Collection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	query_embedding, err := services.GetGeminiEmbedding(model, query)
	// query_embedding = query_embedding[:768]
	log.Printf("Query embedding length: %d:", len(query_embedding))
	if err != nil {
		log.Printf("Error generating embedding for %s: %v\n", query, err)
	}

	pipeline := mongo.Pipeline{
        {{Key:"$vectorSearch", Value:bson.D{
            {Key:"index", Value:"newindex"}, // Use the name of your vector index
            {Key:"path", Value:"plot_embedding_hf"},
            {Key:"queryVector", Value:query_embedding},
            {Key:"numCandidates", Value:100}, // Number of items to consider before top-K selection
            {Key:"limit", Value:2},
        }}},
    	bson.D{{Key: "$project", Value: bson.D{
        	{Key: "title", Value: 1},
        	{Key: "plot", Value: 1},  // Include only title and plot
    }}},
    }

    opts := options.Aggregate().SetMaxTime(1000 * time.Second)
    cursor, err := collection.Aggregate(ctx, pipeline, opts)
	if cursor.RemainingBatchLength() == 0 {
		log.Println("No results found for vector search.")
		c.JSON(http.StatusOK, gin.H{"message": "No results found for vector search." })
		return
	}else{
		log.Printf("Length is %d",cursor.RemainingBatchLength())
	}
    if err != nil {
        log.Fatalf("Vector search failed: %v", err)
    }
    defer cursor.Close(ctx)

	var response []bson.M
    // Iterate over results
    for cursor.Next(ctx) {
        var result bson.M
        if err := cursor.Decode(&result); err != nil {
            fmt.Println("Error decoding result:", err)
            continue
        }
		fmt.Printf("Movie Name: %s,\nMovie Plot: %s\n\n", result["title"], result["plot"])
		response = append(response, result)
    }

	c.JSON(http.StatusOK, gin.H{"message": response })
}