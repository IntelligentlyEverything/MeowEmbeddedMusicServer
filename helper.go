package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Helper function to compress and segment audio file
func compressAndSegmentAudio(inputFile, outputDir string) error {
	fmt.Printf("[Info] Compress and segment audio file %s\n", inputFile)
	// Compress music files
	outputFile := filepath.Join(outputDir, "music.mp3")
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-ac", "1", "-ab", "32k", "-ar", "24000", outputFile)
	err := cmd.Run()
	if err != nil {
		return err
	}

	// Split music files
	chunkDir := filepath.Join(outputDir, "chunk")
	err = os.MkdirAll(chunkDir, 0755)
	if err != nil {
		return err
	}

	// Using ffmpeg for segmentation
	segmentedFilePattern := filepath.Join(chunkDir, "%03d.mp3") // e.g. 001.mp3, 002.mp3, ...
	cmd = exec.Command("ffmpeg", "-i", outputFile, "-ac", "1", "-ab", "32k", "-ar", "16000", "-f", "segment", "-segment_time", "10", segmentedFilePattern)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Helper function to create M3U8 playlist file
func createM3U8Playlist(outputDir string) error {
	fmt.Printf("[Info] Create M3U8 playlist file for %s\n", outputDir)
	playlistFile := filepath.Join(outputDir, "music.m3u8")
	file, err := os.Create(playlistFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("#EXTM3U\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString("#EXT-X-VERSION:3\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString("#EXT-X-TARGETDURATION:10\n")
	if err != nil {
		return err
	}

	chunkDir := filepath.Join(outputDir, "chunk")
	files, err := os.ReadDir(chunkDir)
	if err != nil {
		return err
	}

	var chunkFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".mp3") {
			chunkFiles = append(chunkFiles, file.Name())
		}
	}

	// Sort by file name
	for i := 0; i < len(chunkFiles); i++ {
		for j := i + 1; j < len(chunkFiles); j++ {
			if chunkFiles[i] > chunkFiles[j] {
				chunkFiles[i], chunkFiles[j] = chunkFiles[j], chunkFiles[i]
			}
		}
	}

	for _, chunkFile := range chunkFiles {
		_, err = file.WriteString("#EXTINF:10.000\n")
		if err != nil {
			return err
		}
		url := fmt.Sprintf("%s/cache/music/%s/%s/%s\n", os.Getenv("EMBEDDED_WEBSITE_URL"), filepath.Base(outputDir), "chunk", chunkFile)
		_, err = file.WriteString(url)
	}

	return err
}

// Helper function to download files from URL
func downloadFile(filename string, url string) error {
	fmt.Printf("[Info] Download file %s from URL %s\n", filename, url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// Helper function to get duration of obtaining music files
func getMusicDuration(filePath string) int {
	fmt.Printf("[Info] Get duration of obtaining music file %s\n", filePath)
	// Use ffprobe to get audio duration
	output, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", filePath).Output()
	if err != nil {
		fmt.Println("[Error] Error getting audio duration:", err)
		return 0
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		fmt.Println("[Error] Error converting duration to float:", err)
		return 0
	}

	return int(duration)
}

// Helper function for identifying file formats
func getMusicFileExtension(url string) (string, error) {
	resp, err := http.Head(url)
	if err != nil {
		return "", err
	}
	// Get file format from Content-Type header
	contentType := resp.Header.Get("Content-Type")
	ext, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}
	// Identify file extension based on file format
	switch ext {
	case "audio/mpeg":
		return ".mp3", nil
	case "audio/flac":
		return ".flac", nil
	case "audio/x-flac":
		return ".flac", nil
	case "audio/wav":
		return ".wav", nil
	case "audio/aac":
		return ".aac", nil
	case "audio/ogg":
		return ".ogg", nil
	case "application/octet-stream":
		// Try to guess file format from URL or other information
		if strings.Contains(url, ".mp3") {
			return ".mp3", nil
		} else if strings.Contains(url, ".flac") {
			return ".flac", nil
		} else if strings.Contains(url, ".wav") {
			return ".wav", nil
		} else if strings.Contains(url, ".aac") {
			return ".aac", nil
		} else if strings.Contains(url, ".ogg") {
			return ".ogg", nil
		} else {
			return "", fmt.Errorf("unknown file format from Content-Type and URL: %s", contentType)
		}
	default:
		return "", fmt.Errorf("unknown file format: %s", ext)
	}
}

// Helper function to obtain music data from local folder
func getLocalMusicItem(song, singer string) MusicItem {
	musicDir := "./files/music"
	fmt.Println("[Info] Reading local folder music.")
	files, err := os.ReadDir(musicDir)
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
						musicItem.AudioURL = "/music/" + url.QueryEscape(file.Name()) + "/music.mp3"
						musicItem.Duration = getMusicDuration(musicFilePath)
					}

					for _, audioFormat := range []string{"music_full.mp3", "music_full.flac", "music_full.wav", "music_full.aac", "music_full.ogg"} {
						audioFilePath := filepath.Join(dirPath, audioFormat)
						if _, err := os.Stat(audioFilePath); err == nil {
							musicItem.AudioFullURL = "/music/" + url.QueryEscape(file.Name()) + "/" + audioFormat
							break
						}
					}

					m3u8FilePath := filepath.Join(dirPath, "music.m3u8")
					if _, err := os.Stat(m3u8FilePath); err == nil {
						musicItem.M3U8URL = "/music/" + url.QueryEscape(file.Name()) + "/music.m3u8"
					}

					lyricFilePath := filepath.Join(dirPath, "lyric.lrc")
					if _, err := os.Stat(lyricFilePath); err == nil {
						musicItem.LyricURL = "/music/" + url.QueryEscape(file.Name()) + "/lyric.lrc"
					}

					coverJpgFilePath := filepath.Join(dirPath, "cover.jpg")
					if _, err := os.Stat(coverJpgFilePath); err == nil {
						musicItem.CoverURL = "/music/" + url.QueryEscape(file.Name()) + "/cover.jpg"
					} else {
						coverPngFilePath := filepath.Join(dirPath, "cover.png")
						if _, err := os.Stat(coverPngFilePath); err == nil {
							musicItem.CoverURL = "/music/" + url.QueryEscape(file.Name()) + "/cover.png"
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
						musicItem.AudioURL = "/music/" + url.QueryEscape(file.Name()) + "/music.mp3"
						musicItem.Duration = getMusicDuration(musicFilePath)
					}

					for _, audioFormat := range []string{"music_full.mp3", "music_full.flac", "music_full.wav", "music_full.aac", "music_full.ogg"} {
						audioFilePath := filepath.Join(dirPath, audioFormat)
						if _, err := os.Stat(audioFilePath); err == nil {
							musicItem.AudioFullURL = "/music/" + url.QueryEscape(file.Name()) + "/" + audioFormat
							break
						}
					}

					m3u8FilePath := filepath.Join(dirPath, "music.m3u8")
					if _, err := os.Stat(m3u8FilePath); err == nil {
						musicItem.M3U8URL = "/music/" + url.QueryEscape(file.Name()) + "/music.m3u8"
					}

					lyricFilePath := filepath.Join(dirPath, "lyric.lrc")
					if _, err := os.Stat(lyricFilePath); err == nil {
						musicItem.LyricURL = "/music/" + url.QueryEscape(file.Name()) + "/lyric.lrc"
					}

					coverJpgFilePath := filepath.Join(dirPath, "cover.jpg")
					if _, err := os.Stat(coverJpgFilePath); err == nil {
						musicItem.CoverURL = "/music/" + url.QueryEscape(file.Name()) + "/cover.jpg"
					} else {
						coverPngFilePath := filepath.Join(dirPath, "cover.png")
						if _, err := os.Stat(coverPngFilePath); err == nil {
							musicItem.CoverURL = "/music/" + url.QueryEscape(file.Name()) + "/cover.png"
						}
					}

					return musicItem
				}
			}
		}
	}

	return MusicItem{} // If no matching folder is found, return an empty MusicItem
}

