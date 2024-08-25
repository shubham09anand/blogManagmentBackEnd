package connection

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connection struct holds the MongoDB client, context, cancel function, and port number
type Connection struct {
	Client *mongo.Client      // MongoDB client instance
	Ctx    context.Context    // Context for managing request lifetime
	Cancel context.CancelFunc // Function to cancel the context
	Port   string             // Port number MongoDB is running on
}

// ConnectDB initializes a new MongoDB connection and returns a Connection struct
func ConnectDB() (*Connection, error) {

	urlStr := "mongodb://3.6.164.210:27017"

	// urlStr := databaseURL
	port := "27017"

	// Create a new MongoDB client with the specified URI
	clientOptions := options.Client().ApplyURI(urlStr)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		// If client creation fails, return an error
		return nil, err
	}

	// Create a context with a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is cancelled when done

	// Ping the MongoDB server to verify the connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		// If ping fails, attempt to disconnect the client and return an error
		client.Disconnect(ctx)
		return nil, err
	}

	// Print the port number and confirm the connection
	fmt.Printf("Connected to MongoDB on port %s\n", port)

	// Return a Connection struct with the client, context, cancel function, and port number
	return &Connection{Client: client, Ctx: ctx, Cancel: cancel, Port: port}, nil
}

// Close terminates the MongoDB connection and cancels the context
func Close(conn *Connection) {
	if conn != nil {
		// Disconnect the MongoDB client
		if err := conn.Client.Disconnect(conn.Ctx); err != nil {
			// Log an error if disconnect fails
			fmt.Println("Error while disconnecting:", err)
		}
		// Cancel the context to clean up resources
		fmt.Println("Connection Closed")
		conn.Cancel()
	}
}
