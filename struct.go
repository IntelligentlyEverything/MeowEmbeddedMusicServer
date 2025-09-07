package main

// MusicItem represents a music item.
type MusicItem struct {
	Title        string `json:"title"`
	Artist       string `json:"artist"`
	AudioURL     string `json:"audio_url"`
	AudioFullURL string `json:"audio_full_url"`
	M3U8URL      string `json:"m3u8_url"`
	LyricURL     string `json:"lyric_url"`
	CoverURL     string `json:"cover_url"`
	Duration     int    `json:"duration"`
	FromCache    bool   `json:"from_cache"`
	IP           string `json:"ip"`
}
