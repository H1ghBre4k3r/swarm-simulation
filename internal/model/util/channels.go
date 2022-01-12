package util

func IsChannelClosed(ch chan string) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
