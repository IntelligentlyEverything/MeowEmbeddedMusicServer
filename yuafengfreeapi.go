package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type YuafengAPIFreeResponse struct {
	Data struct {
		Song      string `json:"song"`
		Singer    string `json:"singer"`
		Cover     string `json:"cover"`
		AlbumName string `json:"album_name"`
		Music     string `json:"music"`
		Lyric     string `json:"lyric"`
	} `json:"data"`
}

// 枫雨API response handler.
func YuafengAPIResponseHandler(sources, song, singer string) MusicItem {
	fmt.Printf("[Info] Fetching music data from 枫林 free API for %s by %s\n", song, singer)
	var url string
	switch sources {
	case "kuwo":
		url = "https://api.yuafeng.cn/API/ly/kwmusic.php"
	case "netease":
		url = "https://api.yuafeng.cn/API/ly/wymusic.php"
	case "migu":
		url = "https://api.yuafeng.cn/API/ly/mgmusic.php"
	case "baidu":
		url = "https://api.yuafeng.cn/API/ly/bdmusic.php"
	default:
		return MusicItem{}
	}
	resp, err := http.Get(url + "?msg=" + song + "&n=1")
	if err != nil {
		fmt.Println("[Error] Error fetching the data from Yuafeng free API:", err)
		return MusicItem{}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[Error] Error reading the response body from Yuafeng free API:", err)
		return MusicItem{}
	}
	var response YuafengAPIFreeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("[Error] Error unmarshalling the data from Yuafeng free API:", err)
		return MusicItem{}
	}

	// Create directory
	dirName := fmt.Sprintf("./files/cache/music/%s-%s", response.Data.Singer, response.Data.Song)
	err = os.MkdirAll(dirName, 0755)
	if err != nil {
		fmt.Println("[Error] Error creating directory:", err)
		return MusicItem{}
	}

	if response.Data.Music == "" {
		fmt.Println("[Warning] Music URL is empty")
		return MusicItem{}
	}

	// Identify music file format
	musicExt, err := getMusicFileExtension(response.Data.Music)
	if err != nil {
		fmt.Println("[Error] Error identifying music file format:", err)
		return MusicItem{}
	}

	// Download music files
	err = downloadFile(filepath.Join(dirName, "music_full"+musicExt), response.Data.Music)
	if err != nil {
		fmt.Println("[Error] Error downloading music file:", err)
	}

	// Retrieve music file duration
	musicFilePath := filepath.Join(dirName, "music_full"+musicExt)
	duration := getMusicDuration(musicFilePath)

	// Download cover image
	ext := filepath.Ext(response.Data.Cover)
	err = downloadFile(filepath.Join(dirName, "cover"+ext), response.Data.Cover)
	if err != nil {
		fmt.Println("[Error] Error downloading cover image:", err)
	}

	// Check if the lyrics format is in link format
	lyricData := response.Data.Lyric
	if lyricData == "获取歌词失败" {
		// If it is "获取歌词失败", do nothing
		fmt.Println("[Warning] Lyric retrieval failed, skipping lyric file creation and download.")
	} else if !strings.HasPrefix(lyricData, "http://") && !strings.HasPrefix(lyricData, "https://") {
		// If it is not in link format, write the lyrics to the file line by line
		lines := strings.Split(lyricData, "\n")
		lyricFilePath := filepath.Join(dirName, "lyric.lrc")
		file, err := os.Create(lyricFilePath)
		if err != nil {
			fmt.Println("[Error] Error creating lyric file:", err)
			return MusicItem{}
		}
		defer file.Close()

		timeTagRegex := regexp.MustCompile(`^\[(\d+(?:\.\d+)?)\]`)
		for _, line := range lines {
			// Check if the line starts with a time tag
			match := timeTagRegex.FindStringSubmatch(line)
			if match != nil {
				// Convert the time tag to [mm:ss.ms] format
				timeInSeconds, _ := strconv.ParseFloat(match[1], 64)
				minutes := int(timeInSeconds / 60)
				seconds := int(timeInSeconds) % 60
				milliseconds := int((timeInSeconds-float64(seconds))*1000) / 100 % 100
				formattedTimeTag := fmt.Sprintf("[%02d:%02d.%02d]", minutes, seconds, milliseconds)
				line = timeTagRegex.ReplaceAllString(line, formattedTimeTag)
			}
			_, err := file.WriteString(line + "\r\n")
			if err != nil {
				fmt.Println("[Error] Error writing to lyric file:", err)
				return MusicItem{}
			}
		}
	} else {
		// If it is in link format, download the lyrics file
		err = downloadFile(filepath.Join(dirName, "lyric.lrc"), lyricData)
		if err != nil {
			fmt.Println("[Error] Error downloading lyric file:", err)
		}
	}

	// Compress and segment audio file
	err = compressAndSegmentAudio(filepath.Join(dirName, "music_full"+musicExt), dirName)
	if err != nil {
		fmt.Println("[Error] Error compressing and segmenting audio:", err)
	}

	// Create m3u8 playlist
	err = createM3U8Playlist(dirName)
	if err != nil {
		fmt.Println("[Error] Error creating m3u8 playlist:", err)
	}

	websiteURL := os.Getenv("WEBSITE_URL")
	return MusicItem{
		Title:        response.Data.Song,
		Artist:       response.Data.Singer,
		CoverURL:     websiteURL + "/cache/music/" + response.Data.Singer + "-" + response.Data.Song + "/cover" + ext,
		LyricURL:     websiteURL + "/cache/music/" + response.Data.Singer + "-" + response.Data.Song + "/lyric.lrc",
		AudioFullURL: websiteURL + "/cache/music/" + response.Data.Singer + "-" + response.Data.Song + "/music_full" + musicExt,
		AudioURL:     websiteURL + "/cache/music/" + response.Data.Singer + "-" + response.Data.Song + "/music.mp3",
		M3U8URL:      websiteURL + "/cache/music/" + response.Data.Singer + "-" + response.Data.Song + "/music.m3u8",
		Duration:     duration,
	}
}
