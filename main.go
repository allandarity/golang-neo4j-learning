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
		password: "test1234",
		database: "",
	}

	driver, err := config.createDriver()
	if err != nil {
		fmt.Println(err)
	}
	defer driver.Close(ctx)

	result, err := neo4j.ExecuteQuery(
		ctx, driver, "testststest", map[string]interface{}{"limit": 5}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithReadersRouting(), neo4j.ExecuteQueryWithDatabase(config.database))

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	fmt.Println("hello")
}
