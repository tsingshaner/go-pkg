package formatter

func getLastGroup(log Log) Group {
	return log.Groups()[len(log.Groups())-1]
}
