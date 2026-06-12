package chain

import (
	"log/slog"
	"math/rand"
	"strings"
	"sync"
)

// Chain represents a Markov chain.
type Chain struct {
	Data      map[string][]string
	StateSize int
	mu        sync.RWMutex
}

// New creates a new chain with the given state size.
func New(stateSize int) *Chain {
	slog.Debug("New Chain object with state size", slog.Int("stateSize", stateSize))
	return &Chain{
		Data:      make(map[string][]string),
		StateSize: stateSize,
	}
}

// Add appends data to the currently loaded chain.
func (c *Chain) Add(message string) {
	words := strings.Fields(message)
	if len(words) < c.StateSize+1 {
		slog.Debug("Text was not added to chain", slog.String("text", message))
		return
	}

	c.mu.Lock() // Lock writes but not reads
	slog.Debug("Locked Chain during Add operation")
	defer c.mu.Unlock()

	// StateSize will determine how "smart" the bot is
	for i := 0; i <= len(words)-c.StateSize; i++ {
		key := strings.Join(words[i:i+c.StateSize], "")
		next := ""
		if i+c.StateSize < len(words) {
			next = words[i+c.StateSize]
		}
		slog.Debug("Adding text to chain", slog.String("key", key), slog.String("next", next))
		c.Data[key] = append(c.Data[key], next)
	}
	slog.Debug("Unlocking Chain after Add operation")
}

// Generate alllows the bot to create a response from the current chain
func (c *Chain) Generate() string {
	c.mu.RLock()
	slog.Debug("Read-locked chain for generation")
	defer c.mu.RUnlock()

	if len(c.Data) == 0 {
		return ""
	}

	const maxAttempts = 10
	for range maxAttempts {
		// Pick random starting key
		keys := make([]string, 0, len(c.Data))
		for k := range c.Data {
			keys = append(keys, k)
		}
		current := keys[rand.Intn(len(keys))]
		slog.Debug("Starting generation with key", "key", current)

		words := strings.Fields(current)
		const maxWords = 30

		for len(words) < maxWords {
			next, ok := c.Data[current]
			if !ok || len(next) == 0 {
				break
			}
			nextWord := next[rand.Intn(len(next))]
			if nextWord == "" {
				break
			}
			words = append(words, nextWord)

			// Slide window
			keyWords := append(strings.Fields(current)[1:], nextWord)
			current = strings.Join(keyWords, " ")
		}
	}
	slog.Debug("Unlocking read lock")
	return ""
}

// Seed the chain with text. Mostly useful to ensure that the chain doesn't begin empty.
func (c *Chain) Seed(text string) {
	lines := strings.SplitSeq(text, "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			c.Add(line)
		}
	}
}

// Size returns the length of the chain.
func (c *Chain) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.Data)
}
