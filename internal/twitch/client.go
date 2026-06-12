package twitch

import (
	"log/slog"

	"github.com/aomarai/jam/internal/chain"
	"github.com/aomarai/jam/internal/config"
	twitchirc "github.com/gempir/go-twitch-irc/v4"
)

type Client struct {
	cfg    *config.Config
	chain  *chain.Chain
	client *twitchirc.Client
}

func NewClient(cfg *config.Config, c *chain.Chain) *Client {
	slog.Info("Creating new Twitch IRC Client", "username", cfg.TwitchBotUsername)
	tc := twitchirc.NewClient(cfg.TwitchBotUsername, cfg.TwitchOAuthToken)

	bot := &Client{
		cfg:    cfg,
		chain:  c,
		client: tc,
	}

	tc.OnPrivateMessage(func(msg twitchirc.PrivateMessage) {
		// Don't learn from bot's own messages
		if msg.User.Name == cfg.TwitchBotUsername {
			return
		}
		c.Add(msg.Message)
		slog.Debug("Added private message to chain", "text", msg.Message)
	})

	tc.OnConnect(func() {
		slog.Info("Connected to chat channel", "channel", cfg.TwitchTargetChannel)
	})

	tc.Join(cfg.TwitchTargetChannel)

	return bot
}

func (c *Client) Connect() error {
	return c.client.Connect()
}

func (c *Client) Post() {
	msg := c.chain.Generate()
	if msg == "" {
		slog.Info("Generate returned empty, skipping post")
		return
	}
	c.client.Say(c.TwitchTargetChannel(), msg)
	slog.Info("Posted message", "message", msg)
}

func (c *Client) TwitchTargetChannel() string {
	return c.cfg.TwitchTargetChannel
}
