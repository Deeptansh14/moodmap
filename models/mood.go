package models

type Mood struct {
	Username  string  `json:"username"`
	Emoji     string  `json:"emoji"`
	Message   string  `json:"message"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Timestamp string  `json:"timestamp"`
	ClusterID int     `json:"cluster_id"`
}