// Helper function to obtain IP address of the client
func IPhandler(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip, nil
	}
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0]), nil
	}
	ip = r.RemoteAddr
	if ip != "" {
		return strings.Split(ip, ":")[0], nil
	}

	return "", fmt.Errorf("unable to obtain IP address information")
}

// Helper function to read music sources from sources.json file
func readSources() []MusicItem {
	data, err := os.ReadFile("./sources.json")
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

// Helper function to request and cache music from API sources
func requestAndCacheMusic(song, singer string) MusicItem {
	fmt.Printf("[Info] Requesting and caching music for %s", song)
	// Create cache directory if it doesn't exist
	err := os.MkdirAll("./cache", 0755)
	if err != nil {
		fmt.Println("[Error] Error creating cache directory:", err)
		return MusicItem{}
	}

	// Get API_SOURCES and any subsequent environment variables (e.g. API_SOURCES_1, API_SOURCES_2, etc.)
	var sources []string
	for i := 0; ; i++ {
		var key string
		if i == 0 {
			key = "API_SOURCES"
		} else {
			key = "API_SOURCES_" + strconv.Itoa(i)
		}
		source := os.Getenv(key)
		if source == "" {
			break
		}
		sources = append(sources, source)
	}

	// Request and cache music from each source in turn
	var musicItem MusicItem
	for _, source := range sources {
		fmt.Printf("[Info] Requesting music from source: %s\n", source)
		musicItem = YuafengAPIResponseHandler(strings.TrimSpace(source), song, singer)
		if musicItem.Title != "" {
			// If music item is valid, stop searching for sources
			break
		}
	}

	// If no valid music item was found, return an empty MusicItem
	if musicItem.Title == "" {
		fmt.Printf("[Warning] No valid music item retrieved.\n")
		return MusicItem{}
	}

	// Create cache file path based on artist and title
	cacheFile := fmt.Sprintf("./cache/%s-%s.json", musicItem.Artist, musicItem.Title)

	// Write cache data to cache file
	cacheData, err := json.MarshalIndent(musicItem, "", "  ")
	if err != nil {
		fmt.Println("[Error] Error marshalling cache data:", err)
		return MusicItem{}
	}
	err = os.WriteFile(cacheFile, cacheData, 0644)
	if err != nil {
		fmt.Println("[Error] Error writing cache file:", err)
		return MusicItem{}
	}

	fmt.Println("[Info] Music request and caching completed successfully.")
	return musicItem
}

// Helper function to read music data from cache file
func readFromCache(filePath string) (MusicItem, bool) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("[Error] Failed to read cache file:", err)
		return MusicItem{}, false
	}

	var musicItem MusicItem
	err = json.Unmarshal(data, &musicItem)
	if err != nil {
		fmt.Println("[Error] Failed to parse cache file:", err)
		return MusicItem{}, false
	}

	return musicItem, true
}
