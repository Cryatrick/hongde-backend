package service

import(
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"hongde_backend/internal/model"
	"hongde_backend/internal/database"
	"hongde_backend/internal/middleware"

)

func UpsertTokenData(userId string, dbData bson.M) (bool) {
	opts := options.UpdateOne().SetUpsert(true)
	_, err := database.DbMongo.Collection("access_tokens").UpdateOne(context.TODO(),
		bson.M{"user_id":userId},
		bson.M{"$set" : dbData},
		opts,
	)
	if err != nil {
		middleware.LogError(err,"Mongo DB Failed to Save Data")
		return false
	}
	return true
}
func DeleteTokenData(userId string) (bool) {
	_,err := database.DbMongo.Collection("access_tokens").DeleteOne(context.TODO(),bson.M{"user_id":userId})
	if err != nil {
		middleware.LogError(err,"Mongo DB Failed to Delete Data")
		return false
	}
	return true
}
func GetTokenData(whereParam bson.M) (*model.TokenData) {
	var storedToken *model.TokenData
	err := database.DbMongo.Collection("access_tokens").FindOne(context.TODO(), whereParam).Decode(&storedToken)
	if err != nil {
		// middleware.LogError(err,"User Token Not Found")
		return nil
	}
	return storedToken
}

func TestPing() {
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := database.DbMongo.RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}