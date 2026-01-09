package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RyanSikandar/orders-api/model"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	// Redis Client
	Client *redis.Client
}

func generateOrderKey(id int) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Create(ctx context.Context, order model.Order) error {
	// We marshal the order struct

	data, err := json.Marshal(order)

	if err != nil {
		return fmt.Errorf("failed to encode the order: %w", err)
	}

	key := generateOrderKey(order.ID)

	// We use a transaction to ensure atomicity
	txn := r.Client.TxPipeline()

	// Save the order data in Redis
	err = txn.SetNX(ctx, key, string(data), 0).Err()

	if err != nil {
		txn.Discard()
		return fmt.Errorf("failed to insert order into redis: %w", err)
	}

	// We also add the order ID to a set for listing purposes
	err = txn.SAdd(ctx, "orders", key).Err()

	if err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add order ID to orders set: %w", err)
	}

	_, err = txn.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	return nil
}

func (r *RedisRepo) GetByID(ctx context.Context, id int) (model.Order, error) {
	var order model.Order

	key := generateOrderKey(id)

	data, err := r.Client.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return model.Order{}, fmt.Errorf("order with ID %d not found", id)
		}
		return model.Order{}, fmt.Errorf("failed to get order from redis: %w", err)
	}

	err = json.Unmarshal([]byte(data), &order) // we use pointer so that unmarshal can modify the original order variable

	if err != nil {
		return model.Order{}, fmt.Errorf("failed to decode order data: %w", err)
	}

	return order, nil
}

func (r *RedisRepo) UpdateByID(ctx context.Context, id int, updatedOrder model.Order) error {
	key := generateOrderKey(id)

	// Check if the order exists
	_, err := r.Client.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("order with ID %d not found", id)
		}
		return fmt.Errorf("failed to get order from redis: %w", err)
	}

	data, err := json.Marshal(updatedOrder)
	
	if err != nil {
		return fmt.Errorf("failed to encode the updated order: %w", err)
	}
	// Update the order data in Redis
	err = r.Client.Set(ctx, key, string(data), 0).Err()

	if err != nil {
		return fmt.Errorf("failed to update order in redis: %w", err)
	}

	return nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, id int) error {
	key := generateOrderKey(id)

	//We use a transaction to ensure atomicity
	txn := r.Client.TxPipeline()

	// Delete the order from Redis
	result, err := txn.Del(ctx, key).Result()

	if err != nil {
		txn.Discard()
		return fmt.Errorf("failed to delete order from redis: %w", err)
	}

	if result == 0 {
		return fmt.Errorf("order with ID %d not found", id)
	}

	// Remove the order ID from the set
	err = txn.SRem(ctx, "orders", key).Err()
	
	if err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove order ID from orders set: %w", err)
	}

	_, err = txn.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	return nil
}

type OrderIterator struct {
	Size int64
	Offset uint64
}

type FindResults struct {
	Orders []model.Order
	cursor int
}

func (r *RedisRepo) List(ctx context.Context, page OrderIterator) (FindResults, error){
	res := r.Client.SScan(ctx, "orders", page.Offset, "*", page.Size)

	keys, cursor, err := res.Result()

	if err != nil {
		return FindResults{}, fmt.Errorf("failed to scan orders set: %w", err)
	}

	if len(keys) == 0 {
		return FindResults{
			Orders: []model.Order{},
			cursor: int(cursor),
		}, nil
	}
	xs, err := r.Client.MGet(ctx, keys...).Result()

	if err != nil {
		return FindResults{}, fmt.Errorf("failed to get orders from redis: %w", err)
	}

	orders := make([]model.Order, 0, len(xs))

	for _, x := range xs {
		if x == nil {
			continue
		}
		
		var order model.Order
		err := json.Unmarshal([]byte(x.(string)), &order)
		if err != nil {
			return FindResults{}, fmt.Errorf("failed to decode order data: %w", err)
		}
		orders = append(orders, order)
	}
	
	return FindResults{
		Orders: orders,
		cursor: int(cursor),
	}, nil
}