// gomongoclient implements and simplify the mongodb-driver usage
package gomongoclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	client   mongo.Client
}

// Init the basic connection
func (m *MongoDB) Init() (err error) {

	if m.client, err = m.connect(); err != nil {
		fmt.Printf("[ERROR]: Can't connect to the host: [%s]. [Err Detail: %s]\n ", m.Host, err)
	}

	return
}

// Done close all connections to this stance and free all the cache
func (m *MongoDB) Done() (err error) {

	if err = m.client.Disconnect(context.TODO()); err != nil {
		fmt.Printf("[ERROR]: No connection to close. [Err Detail: %s]\n", err)
	}

	return
}

// connect to the mongo databse
func (m *MongoDB) connect() (mongo.Client, error) {

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%v", m.User, m.Password, m.Host, m.Port)

	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		fmt.Printf("failed to connect to the mongo server %v", err)
	}

	// Test Connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())

	return *client, err

}

func (m *MongoDB) Test() string {

	return "test"
}

// Ping tests the connection with the mongo Server
func (m *MongoDB) Ping(connectOnFailure bool) (err error) {

	if err = m.client.Ping(context.TODO(), nil); err != nil {

		if connectOnFailure {
			err = m.Init()
		}
	}

	return
}

// Find the values from given query
func (m *MongoDB) Find(collection string, query interface{}, limit int64) (result []bson.M, err error) {

	if err = m.Ping(true); err != nil {
		return
	}

	options := options.Find()
	options.SetLimit(limit)
	col := m.client.Database(m.Database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err := col.Find(ctx, query, options)

	if err != nil {
		log.Printf("[ERROR]: Failed to execute the search query. [Err Detail: %s]", err)
		return
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {

		var res bson.M
		err = cur.Decode(&res)

		result = append(result, res)

		if err != nil {
			log.Printf("[ERROR]: ToDo1 error message. [Err Detail: %s]", err)
			return
		}
	}

	return

}

// FindAndConvert return the mongo cursor to use in struct conversion
func (m *MongoDB) FindAndConvert(collection string, query interface{}, sort interface{}, limit int64) (cur *mongo.Cursor, err error) {

	if err = m.Ping(true); err != nil {
		return
	}

	options := options.Find()
	options.SetLimit(limit)
	options.SetSort(sort)
	col := m.client.Database(m.Database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cur, err = col.Find(ctx, query, options)

	if err != nil {
		log.Printf("[ERROR]: Failed to execute the search query. [Err Detail: %s]", err)
	}

	return

}

// Insert a new Value to the Collection
func (m *MongoDB) Insert(collection string, query interface{}) error {

	if err := m.Ping(true); err != nil {
		return err
	}

	col := m.client.Database(m.Database).Collection(collection)

	// TODO - TTL Index
	// mod := mongo.IndexModel{
	// 	Keys: bson.M{
	// 		"ExpireAt": 1, // index in ascending order
	// 	}, Options: options.Index().SetExpireAfterSeconds(0),
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ind, err := col.Indexes().CreateOne(ctx, mod)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(ind)

	_, err := col.InsertOne(ctx, query)

	if err != nil {
		return fmt.Errorf("[MONGO]: Insert failed. [Err Detail: %s]", err)
	}

	return nil
}

func (m *MongoDB) Upsert(collection string, values bson.M, filter bson.M) error {

	if err := m.Ping(true); err != nil {
		return err
	}

	col := m.client.Database(m.Database).Collection(collection)
	opts := options.Update()
	opts.SetUpsert(true)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := col.UpdateOne(ctx, filter, values, opts)

	if err != nil {
		return fmt.Errorf("[MONGO]: Upsert failed. [Err Detail: %s]", err)
	}

	return nil
}

func (m *MongoDB) Update(collection string, filter bson.M, update bson.M) {

	if err := m.Ping(true); err != nil {
		return
	}

	col := m.client.Database(m.Database).Collection(collection)

	_, err := col.UpdateMany(context.TODO(), filter, update)

	if err != nil {
		fmt.Println(err)
	}

}

// UpsertMany records on the mongo database
func (m *MongoDB) UpsertMany(collection string, query interface{}, filter interface{}) {

	if err := m.Ping(true); err != nil {
		return
	}

	col := m.client.Database(m.Database).Collection(collection)
	var operations []mongo.WriteModel

	operation := mongo.NewUpdateOneModel()
	operation.SetFilter(filter)
	operation.SetUpdate(bson.M{"$set": query})
	operation.SetUpsert(true)

	operations = append(operations, operation)

	// Specify an option to turn the bulk insertion in order of operation
	bulkOption := options.BulkWriteOptions{}
	// bulkOption.SetOrdered(true)

	_, err := col.BulkWrite(context.TODO(), operations, &bulkOption)
	if err != nil {
		(fmt.Printf("Mongo failed to Upsert Many, %v\n", err))

	}

}

// InsertMany values to the given collection
func (m *MongoDB) InsertMany(collection string, query []interface{}) {

	if err := m.Ping(true); err != nil {
		return
	}

	col := m.client.Database(m.Database).Collection(collection)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := col.InsertMany(ctx, query)

	if err != nil {
		fmt.Printf("[MONGO]: Insert failed. [Err Detail: %s]", err)
	}
}
