package main

var minShards = 1024

type cache struct {
	shards []*cacheShard //分片
	hash   fnv64a        //hash算法
}

func newCache() *cache {
	cache := &cache{
		hash:   newDefaultHasher(),
		shards: make([]*cacheShard, minShards),
	}

	for i := 0; i < minShards; i++ {
		cache.shards[i] = initNewShard()
	}
	return cache
}

func (c *cache) getShard(hashedKey uint64) (shard *cacheShard) {
	return c.shards[hashedKey&uint64(minShards-1)]
}

func (c *cache) set(key string, value []byte) {
	hashedKey := c.hash.Sum64(key)
	shard := c.getShard(hashedKey)
	shard.set(hashedKey, value)
}

func (c *cache) get(key string) ([]byte, error) {
	hashKey := c.hash.Sum64(key)
	shard := c.getShard(hashKey)

	return shard.get(key, hashKey)
}
