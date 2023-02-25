package bot

import "errors"

var (
	errorNotifyDevelopers = errors.New("the request could not be processed. It may be a good idea to contact the developer")
)
