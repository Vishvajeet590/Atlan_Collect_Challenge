package entity

type Job struct {
	JobId         int    `json:"job_id"`
	JobStatus     string `json:"job_status"`
	JobStatusCode int    `json:"job_status_code"`
	PluginCode    int    `json:"job_plugin_code"`
}
