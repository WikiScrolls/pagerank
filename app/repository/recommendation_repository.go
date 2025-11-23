package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type RecommendationRepository interface {
	GetRecommendationTitles(ctx context.Context, chainLength int) ([]string, error)
}

type Neo4jRecommendationRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jRecommendationRepository(driver neo4j.DriverWithContext) *Neo4jRecommendationRepository {
	return &Neo4jRecommendationRepository{driver: driver}
}

func (r *Neo4jRecommendationRepository) GetRecommendationTitles(ctx context.Context, chainLength int) ([]string, error) {
	query := fmt.Sprintf(
		`
		MATCH (n)
		WITH n, rand() AS r
		ORDER BY r
		LIMIT 1
		MATCH p = (n)-[*1..%d]-(m)
		WITH p, rand() AS r2
		ORDER BY r2
		LIMIT 1
		RETURN [x IN nodes(p) | x.title] AS pathTitles
		`, chainLength)

	result, err := neo4j.ExecuteQuery(ctx, r.driver,
		query,
		map[string]any{
			"chainLength": chainLength,
		},
		neo4j.EagerResultTransformer,
	)

	if err != nil {
		return nil, err
	}

	if len(result.Records) == 0 {
		return []string{}, nil
	}

	titles, _ := result.Records[0].Get("pathTitles")
	titleArray, isArray := titles.([]any)
	if !isArray {
		return []string{}, nil
	}

	var cleanTitles []string
	for _, title := range titleArray {
		titleString, isString := title.(string)
		if !isString || titleString == "" {
			continue
		}
		cleanTitles = append(cleanTitles, cleanTitle(titleString))
	}

	return cleanTitles, nil
}

func cleanTitle(title string) string {
	title = strings.Trim(title, "'")
	title = strings.ReplaceAll(title, " ", "_")
	if strings.Contains(title, "(") {
		title += ")"
	}
	return title
}
