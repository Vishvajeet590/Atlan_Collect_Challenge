package entity

type JobRequest struct {
	FormId     int    `json:"form_Id"`
	OAuthCode  string `json:"OAuth_code"`
	PluginCode int    `json:"plugin_code"`
}
