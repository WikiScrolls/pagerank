package database

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NewNeo4jClient(ctx context.Context, uri string, user string, password string) (neo4j.DriverWithContext, error) {
	driver, err := neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(user, password, ""),
	)
	if err != nil {
		return nil, err
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, err
	}

	return driver, nil
}
