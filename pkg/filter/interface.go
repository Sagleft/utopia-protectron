package filter

type Filter interface {
	Use(message string) (isDetected bool)
}
