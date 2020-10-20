package with

func StopChan(stop chan struct{}, fn func()) {
	for {
		select {
		case <-stop:
			return
		default:
			fn()
		}
	}
}
