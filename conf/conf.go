package conf

type Config struct {
	Version, TgBotToken, ApiLogin, ApiPassword, ApiBasePath, LogPath string
}

var DefaultConfig Config = Config{
	"0.0.1",
	"TELEGRAM BOT TOKEN",
	"USERNAME",
	"PASSWORD",
	"BASE PATH (include last slash)",
	"ERROR LOG PATH",
}
