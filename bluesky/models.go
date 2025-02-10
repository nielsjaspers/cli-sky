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

type BlueskyAuthResponse struct {
	AccessJwt  string `json:"accessJwt"`
	RefreshJwt string `json:"refreshJwt"`
	Handle     string `json:"handle"`
	Did        string `json:"did"`
}
