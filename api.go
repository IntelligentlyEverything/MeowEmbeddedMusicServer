package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

	// Attempt to retrieve music items from sources.json
	sources := readSources()

	var musicItem MusicItem
	var found bool = false

	for _, source := range sources {
		if source.Title == song {
			if singer == "" || source.Artist == singer {
				musicItem = MusicItem{
					Title:     source.Title,
					Artist:    source.Artist,
					AudioURL:  source.AudioURL,
					M3U8URL:   source.M3U8URL,
					LyricURL:  source.LyricURL,
					CoverURL:  source.CoverURL,
					Duration:  source.Duration,
					FromCache: false,
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
					musicItem.FromCache = true
					break
				}
			}
		}
	}

	// If still not found, request and cache the music item in a separate goroutine
	if !found {
		fmt.Println("[Info] Updating music item cache from API request.")
		go func() {
			requestAndCacheMusic(song, singer)
			fmt.Println("[Info] Music item cache updated.")
		}()
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

// Read sources.json file and return a list of SourceItem.
func readSources() []MusicItem {
	data, err := ioutil.ReadFile("./sources.json")
	fmt.Println("[Info] Reading local sources.json")
	if err != nil {
		fmt.Println("[Error] Failed to read sources.json:", err)
		return nil
	}

	var sources []MusicItem
	err = json.Unmarshal(data, &sources)
	if err != nil {
		fmt.Println("[Error] Failed to parse sources.json:", err)
		return nil
	}

	return sources
}

// Retrieve music items from local folder
func getLocalMusicItem(song, singer string) MusicItem {
	musicDir := "./files/music"
	fmt.Println("[Info] Reading local folder music.")
	files, err := ioutil.ReadDir(musicDir)
	if err != nil {
		fmt.Println("[Error] Failed to read local music directory:", err)
		return MusicItem{}
	}

	for _, file := range files {
		if file.IsDir() {
			if singer == "" {
				if strings.Contains(file.Name(), song) {
					dirPath := filepath.Join(musicDir, file.Name())
					// Extract artist and title from the directory name
					parts := strings.SplitN(file.Name(), "-", 2)
					if len(parts) != 2 {
						continue // Skip if the directory name doesn't contain exactly one "-"
					}
					artist := parts[0]
					title := parts[1]
					musicItem := MusicItem{
						Title:        title,
						Artist:       artist,
						AudioURL:     "",
						AudioFullURL: "",
						M3U8URL:      "",
						LyricURL:     "",
						CoverURL:     "",
						Duration:     0,
					}

					musicFilePath := filepath.Join(dirPath, "music.mp3")
					if _, err := os.Stat(musicFilePath); err == nil {
						musicItem.AudioURL = os.Getenv("WEBSITE_URL") + "/music/" + file.Name() + "/music.mp3"
						musicItem.Duration = getMusicDuration(musicFilePath)
					}

					for _, audioFormat := range []string{"music_full.mp3", "music_full.flac", "music_full.wav", "music_full.aac", "music_full.ogg"} {
						audioFilePath := filepath.Join(dirPath, audioFormat)
						if _, err := os.Stat(audioFilePath); err == nil {
							musicItem.AudioFullURL = os.Getenv("WEBSITE_URL") + "/music/" + file.Name() + "/" + audioFormat
							break
						}
					}

					m3u8FilePath := filepath.Join(dirPath, "music.m3u8")
					if _, err := os.Stat(m3u8FilePath); err == nil {
						musicItem.M3U8URL = os.Getenv("WEBSITE_URL") + "/music/" + file.Name() + "/music.m3u8"
					}

					lyricFilePath := filepath.Join(dirPath, "lyric.lrc")
					if _, err := os.Stat(lyricFilePath); err == nil {
						musicItem.LyricURL = os.Getenv("WEBSITE_URL") + "/music/" + file.Name() + "/lyric.lrc"
					}

					coverJpgFilePath := filepath.Join(dirPath, "cover.jpg")
					if _, err := os.Stat(coverJpgFilePath); err == nil {
						musicItem.CoverURL = os.Getenv("WEBSITE_URL") + "/music/" + file.Name() + "/cover.jpg"
					} else {
						coverPngFilePath := filepath.Join(dirPath, "cover.png")
						if _, err := os.Stat(coverPngFilePath); err == nil {
							musicItem.CoverURL = os.Getenv("WEBSITE_URL") + "/music/" + file.Name() + "/cover.png"
						}
					}

					return musicItem
				}
			} else {
				if strings.Contains(file.Name(), singer) && strings.Contains(file.Name(), song) {
					dirPath := filepath.Join(musicDir, file.Name())
					// Extract artist and title from the directory name
					parts := strings.SplitN(file.Name(), "-", 2)
					if len(parts) != 2 {
						continue // Skip if the directory name doesn't contain exactly one "-"
					}
					artist := parts[0]
					title := parts[1]
					musicItem := MusicItem{
						Title:        title,
						Artist:       artist,
						AudioURL:     "",
						AudioFullURL: "",
						M3U8URL:      "",
						LyricURL:     "",
						CoverURL:     "",
						Duration:     0,
					}

					musicFilePath := filepath.Join(dirPath, "music.mp3")
					if _, err := os.Stat(musicFilePath); err == nil {
						musicItem.AudioURL = os.Getenv("WEBSITE_URL") + "/music/" + file.Name() + "/music.mp3"
						musicItem.Duration = getMusicDuration(musicFilePath)
					}

					for _, audioFormat := range []string{"music_full.mp3", "music_full.flac", "music_full.wav", "music_full.aac", "music_full.ogg"} {
						audioFilePath := filepath.Join(dirPath, audioFormat)
						if _, err := os.Stat(audioFilePath); err == nil {
							musicItem.AudioFullURL = os.Getenv("WEBSITE_URL") + "/music/" + file.Name() + "/" + audioFormat
							break
						}
					}

					m3u8FilePath := filepath.Join(dirPath, "music.m3u8")
					if _, err := os.Stat(m3u8FilePath); err == nil {
						musicItem.M3U8URL = os.Getenv("WEBSITE_URL") + "/music/" + file.Name() + "/music.m3u8"
					}

					lyricFilePath := filepath.Join(dirPath, "lyric.lrc")
					if _, err := os.Stat(lyricFilePath); err == nil {
						musicItem.LyricURL = os.Getenv("WEBSITE_URL") + "/music/" + file.Name() + "/lyric.lrc"
					}

					coverJpgFilePath := filepath.Join(dirPath, "cover.jpg")
					if _, err := os.Stat(coverJpgFilePath); err == nil {
						musicItem.CoverURL = os.Getenv("WEBSITE_URL") + "/music/" + file.Name() + "/cover.jpg"
					} else {
						coverPngFilePath := filepath.Join(dirPath, "cover.png")
						if _, err := os.Stat(coverPngFilePath); err == nil {
							musicItem.CoverURL = os.Getenv("WEBSITE_URL") + "/music/" + file.Name() + "/cover.png"
						}
					}

					return musicItem
				}
			}
		}
	}

	return MusicItem{} // If no matching folder is found, return an empty MusicItem
}
