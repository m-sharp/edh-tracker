package format

func resetCache() {
	cache.Lock()
	defer cache.Unlock()
	cache.m = nil
}
