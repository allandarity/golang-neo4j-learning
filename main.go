package main

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type neo4jConfiguration struct {
	url      string
	username string
	password string
	database string
}

func (config *neo4jConfiguration) createDriver() (neo4j.DriverWithContext, error) {
	return neo4j.NewDriverWithContext(config.url, neo4j.BasicAuth(config.username, config.password, "")) //TODO: set these from env variables
}

func main() {
	ctx := context.Background()
	config := &neo4jConfiguration{
		url:      "bolt://localhost:7687",
		username: "neo4j",
		password: "password",
		database: "",
	}

	driver, err := config.createDriver()
	if err != nil {
		fmt.Println(err)
	}
	defer driver.Close(ctx)
	//matchAll(driver, config.database, "Person")
	example(driver)
}

func matchAll(driver neo4j.DriverWithContext, database string, node string) {
	queryString := fmt.Sprintf("MATCH (n: %s) return *", node)
	result, err := neo4j.ExecuteQuery(context.Background(), driver, queryString, nil, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(database), neo4j.ExecuteQueryWithReadersRouting())

	if err != nil {
		fmt.Println(err)
	}

	for _, record := range result.Records {
		fmt.Println(record.AsMap())
		x, _ := record.Get("from")
		fmt.Println(x)
	}
}

func example(driver neo4j.DriverWithContext) {
	ctx := context.Background()
	// Prepare data
	people := []map[string]any{
		{"name": "Alice", "age": 42, "friends": []string{"Bob", "Peter", "Anna"}},
		{"name": "Bob", "age": 19},
		{"name": "Peter", "age": 50},
		{"name": "Anna", "age": 30},
	}

	// Create some nodes
	for _, person := range people {
		_, err := neo4j.ExecuteQuery(ctx, driver,
			"MERGE (p:Person {name: $person.name, age: $person.age})",
			map[string]any{
				"person": person,
			}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))
		if err != nil {
			panic(err)
		}
	}

	// Create some relationships
	for _, person := range people {
		if person["friends"] != "" {
			_, err := neo4j.ExecuteQuery(ctx, driver, `
                MATCH (p:Person {name: $person.name})
                UNWIND $person.friends AS friend_name
                MATCH (friend:Person {name: friend_name})
                MERGE (p)-[:KNOWS]->(friend)
                `, map[string]any{
				"person": person,
			}, neo4j.EagerResultTransformer,
				neo4j.ExecuteQueryWithDatabase("neo4j"))
			if err != nil {
				panic(err)
			}
		}
	}

	// Retrieve Alice's friends who are under 40
	result, err := neo4j.ExecuteQuery(ctx, driver, `
        MATCH (p:Person {name: $name})-[:KNOWS]-(friend:Person)
        WHERE friend.age < $age
        RETURN friend
        `, map[string]any{
		"name": "Alice",
		"age":  40,
	}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		panic(err)
	}

	// Loop through results and do something with them
	for _, record := range result.Records {
		fmt.Println(record.AsMap())
	}

	// Summary information
	fmt.Printf("\nThe query `%v` returned %v records in %+v.\n",
		result.Summary.Query().Text(), len(result.Records),
		result.Summary.ResultAvailableAfter())
}
