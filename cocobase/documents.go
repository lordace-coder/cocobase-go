package cocobase

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetDocument(ctx context.Context, collection, docID string) (*Document, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", collection, docID)
	
	resp, err := c.request(ctx, http.MethodGet, path, nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var doc Document
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &doc, nil
}

func (c *Client) CreateDocument(ctx context.Context, collection string, data map[string]interface{}) (*Document, error) {
	path := fmt.Sprintf("/collections/documents?collection=%s", collection)
	
	resp, err := c.request(ctx, http.MethodPost, path, data, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var doc Document
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &doc, nil
}

func (c *Client) UpdateDocument(ctx context.Context, collection, docID string, data map[string]interface{}) (*Document, error) {
	path := fmt.Sprintf("/collections/%s/documents/%s", collection, docID)
	
	resp, err := c.request(ctx, http.MethodPatch, path, data, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var doc Document
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &doc, nil
}

func (c *Client) DeleteDocument(ctx context.Context, collection, docID string) error {
	path := fmt.Sprintf("/collections/%s/documents/%s", collection, docID)
	
	resp, err := c.request(ctx, http.MethodDelete, path, nil, true)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) ListDocuments(ctx context.Context, collection string, query *QueryBuilder) ([]Document, error) {
	path := fmt.Sprintf("/collections/%s/documents", collection)
	
	if query != nil {
		queryStr := query.Build()
		if queryStr != "" {
			path += "?" + queryStr
		}
	}
	
	resp, err := c.request(ctx, http.MethodGet, path, nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var docs []Document
	if err := json.NewDecoder(resp.Body).Decode(&docs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return docs, nil
}

func (c *Client) QueryDocuments(ctx context.Context, collection, rawQuery string) ([]Document, error) {
	path := fmt.Sprintf("/collections/%s/documents", collection)
	
	if rawQuery != "" {
		path += "?" + rawQuery
	}
	
	resp, err := c.request(ctx, http.MethodGet, path, nil, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var docs []Document
	if err := json.NewDecoder(resp.Body).Decode(&docs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return docs, nil
}
