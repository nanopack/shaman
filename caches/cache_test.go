package cache_test

type MapCacher map[string]string

func (self MapCacher) GetRecord(s string) (string, error) {
	return self[s], nil
}

func (self MapCacher) SetRecord(key, val string) error {
	self[key] = val
	return nil
}

func (self MapCacher) ReviseRecord(key, val string) error {
	return self.SetRecord(key, val)
}

func (self MapCacher) DeleteRecord(key string) error {
	delete(self, key)
}

func TestL1(t *testing.T) {
	
}