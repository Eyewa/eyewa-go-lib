package amqp

// Get returns the value associated with the passed key.
func (c *HeaderCarrier) Get(key string) string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	c.gets = append(c.gets, key)
	return c.data[key]
}

// Set stores the key-value pair.
func (c *HeaderCarrier) Set(key, value string) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	c.sets = append(c.sets, [2]string{key, value})
	c.data[key] = value
}

// Keys returns the keys for which this carrier has a value.
func (c *HeaderCarrier) Keys() []string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	result := make([]string, 0, len(c.data))
	for k := range c.data {
		result = append(result, k)
	}
	return result
}
