package chain

import (
	"encoding/json"
	"os"
)

// Save writes the markov chain to a JSON model file.
func (c *Chain) Save(path string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := json.Marshal(c.Data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads the markov chain from a JSON model file.
func (c *Chain) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return json.Unmarshal(data, &c.Data)
}
