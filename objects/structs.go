package objects

//WordDataSlice represents a list of words in JSON.
type WordDataSlice struct {
	List []WordData `json:"list"`
}

type SlackIncoming struct {
	Text         string `form:"text"`
	SlackUser    string `form:"user_name"`
	SlackChannel string `form:"channel_name"`
	SlackTeam    string `form:"team_id"`
	Token        string `form:"token"`
	Channel_id   string `form:"channel_id"`
	Command      string `form:"command"`
	Response_url string `form:"response_url"`
	Team_domain  string `form:"team_domain"`
	User_id      string `form:"user_id"`
}

type Response struct {
	Response   string `json:"response"`
	BotVersion string `json:"bot_version"`
}

//WordData represents the JSON struct sent by Urban Dictionary with the word.
type WordData struct {
	Author      string `json:"author"`
	CurrentVote string `json:"current_vote"`
	Defid       int    `json:"defid"`
	Definition  string `json:"definition"`
	Example     string `json:"example"`
	Permalink   string `json:"permalink"`
	ThumbsUp    int    `json:"thumbs_up"`
	ThumbsDown  int    `json:"thumbs_down"`
	Word        string `json:"word"`
}

type SlackResponse struct {
	Text         string `json:"text"`
	ResponseType string `json:"response_type"`
	BotVersion   string `json:"bot_version"`
}
