package cluster

import "github.com/boss-net/api/boss-plugin/internal/utils/cache"

func clearClusterState() {
	cache.Del(CLUSTER_STATUS_HASH_MAP_KEY)
	cache.Del(PREEMPTION_LOCK_KEY)
}
