package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// APIHandler handles API requests.
func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "MeowMusicEmbeddedServer")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	queryParams := r.URL.Query()
	fmt.Printf("[Web Access] Handling request for %s?%s\n", r.URL.Path, queryParams.Encode())
	song := queryParams.Get("song")
	singer := queryParams.Get("singer")

	ip, err := IPhandler(r)
	if err != nil {
		ip = "0.0.0.0"
	}

	if song == "" {
		musicItem := MusicItem{
			FromCache: false,
			IP:        ip,
		}
		json.NewEncoder(w).Encode(musicItem)
		return
	}

	// Attempt to retrieve music items from sources.json
	sources := readSources()

	var musicItem MusicItem
	var found bool = false

	// Build request scheme
	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	for _, source := range sources {
		if source.Title == song {
			if singer == "" || source.Artist == singer {
				// Determine the protocol for each URL and build accordingly
				var audioURL, audioFullURL, m3u8URL, lyricURL, coverURL string
				if strings.HasPrefix(source.AudioURL, "http://") {
					audioURL = scheme + "://" + r.Host + "/url/http/" + url.QueryEscape(strings.TrimPrefix(source.AudioURL, "http://"))
				} else if strings.HasPrefix(source.AudioURL, "https://") {
					audioURL = scheme + "://" + r.Host + "/url/https/" + url.QueryEscape(strings.TrimPrefix(source.AudioURL, "https://"))
				} else {
					audioURL = scheme + "://" + r.Host + "/" + url.QueryEscape(source.AudioURL)
				}
				if strings.HasPrefix(source.AudioFullURL, "http://") {
					audioFullURL = scheme + "://" + r.Host + "/url/http/" + url.QueryEscape(strings.TrimPrefix(source.AudioFullURL, "http://"))
				} else if strings.HasPrefix(source.AudioFullURL, "https://") {
					audioFullURL = scheme + "://" + r.Host + "/url/https/" + url.QueryEscape(strings.TrimPrefix(source.AudioFullURL, "https://"))
				} else {
					audioFullURL = scheme + "://" + r.Host + "/" + url.QueryEscape(source.AudioFullURL)
				}
				if strings.HasPrefix(source.M3U8URL, "http://") {
					m3u8URL = scheme + "://" + r.Host + "/url/http/" + url.QueryEscape(strings.TrimPrefix(source.M3U8URL, "http://"))
				} else if strings.HasPrefix(source.M3U8URL, "https://") {
					m3u8URL = scheme + "://" + r.Host + "/url/https/" + url.QueryEscape(strings.TrimPrefix(source.M3U8URL, "https://"))
				} else {
					m3u8URL = scheme + "://" + r.Host + "/" + url.QueryEscape(source.M3U8URL)
				}
				if strings.HasPrefix(source.LyricURL, "http://") {
					lyricURL = scheme + "://" + r.Host + "/url/http/" + url.QueryEscape(strings.TrimPrefix(source.LyricURL, "http://"))
				} else if strings.HasPrefix(source.LyricURL, "https://") {
					lyricURL = scheme + "://" + r.Host + "/url/https/" + url.QueryEscape(strings.TrimPrefix(source.LyricURL, "https://"))
				} else {
					lyricURL = scheme + "://" + r.Host + "/" + url.QueryEscape(source.LyricURL)
				}
				if strings.HasPrefix(source.CoverURL, "http://") {
					coverURL = scheme + "://" + r.Host + "/url/http/" + url.QueryEscape(strings.TrimPrefix(source.CoverURL, "http://"))
				} else if strings.HasPrefix(source.CoverURL, "https://") {
					coverURL = scheme + "://" + r.Host + "/url/https/" + url.QueryEscape(strings.TrimPrefix(source.CoverURL, "https://"))
				} else {
					coverURL = scheme + "://" + r.Host + "/" + url.QueryEscape(source.CoverURL)
				}
				musicItem = MusicItem{
					Title:        source.Title,
					Artist:       source.Artist,
					AudioURL:     audioURL,
					AudioFullURL: audioFullURL,
					M3U8URL:      m3u8URL,
					LyricURL:     lyricURL,
					CoverURL:     coverURL,
					Duration:     source.Duration,
					FromCache:    false,
				}
				found = true
				break
			}
		}
	}

	// If not found in sources.json, attempt to retrieve from local folder
	if !found {
		musicItem = getLocalMusicItem(song, singer)
		musicItem.FromCache = false
		if musicItem.Title != "" {
			if musicItem.AudioURL != "" {
				musicItem.AudioURL = scheme + "://" + r.Host + musicItem.AudioURL
			}
			if musicItem.AudioFullURL != "" {
				musicItem.AudioFullURL = scheme + "://" + r.Host + musicItem.AudioFullURL
			}
			if musicItem.M3U8URL != "" {
				musicItem.M3U8URL = scheme + "://" + r.Host + musicItem.M3U8URL
			}
			if musicItem.LyricURL != "" {
				musicItem.LyricURL = scheme + "://" + r.Host + musicItem.LyricURL
			}
			if musicItem.CoverURL != "" {
				musicItem.CoverURL = scheme + "://" + r.Host + musicItem.CoverURL
			}
			found = true
		}
	}

	// If still not found, attempt to retrieve from cache file
	if !found {
		fmt.Println("[Info] Reading music from cache.")
		// Fuzzy matching for singer and song
		files, err := filepath.Glob("./cache/*.json")
		if err != nil {
			fmt.Println("[Error] Error reading cache directory:", err)
			return
		}
		for _, file := range files {
			if strings.Contains(filepath.Base(file), song) && (singer == "" || strings.Contains(filepath.Base(file), singer)) {
				musicItem, found = readFromCache(file)
				if found {
					if musicItem.AudioURL != "" {
						musicItem.AudioURL = scheme + "://" + r.Host + musicItem.AudioURL
					}
					if musicItem.AudioFullURL != "" {
						musicItem.AudioFullURL = scheme + "://" + r.Host + musicItem.AudioFullURL
					}
					if musicItem.M3U8URL != "" {
						musicItem.M3U8URL = scheme + "://" + r.Host + musicItem.M3U8URL
					}
					if musicItem.LyricURL != "" {
						musicItem.LyricURL = scheme + "://" + r.Host + musicItem.LyricURL
					}
					if musicItem.CoverURL != "" {
						musicItem.CoverURL = scheme + "://" + r.Host + musicItem.CoverURL
					}
					musicItem.FromCache = true
					break
				}
			}
		}
	}

	// If still not found, request and cache the music item in a separate goroutine
	if !found {
		fmt.Println("[Info] Updating music item cache from API request.")
		musicItem = requestAndCacheMusic(song, singer)
		fmt.Println("[Info] Music item cache updated.")
		musicItem.FromCache = false
		musicItem.AudioURL = scheme + "://" + r.Host + musicItem.AudioURL
		musicItem.AudioFullURL = scheme + "://" + r.Host + musicItem.AudioFullURL
		musicItem.M3U8URL = scheme + "://" + r.Host + musicItem.M3U8URL
		musicItem.LyricURL = scheme + "://" + r.Host + musicItem.LyricURL
		musicItem.CoverURL = scheme + "://" + r.Host + musicItem.CoverURL
		found = true
	}

	// If still not found, return an empty MusicItem
	if !found {
		musicItem = MusicItem{
			FromCache: false,
			IP:        ip,
		}
	} else {
		musicItem.IP = ip
	}

	json.NewEncoder(w).Encode(musicItem)
}
