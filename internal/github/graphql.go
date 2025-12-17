package github

import (
	"context"
	"encoding/json"
	"fmt"
)

// GraphQLクエリを実行する共通メソッド
func (c *Client) executeGraphQL(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	req := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: variables,
	}

	// HTTPリクエストを作成
	httpReq, err := c.client.NewRequest("POST", "graphql", req)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// GraphQLレスポンス構造
	var graphQLResp struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	// リクエスト実行
	_, err = c.client.Do(ctx, httpReq, &graphQLResp)
	if err != nil {
		return fmt.Errorf("failed to execute GraphQL query: %w", err)
	}

	// GraphQLエラーチェック
	if len(graphQLResp.Errors) > 0 {
		return fmt.Errorf("GraphQL error: %s", graphQLResp.Errors[0].Message)
	}

	// 結果をUnmarshal
	if err := json.Unmarshal(graphQLResp.Data, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
