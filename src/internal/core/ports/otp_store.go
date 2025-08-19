package ports

type OTPStore interface {
	Save(phone, hash string, ttlSeconds int) error
	Get(phone string) (hash string, attempts int, exists bool, err error)
	IncreaseAttempt(phone string) (int, error)
	Delete(phone string) error
}
