package redis

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// add k & v to redis
func add(c *redis.Client, k string, v int64) {
	tx := c.TxPipeline()
	tx.SAdd(ctx, k, v)
	tx.Expire(ctx, k, expireTime)
	tx.Exec(ctx)
}

// del k & v
func del(c *redis.Client, k string, v int64) {
	tx := c.TxPipeline()
	tx.SRem(ctx, k, v)
	tx.Expire(ctx, k, expireTime)
	tx.Exec(ctx)
}

// check if the set of k exists
func check(c *redis.Client, k string) bool {
	if e, _ := c.Exists(ctx, k).Result(); e > 0 {
		return true
	}
	return false
}

// exist checks if the relation k and v exists
func exist(c *redis.Client, k string, v int64) bool {
	if e, _ := c.SIsMember(ctx, k, v).Result(); e {
		c.Expire(ctx, k, expireTime)
		return true
	}
	return false
}

// count gets the size of the set of key
func count(c *redis.Client, k string) (sum int64, err error) {
	if sum, err = c.SCard(ctx, k).Result(); err == nil {
		c.Expire(ctx, k, expireTime)
		return sum, err
	}
	return sum, err
}

// get gets the elements from the set of key
func get(c *redis.Client, k string) (vt []int64) {
	v, _ := c.SMembers(ctx, k).Result()
	c.Expire(ctx, k, expireTime)
	for _, vs := range v {
		v_i64, _ := strconv.ParseInt(vs, 10, 64)
		vt = append(vt, v_i64)
	}
	return vt
}
