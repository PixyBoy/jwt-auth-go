package ports

type RateLimiter interface {
	Allow(key string, max int, windowSeconds int) (bool, int, error)
}
