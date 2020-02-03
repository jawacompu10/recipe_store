package repo

import (
	"context"
	"log"

	"github.com/jawacompu10/juice_shop/recipe_store/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo is the type that will implement the database layer for MongoDB
type Mongo struct {
	client *mongo.Client
	ctx    context.Context
	dbInfo DBInfo
}

// DBInfo stores the database connection details
type DBInfo struct {
	ConnectionString string
	DBName           string
	CollectionName   string
}

// New creates and returns a Mongo repo value
func New(dbInfo DBInfo) (*Mongo, error) {
	ctx := context.Background()
	client, err := mongo.NewClient(options.Client().ApplyURI(dbInfo.ConnectionString))
	if err != nil {
		return nil, err
	}
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to DB")

	return &Mongo{
		client: client,
		ctx:    ctx,
		dbInfo: dbInfo,
	}, nil
}

// GetRecipeByItem gets the recipe of a given item
func (mr *Mongo) GetRecipeByItem(item string) (models.Recipe, error) {
	recipe := models.Recipe{}
	filter := bson.M{
		"item_name": item,
	}

	err := mr.getCollection().FindOne(mr.ctx, filter).Decode(&recipe)
	return recipe, err
}

// AddRecipe adds a recipe to the recipe store
func (mr *Mongo) AddRecipe(recipe models.Recipe) (models.Recipe, error) {
	r, err := mr.getCollection().InsertOne(mr.ctx, recipe)
	if err != nil {
		log.Println("Failed to insert new recipe", err)
		return models.Recipe{}, err
	}
	recipe.ID = r.InsertedID.(primitive.ObjectID).Hex()
	return recipe, err
}

// UpdateRecipe updates an existing recipe
func (mr *Mongo) UpdateRecipe(recipe models.Recipe) (models.Recipe, error) {
	id, err := primitive.ObjectIDFromHex(recipe.ID)
	if err != nil {
		log.Println("Failed to convert ID to ObjectID", err)
		return recipe, err
	}
	filter := bson.M{
		"_id": id,
	}

	recipe.ID = ""
	_, err = mr.getCollection().ReplaceOne(mr.ctx, filter, recipe)
	return recipe, err
}

func (mr *Mongo) getCollection() *mongo.Collection {
	return mr.client.Database(mr.dbInfo.DBName).Collection(mr.dbInfo.CollectionName)
}
