package discord

type DiscordCreateMessageBody struct {
	Content string             `json:"content"`
	Embeds  []ItemEmbedDiscord `json:"embeds"`
}

type ItemEmbedDiscord struct {
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Image       ImageItemEmbedDiscord `json:"image"`
}

type ImageItemEmbedDiscord struct {
	Url string `json:"url"`
}
