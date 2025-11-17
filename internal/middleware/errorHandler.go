package middleware

import (
	"context"
	"runtime"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"hongde_backend/internal/database"
)


func LogError(err error, errContext string) {
	if err == nil {
		return
	}

	// Capture the file and line number
	_, file, line, _ := runtime.Caller(1)

	errorCollection := database.DbMongo.Collection("error_lists")
	// Log activity
	_, _ = errorCollection.InsertOne(context.TODO(), bson.M{
		"timestamp": time.Now(),
		"error_message":err.Error(),
		"context":errContext,
		"affected_file":file,
		"line_number":line,
	})
}