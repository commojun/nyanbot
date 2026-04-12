package config

type Config struct {
	ServerPort         int    `default:"8999" env:"NYAN_SERVER_PORT" help:"Server port"`
	ChannelSecret      string `required:"" env:"NYAN_CHANNEL_SECRET" help:"LINE channel secret"`
	ChannelAccessToken string `required:"" env:"NYAN_ACCESS_TOKEN" help:"LINE access token"`
	DefaultRoomID      string `env:"NYAN_DEFAULT_ROOM_ID" help:"Default LINE room ID"`
	MessageToken       string `env:"NYAN_MESSAGE_TOKEN" help:"Message API auth token"`
	GoogleClientEmail  string `required:"" env:"NYAN_GOOGLE_CLIENT_EMAIL" help:"Google client email"`
	GooglePrivateKey   string `required:"" env:"NYAN_GOOGLE_PRIVATE_KEY" help:"Google private key"`
	GooglePrivateKeyID string `env:"NYAN_GOOGLE_PRIVATE_KEY_ID" help:"Google private key ID"`
	GoogleTokenURL     string `env:"NYAN_GOOGLE_TOKEN_URL" help:"Google token URL"`
	SheetID            string `required:"" env:"NYAN_SHEET_ID" help:"Google Sheet ID"`
}
