package memory

type Memory interface {
	IsUserExists(u User) (bool, error)
	SaveUser(u User) error
}
