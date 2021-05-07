package rds_cache_go

import "time"

//	@method TTL
//	@description: 过期时间校验
//	@param ttl time.Duration
//	@return time.Duration
func TTL(ttl time.Duration) time.Duration {
	if ttl <= 0 {
		return 0
	}

	return ttl
}
