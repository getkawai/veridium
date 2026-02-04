package store

import (
	"context"
	"fmt"
	"io"

	cfv6 "github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/kv"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

// KVClient wraps Cloudflare KV operations with v6 SDK
type KVClient struct {
	client    *cfv6.Client
	accountID string
}

// NewKVClient creates a new KV client using v6 SDK
func NewKVClient(apiToken, accountID string) (*KVClient, error) {
	if apiToken == "" {
		return nil, fmt.Errorf("API token is required")
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID is required")
	}

	client := cfv6.NewClient(
		option.WithAPIToken(apiToken),
	)

	return &KVClient{
		client:    client,
		accountID: accountID,
	}, nil
}

// KeyInfo represents a KV key metadata
type KeyInfo struct {
	Name       string
	Expiration int64
	Metadata   interface{}
}

// ListResult represents paginated list result
type ListResult struct {
	Result     []KeyInfo
	ResultInfo ResultInfo
}

// ResultInfo contains pagination information
type ResultInfo struct {
	Cursor string
	Count  int
}

// GetValue retrieves a value from KV store
func (c *KVClient) GetValue(ctx context.Context, namespaceID, key string) ([]byte, error) {
	resp, err := c.client.KV.Namespaces.Values.Get(ctx, namespaceID, key, kv.NamespaceValueGetParams{
		AccountID: cfv6.F(c.accountID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get KV value: %w", err)
	}

	// Read the response body
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// SetValue writes a value to KV store
func (c *KVClient) SetValue(ctx context.Context, namespaceID, key string, value []byte) error {
	_, err := c.client.KV.Namespaces.Values.Update(ctx, namespaceID, key, kv.NamespaceValueUpdateParams{
		AccountID: cfv6.F(c.accountID),
		Value:     cfv6.F(string(value)),
	})
	if err != nil {
		return fmt.Errorf("failed to set KV value: %w", err)
	}
	return nil
}

// SetValueWithTTL writes a value to KV store with expiration
func (c *KVClient) SetValueWithTTL(ctx context.Context, namespaceID, key string, value []byte, ttlSeconds int) error {
	_, err := c.client.KV.Namespaces.Values.Update(ctx, namespaceID, key, kv.NamespaceValueUpdateParams{
		AccountID:     cfv6.F(c.accountID),
		Value:         cfv6.F(string(value)),
		ExpirationTTL: cfv6.F(float64(ttlSeconds)),
	})
	if err != nil {
		return fmt.Errorf("failed to set KV value with TTL: %w", err)
	}
	return nil
}

// DeleteValue removes a value from KV store
func (c *KVClient) DeleteValue(ctx context.Context, namespaceID, key string) error {
	_, err := c.client.KV.Namespaces.Values.Delete(ctx, namespaceID, key, kv.NamespaceValueDeleteParams{
		AccountID: cfv6.F(c.accountID),
	})
	if err != nil {
		return fmt.Errorf("failed to delete KV value: %w", err)
	}
	return nil
}

// ListKeys lists keys in a namespace with optional prefix and cursor
func (c *KVClient) ListKeys(ctx context.Context, namespaceID, prefix, cursor string) (*ListResult, error) {
	params := kv.NamespaceKeyListParams{
		AccountID: cfv6.F(c.accountID),
	}

	if prefix != "" {
		params.Prefix = cfv6.F(prefix)
	}

	if cursor != "" {
		params.Cursor = cfv6.F(cursor)
	}

	resp, err := c.client.KV.Namespaces.Keys.List(ctx, namespaceID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list KV keys: %w", err)
	}

	// Convert response to our format
	result := &ListResult{
		Result: make([]KeyInfo, 0, len(resp.Result)),
		ResultInfo: ResultInfo{
			Count: len(resp.Result),
		},
	}

	for _, key := range resp.Result {
		result.Result = append(result.Result, KeyInfo{
			Name:       key.Name,
			Expiration: int64(key.Expiration),
			Metadata:   key.Metadata,
		})
	}

	// Handle pagination cursor from v6 SDK
	// The cursor is in ResultInfo.Cursors.After field
	if resp.ResultInfo.Cursors.After != "" {
		result.ResultInfo.Cursor = resp.ResultInfo.Cursors.After
	}

	return result, nil
}

// ListKeysSimple is a convenience wrapper that returns just the key names
// Note: This only returns the first page. Use ListAllKeys for complete results.
func (c *KVClient) ListKeysSimple(ctx context.Context, namespaceID, prefix string) ([]string, error) {
	result, err := c.ListKeys(ctx, namespaceID, prefix, "")
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(result.Result))
	for _, key := range result.Result {
		keys = append(keys, key.Name)
	}

	return keys, nil
}

// ListAllKeys returns all keys with pagination support
func (c *KVClient) ListAllKeys(ctx context.Context, namespaceID, prefix string) ([]string, error) {
	var allKeys []string
	cursor := ""

	for {
		result, err := c.ListKeys(ctx, namespaceID, prefix, cursor)
		if err != nil {
			return nil, err
		}

		for _, key := range result.Result {
			allKeys = append(allKeys, key.Name)
		}

		// Check if there are more pages
		if result.ResultInfo.Cursor == "" {
			break
		}
		cursor = result.ResultInfo.Cursor
	}

	return allKeys, nil
}
