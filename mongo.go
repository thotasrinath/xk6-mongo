package xk6_mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	k6modules "go.k6.io/k6/js/modules"
)

// Register the extension on module initialization, available to
// import from JS as "k6/x/mongo".
func init() {
	k6modules.Register("k6/x/mongo", new(Mongo))
}

// Mongo is the k6 extension for a Mongo client.
type Mongo struct{}

// Client is the Mongo client wrapper.
type Client struct {
	client *mongo.Client
}

// NewClient represents the Client constructor (i.e. `new mongo.Client()`) and
// returns a new Mongo client object.
// connURI -> mongodb://username:password@address:port/db?connect=direct
func (*Mongo) NewClient(connURI string) interface{} {

	clientOptions := options.Client().ApplyURI(connURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	return &Client{client: client}

}

func (c *Client) Insert(database string, collection string, doc interface{}) error {
	db := c.client.Database(database)
	col := db.Collection(collection)

	_, err := col.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ConvertToEjsonAndInsert(database string, collection string, doc interface{}) error {
	db := c.client.Database(database)
	col := db.Collection(collection)

	parsedDocBytes, err := bson.MarshalExtJSON(doc, false, false)
	if err != nil {
		return err
	}
	var parsedDoc interface{}
	err = bson.UnmarshalExtJSON(parsedDocBytes, false, &parsedDoc)
	if err != nil {
		return err
	}
	_, err = col.InsertOne(context.TODO(), parsedDoc)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) InsertBatch(database string, collection string, docs []any) error {

	db := c.client.Database(database)
	col := db.Collection(collection)
	_, err := col.InsertMany(context.TODO(), docs)
	if err != nil {
		return err
	}
	return nil

}

func (c *Client) Upsert(database string, collection string, docId string, doc map[string]interface{}) error {
	db := c.client.Database(database)
	col := db.Collection(collection)

	filter := bson.D{{Key: "_id", Value: docId}}

	var upsertDoc bson.D
	for k, v := range doc {
		upsertDoc = append(upsertDoc, bson.E{Key: k, Value: v})

	}

	update := bson.D{{Key: "$set", Value: upsertDoc}}
	opts := options.Update().SetUpsert(true)
	_, err := col.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Find(database string, collection string, filter interface{}, limit int64) ([]bson.M, error) {
	var results []bson.M

	db := c.client.Database(database)
	col := db.Collection(collection)

	opts := options.Find().SetLimit(limit)
	cur, err := col.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
		return results, err
	}

	if err = cur.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return results, err
}

func (c *Client) FindOne(database string, collection string, filter map[string]string) (bson.M, error) {
	db := c.client.Database(database)
	col := db.Collection(collection)
	var result bson.M
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: 1}})
	err := col.FindOne(context.TODO(), filter, opts).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result, nil
}
