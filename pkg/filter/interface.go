package filter

type Filter interface {
	Use(message string) (isSpam bool)
}
