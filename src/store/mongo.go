package store

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoDBStore struct {
	ctx    context.Context
	dbName string
	uri    string
	db   *mongo.Database
}

const MAX_POOL = 100

func NewMongoDBStore(ctx context.Context,dbName string, uri string) *MongoDBStore {
	store := MongoDBStore{
		ctx:    ctx,
		dbName: dbName,
		uri:    uri,
	}
	store.init()
	return &store
}

func (store *MongoDBStore) init()  {
	client, err := mongo.Connect(store.ctx, options.Client().ApplyURI(store.uri).SetMaxPoolSize(MAX_POOL))
	if err != nil {
		log.Printf("Store error %x \n", err)
		return
	}
	store.db = client.Database(store.dbName)
}

func (store *MongoDBStore) Save(m interface{}, collection string) {
	res,err := store.db.Collection(collection).InsertOne(store.ctx, m)
	if err != nil {
		log.Printf("Store %s error %v \n", collection, err)
		return
	}
	log.Printf("insert %s %v", collection, res.InsertedID)
}
