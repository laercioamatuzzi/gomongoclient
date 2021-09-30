# gomongoclient

- gomongoclient implements and simplify the mongodb-driver usage

### Todos

- 
- 


### Usage example

```go
	mongo := gomongoclient.MongoDB{
		Host:     "127.0.0.1",
		User:     "localhost",
		Password: "example",
		Port:     27017,
		Database: "root",
	}

    // Connect to the mongo Database.
    err := mongo.Init()

	if err != nil {
		log.Fatal(err)
	}

    // Closing the connection instance.
	defer mongo.Done()

    // Insert
	collection := "Test"
	query := bson.M{"key": "value"}
	err = mongo.Insert(collection, query)

    // Select
    filter := bson.M{}
	sort := bson.M{}

	cursor, err := mongo.FindAndConvert("Test", filter, sort, 50)

    if err != nil {
		log.Fatal(err)
	}

	type Test struct {
		ObjectID string   `bson:"_id"`
		Value    string `bson:"value"`
	}

	result := []bson.M{}

	for cursor.Next(context.TODO()) {

		var t Test
		cursor.Decode(&t)

        fmt.Println(t)
	}
```