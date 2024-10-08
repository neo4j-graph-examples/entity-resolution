// go mod init main
// go run example.go
package main

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"reflect"
)

func main() {
	results, err := runQuery("neo4j://<HOST>:<BOLTPORT>", "neo4j", "<USERNAME>", "<PASSWORD>")
	if err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}

func runQuery(uri, database, username, password string) (_ []string, err error) {
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	defer func() { err = handleClose(ctx, driver, err) }()
	query := "	MATCH (u:User {state: $state} )-[:WATCHED]->(m)-[:HAS]->(g:Genre)
	
	RETURN g.name as genre, count(g) as freq
	ORDER BY freq DESC
	params := map[string]any{"state": "Texas"}
	result, err := neo4j.ExecuteQuery(ctx, driver, query, params,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(database),
		neo4j.ExecuteQueryWithReadersRouting())
	if err != nil {
		return nil, err
	}
	genres := make([]string, len(result.Records))
	for i, record := range result.Records {
		// this assumes all actors have names, hence ignoring the 2nd returned value
		name, _, err := neo4j.GetRecordValue[string](record, "genre")
		if err != nil {
			return nil, err
		}
		genres[i] = name
	}
	return genres, nil
}

func handleClose(ctx context.Context, closer interface{ Close(context.Context) error }, previousError error) error {
	err := closer.Close(ctx)
	if err == nil {
		return previousError
	}
	if previousError == nil {
		return err
	}
	return fmt.Errorf("%v closure error occurred:\n%s\ninitial error was:\n%w", reflect.TypeOf(closer), err.Error(), previousError)
}
