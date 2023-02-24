package bot

import "errors"

var (
	errorChannelIDMustBeSent = errors.New("write me the channel ID, anti-spam filters for which you need to configure")
	errorNotifyDevelopers    = errors.New("the request could not be processed. It may be a good idea to contact the developer")
)
