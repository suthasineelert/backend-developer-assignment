package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// RedisClient is a mock for the Redis client
type RedisClient struct {
	mock.Mock
}

// Get mocks the Get method
func (m *RedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

// Set mocks the Set method
func (m *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *RedisClient) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}