package bluesky

type BlueskyPost struct {
	Type      string `json:"$type"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

type BlueskyAuth struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}
