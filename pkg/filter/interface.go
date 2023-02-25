package filter

type Filter interface {
	Use(message string) (isDetected bool)
	GetTag() string
	GetName() string
}
