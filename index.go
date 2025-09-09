package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "MeowMusicEmbeddedServer")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Printf("[Web Access] Handling request for %s\n", r.URL.Path)
	if r.URL.Path != "/" {
		fileHandler(w, r)
		return
	}
	// Serve index.html in theme directory
	indexPath := filepath.Join("theme", "index.html")

	// Check if index.html exists in theme directory
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		defaultIndexPage(w)
	} else if err != nil {
		defaultIndexPage(w)
	} else {
		http.ServeFile(w, r, indexPath)
	}
}

func defaultIndexPage(w http.ResponseWriter) {
	websiteVersion := "0.0.1 alpha 20250909"
	websiteNameCN := os.Getenv("WEBSITE_NAME_CN")
	if websiteNameCN == "" {
		websiteNameCN = "🎵 音乐搜索"
	}
	websiteNameEN := os.Getenv("WEBSITE_NAME_EN")
	if websiteNameEN == "" {
		websiteNameEN = "🎵 Music Search"
	}
	websiteTitleCN := os.Getenv("WEBSITE_TITLE_CN")
	if websiteTitleCN == "" {
		websiteTitleCN = "为嵌入式设备设计的音乐搜索服务器"
	}
	websiteTitleEN := os.Getenv("WEBSITE_TITLE_EN")
	if websiteTitleEN == "" {
		websiteTitleEN = "Music Search Server for Embedded Devices"
	}
	websiteDescCN := os.Getenv("WEBSITE_DESC_CN")
	if websiteDescCN == "" {
		websiteDescCN = "搜索并播放您喜爱的音乐"
	}
	websiteDescEN := os.Getenv("WEBSITE_DESC_EN")
	if websiteDescEN == "" {
		websiteDescEN = "Search and play your favorite music"
	}
	websiteKeywordsCN := os.Getenv("WEBSITE_KEYWORDS_CN")
	if websiteKeywordsCN == "" {
		websiteKeywordsCN = "音乐, 搜索, 嵌入式"
	}
	websiteKeywordsEN := os.Getenv("WEBSITE_KEYWORDS_EN")
	if websiteKeywordsEN == "" {
		websiteKeywordsEN = "music, search, embedded"
	}
	websiteFavicon := os.Getenv("WEBSITE_FAVICON")
	if websiteFavicon == "" {
		websiteFavicon = "/favicon.ico"
	}
	websiteBackground := os.Getenv("WEBSITE_BACKGROUND")
	if websiteBackground == "" {
		websiteBackground = "/background.webp"
	}
	fontawesomeCDN := os.Getenv("FONTAWESOME_CDN")
	if fontawesomeCDN == "" {
		fontawesomeCDN = "https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css"
	}

	// Build HTML
	fmt.Fprintf(w, "<!DOCTYPE html><html>")
	fmt.Fprintf(w, "<head>")
	fmt.Fprintf(w, "<meta charset=\"UTF-8\">")
	fmt.Fprintf(w, "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">")
	fmt.Fprintf(w, "<meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\">")
	fmt.Fprintf(w, "<link rel=\"icon\" href=\"%s\">", websiteFavicon)
	fmt.Fprintf(w, "<link rel=\"stylesheet\" href=\"%s\">", fontawesomeCDN)
	fmt.Fprintf(w, "<title></title><style>")
	// HTML style
	fmt.Fprintf(w, "body {background-image: url('%s');background-size: cover;background-repeat: no-repeat;background-attachment: fixed;display: flex;justify-content: center;align-items: center;margin: 60px 0;}", websiteBackground)
	fmt.Fprintf(w, ".container {background: rgba(255, 255, 255, 0.4);width: 65%%;border-radius: 20px;box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);backdrop-filter: blur(10px);display: flex;flex-direction: column;}")
	fmt.Fprintf(w, ".title {font-size: 36px;font-weight: bold;margin: 25px auto 0px auto;text-align: center;}.description {font-size: 1.1rem;color: #4f596b;margin: 10px auto;text-align: center;}")
	fmt.Fprintf(w, ".search-form {display: flex;justify-content: center;align-items: center;width: 100%%;margin-bottom: 20px;}.songContainer,.singerContainer,.searchContainer {display: flex;align-items: center;margin: 0 10px;}")
	fmt.Fprintf(w, ".songInput {padding: 10px;border: 2px solid #ccc;border-radius: 20px;height: 45px;margin-left: -6%%;margin-right: 10px;font-size: 1.1rem;width: 110%%;background-color: rgba(255, 255, 255, 0.4);}")
	fmt.Fprintf(w, ".artistInput {padding: 10px;border: 2px solid #ccc;border-radius: 20px;height: 45px;margin-left: 7%%;margin-right: 10px;font-size: 1.1rem;width: 80%%;background-color: rgba(255, 255, 255, 0.4);}")
	fmt.Fprintf(w, ".searchBtn {padding: 10px 20px;border: none;background-image: linear-gradient(to right, pink, deeppink);color: white;margin-left: -20%%;font-size: 1.1rem;border-radius: 20px;width: 128%%;height: 60px;cursor: pointer;transition: all 0.3s ease;}")
	fmt.Fprintf(w, "@media (max-width: 768px) {.search-form {flex-direction: column;align-items: flex-start;text-align: center;}.songContainer,.singerContainer,.searchContainer {display: block;margin: 4px 12%% 0 auto;width: 80%%;}.songInput {margin: 0;width: 100%%;}.artistInput {margin: 0;width: 100%%;}.searchBtn {margin: 0;width: 106%%;height: 40px;}}")
	fmt.Fprintf(w, ".songInput:hover,.artistInput:hover,.songInput:focus,.artistInput:focus {outline: none;border: 2px solid deeppink;}.song-item:hover,.searchBtn:hover {box-shadow: 0 4px 8px rgba(255, 182, 193, 0.7);transform: translateY(-5px);}")
	fmt.Fprintf(w, ".getError,.no-enter,.no-result {width: 80%%;margin: 4px auto;padding: 20px;background-color: rgba(255, 0, 38, 0.4);text-align: center;border: 1px solid rgb(255, 75, 75);border-radius: 15px;color: rgb(205, 0, 0);}")
	fmt.Fprintf(w, ".loading {width: 80%%;margin: 4px auto;padding: 20px;text-align: center;color: deeppink;font-size: 45px;animation: spin 2s linear infinite;}@keyframes spin {from {transform: rotate(0deg);}to {transform: rotate(360deg);}}")
	fmt.Fprintf(w, ".result {width: 85%%;margin: 4px auto;}.result-title {font-size: 24px;font-weight: bold;}.song-item {background-color: rgba(255, 255, 255, 0.4);border: 2px solid deeppink;border-radius: 15px;transition: all 0.3s ease;padding: 10px;}")
	fmt.Fprintf(w, ".song-title-container {display: flex;align-items: center;}.song-name {font-size: 18px;font-weight: bold;}.cache {width: 45px;background-color: deepskyblue;color: #000;font-size: 14px;text-align: center;border-radius: 15px;}")
	fmt.Fprintf(w, ".singer-name-icon,.lyric-icon {font-size: 18px;color: deeppink;}.singer-name,.lyric {font-size: 16px;color: #4f596b;}.playBtn,.pauseBtn {border: none;background-image: linear-gradient(to right, skyblue, deepskyblue);border-radius: 5px;padding: 5px 10px;font-size: 15px;transition: all 0.3s ease;}.playBtn:hover,.pauseBtn:hover {box-shadow: 0 4px 8px rgba(182, 232, 255, 0.7);transform: translateY(-5px);}")
	fmt.Fprintf(w, ".audio-player-container {display: flex;align-items: center;}.audio {display: none;}.progress-bar {width: 70%%;margin: 4px auto;padding: 8px;background-color: rgba(255, 255, 255, 0.4);border: 1px solid deeppink;border-radius: 5px;display: flex;justify-content: space-between;align-items: center;}.progress {width: 0;height: 10px;background-color: deeppink;}.time {margin-left: auto;}")
	fmt.Fprintf(w, ".stream_pcm {width: 80%%;margin: 4px auto;padding: 20px;background-color: rgba(135, 206, 235, 0.4);border: 1px solid skyblue;border-radius: 15px;}.stream_pcm_title {color: rgb(0, 100, 100);font-size: 16px;font-weight: bold;}.stream_pcm_content {margin-top: 10px;font-size: 14px;color: #555;}")
	fmt.Fprintf(w, ".stream_pcm_type_title,.stream_pcm_content_num_title,.stream_pcm_content_time_title,.stream_pcm_response_title {font-weight: bold;}.stream_pcm_response_value {width: 100%%;background-color: rgba(255, 255, 255, 0.4);display: block;white-space: pre-wrap;overflow: auto;height: 200px;border-radius: 6px;padding: 10px;}")
	fmt.Fprintf(w, ".info {width: 80%%;margin: 4px auto;padding: 20px;text-align: center;color: #4f596b;}.info strong {font-weight: bolder;color: #000;}.showStreamPcmBtnContainer,.hideStreamPcmBtnContainer {margin: 0 auto;width: 80%%;display: flex;justify-content: center;}")
	fmt.Fprintf(w, ".showStreamPcmBtn {border: 1px solid deepskyblue;color: deepskyblue;}.showStreamPcmBtn:hover {background-color: deepskyblue;color: #000;}.hideStreamPcmBtn {border: 1px solid deeppink;color: deeppink;}.hideStreamPcmBtn:hover {background-color: deeppink;color: #000;}.showStreamPcmBtn,.hideStreamPcmBtn {background: none;padding: 2px 6px;}")
	fmt.Fprintf(w, ".footer {text-align: center;margin: 10px auto;justify-content: center;align-items: center;width: 80%%;border-top: 1px solid #ccc;}.language-select {background-color: rgba(255, 255, 255, 0.4);border: 1px solid #ccc;text-align: center;width: 120px;height: 40px;border-radius: 10px;margin: 10px auto;}")
	fmt.Fprintf(w, ".language-select:focus,.language-select:hover {outline: none;border: 1px solid deeppink;}.copyright {font-size: 14px;color: #4f596b;}")
	fmt.Fprintf(w, "</style></head>")
	// Build body
	fmt.Fprintf(w, "<body><div class=\"container\"><div id=\"title\" class=\"title\"></div><div id=\"description\" class=\"description\"></div>")
	fmt.Fprintf(w, "<div class=\"search-form\"><div class=\"songContainer\"><div class=\"song\"><input type=\"text\" id=\"songInput\" class=\"songInput\" autocomplete=\"off\"></div></div>")
	fmt.Fprintf(w, "<div class=\"singerContainer\"><div class=\"singer\"><input type=\"text\" id=\"artistInput\" class=\"artistInput\" autocomplete=\"off\"></div></div><div class=\"searchContainer\"><div class=\"search\"><button type=\"button\" id=\"searchBtn\" class=\"searchBtn\"></button></div></div></div>")
	fmt.Fprintf(w, "<div class=\"getError\" id=\"getError\"></div><div class=\"no-enter\" id=\"noEnter\"></div><div class=\"no-result\" id=\"noResult\"></div><div class=\"loading\" id=\"loading\"><i class=\"fa fa-circle-o-notch\"></i></div>")
	fmt.Fprintf(w, "<div class=\"result\" id=\"result\"><div class=\"result-title\" id=\"resultTitle\"></div><div class=\"result-list\"><div class=\"song-item\"><div class=\"song-title-container\"><div class=\"song-name\" id=\"songName\"></div><div class=\"cache\" id=\"cache\"></div></div><div class=\"singer-name\"><span class=\"singer-name-icon\" id=\"singerNameIcon\"><i class=\"fa fa-user-o\"></i></span><span class=\"singer-name-value\" id=\"singerName\"></span></div><div class=\"lyric\"><span class=\"lyric-icon\" id=\"lyricIcon\"><i class=\"fa fa-file-text-o\"></i></span><span class=\"lyric-value\" id=\"noLyric\"></span><span class=\"lyric-value\" id=\"lyric\"></span></div><div class=\"audio-player-container\"><button type=\"button\" class=\"playBtn\" id=\"playBtn\"></button><button type=\"button\" class=\"pauseBtn\" id=\"pauseBtn\"></button><audio class=\"audio\" id=\"audio\"></audio><div class=\"progress-bar\"><div class=\"progress\" id=\"progress\"></div><div class=\"time\" id=\"time\"></div></div></div></div></div></div>")
	fmt.Fprintf(w, "<div class=\"stream_pcm\" id=\"streamPcm\"><div class=\"stream_pcm_title\" id=\"streamPcmTitle\"></div><div class=\"stream_pcm_content\"><div class=\"stream_pcm_type\"><span class=\"stream_pcm_type_title\" id=\"streamPcmTypeTitle\"></span><span class=\"stream_pcm_type_value\" id=\"streamPcmTypeValue\"></span></div><div class=\"stream_pcm_content_num\"><span class=\"stream_pcm_content_num_title\" id=\"streamPcmContentNumTitle\"></span><span class=\"stream_pcm_content_num_value\">1</span></div><div class=\"stream_pcm_content_time\"><span class=\"stream_pcm_content_time_title\" id=\"streamPcmContentTimeTitle\"></span><span class=\"stream_pcm_content_time_value\" id=\"streamPcmContentTimeValue\"></span></div><div class=\"stream_pcm_response\"><span class=\"stream_pcm_response_title\" id=\"streamPcmResponseTitle\"></span><br><span class=\"stream_pcm_response_value\" id=\"streamPcmResponseValue\"></span></div></div></div>")
	fmt.Fprintf(w, "<div class=\"info\" id=\"info\"></div><div class=\"showStreamPcmBtnContainer\" id=\"showStreamPcmBtnContainer\"><button type=\"button\" id=\"showStreamPcmBtn\" class=\"showStreamPcmBtn\"></button></div><div class=\"hideStreamPcmBtnContainer\" id=\"hideStreamPcmBtnContainer\"><button type=\"button\" id=\"hideStreamPcmBtn\" class=\"hideStreamPcmBtn\"></button></div><div class=\"footer\"><select id=\"languageSelect\" class=\"language-select\"><option value=\"zh-CN\">简体中文</option><option value=\"en\">English</option></select><div class=\"copyright\" id=\"copyright\"></div></div></div>")
	fmt.Fprintf(w, "<script>")
	// Set copyright year and read head meta tags
	fmt.Fprintf(w, "const currentYear = new Date().getFullYear();var head = document.getElementsByTagName('head')[0];")
	// language definition
	fmt.Fprintf(w, "const titles = {'zh-CN': '%s','en': '%s'};", websiteNameCN, websiteNameEN)
	fmt.Fprintf(w, "const titles2 = {'zh-CN': '%s','en': '%s'};", websiteTitleCN, websiteTitleEN)
	fmt.Fprintf(w, "const descriptions = {'zh-CN': '%s','en': '%s'};", websiteDescCN, websiteDescEN)
	fmt.Fprintf(w, "const keywords = {'zh-CN': '%s','en': '%s'};", websiteKeywordsCN, websiteKeywordsEN)
	fmt.Fprintf(w, "const songInputs = {'zh-CN': '输入歌曲名称...','en': 'Enter song name...'};")
	fmt.Fprintf(w, "const singerInputs = {'zh-CN': '歌手名称(可选)','en': 'Singer name(optional)'};")
	fmt.Fprintf(w, "const searchBtns = {'zh-CN': '<i class=\"fa fa-search\"></i> 搜索','en': '<i class=\"fa fa-search\"></i> Search'};")
	fmt.Fprintf(w, "const getErrors = {'zh-CN': '获取数据失败<br>可能是因为网络响应出错或其它原因<br>请检查您的网络并稍后再试','en': 'Failed to get data<br>It may be because of network response error or other reasons<br>Please check your network and try again later'};")
	fmt.Fprintf(w, "const noEnters = {'zh-CN': '请输入歌曲名称','en': 'Please enter song name'};")
	fmt.Fprintf(w, "const noResults = {'zh-CN': '没有找到相关歌曲','en': 'No related songs found'};")
	fmt.Fprintf(w, "const resultTitles = {'zh-CN': '<i class=\"fa fa-list-ul\"></i> 搜索结果','en': '<i class=\"fa fa-list-ul\"></i> Search Result'};")
	fmt.Fprintf(w, "const caches = {'zh-CN': '缓存','en': 'Cache'};")
	fmt.Fprintf(w, "const noLyrics = {'zh-CN': '暂无歌词','en': 'No lyrics'};")
	fmt.Fprintf(w, "const playBtns = {'zh-CN': '<i class=\"fa fa-play-circle-o\"></i> 播放','en': '<i class=\"fa fa-play-circle-o\"></i> Play'};")
	fmt.Fprintf(w, "const pauseBtns = {'zh-CN': '<i class=\"fa fa-pause-circle-o\"></i> 暂停','en': '<i class=\"fa fa-pause-circle-o\"></i> Pause'};")
	fmt.Fprintf(w, "const streamPcmTitle = {'zh-CN': '<i class=\"fa fa-info-circle\"></i> stream_pcm 响应讯息：','en': '<i class=\"fa fa-info-circle\"></i> stream_pcm response: '};")
	fmt.Fprintf(w, "const streamPcmTypeTitle = {'zh-CN': '响应类型：','en': 'Response type: '};")
	fmt.Fprintf(w, "const streamPcmTypeValue = {'zh-CN': '单曲播放讯息','en': 'Single song playback message'};")
	fmt.Fprintf(w, "const streamPcmContentNumTitle = {'zh-CN': '响应数量：','en': 'Response number: '};")
	fmt.Fprintf(w, "const streamPcmContentTimeTitle = {'zh-CN': '响应时间：','en': 'Response time: '};")
	fmt.Fprintf(w, "const streamPcmResponseTitle = {'zh-CN': '完整响应：','en': 'Full response: '};")
	fmt.Fprintf(w, "const info = {'zh-CN': '<strong><i class=\"fa fa-info-circle\"></i> 系统讯息</strong><br>嵌入式音乐搜索服务器 | Ver %s<br>支持云端/本地音乐搜索，支持多种音乐格式播放，支持多种语言<br>基于聚合API，支持本地音乐缓存','en': '<strong><i class=\"fa fa-info-circle\"></i> System Information</strong><br>Embedded Music Search Server | Ver %s<br>Support cloud/local music search, support various music formats, support various languages<br>Based on aggregation API, support local music cache'};", websiteVersion, websiteVersion)
	fmt.Fprintf(w, "const showStreamPcmBtns = {'zh-CN': '<i class=\"fa fa-eye\"></i> 显示 stream_pcm 响应','en': '<i class=\"fa fa-eye\"></i> Show stream_pcm response'};")
	fmt.Fprintf(w, "const hideStreamPcmBtns = {'zh-CN': '<i class=\"fa fa-eye-slash\"></i> 隐藏 stream_pcm 响应','en': '<i class=\"fa fa-eye-slash\"></i> Hide stream_pcm response'};")
	// Get browser language, set HTML lang attribute and Set default language
	fmt.Fprintf(w, "const browserLang = navigator.language || 'en';document.documentElement.lang = browserLang || \"en\";document.getElementById('languageSelect').value = browserLang;")
	// Initialize title
	fmt.Fprintf(w, "document.title = (titles[browserLang] || '%s') + \" - \" + (titles2[browserLang] || '%s');", websiteNameEN, websiteTitleEN)
	// Initialize meta description
	fmt.Fprintf(w, "var existingMetaDescription = document.querySelector('meta[name=\"description\"]');if (existingMetaDescription) {existingMetaDescription.content = descriptions[browserLang] || '%s';} else {var metaDescription = document.createElement('meta');metaDescription.name = 'description';metaDescription.content = descriptions[browserLang] || '%s';head.appendChild(metaDescription);};", websiteDescEN, websiteDescEN)
	// Initialize meta keywords
	fmt.Fprintf(w, "var existingMetaKeywords = document.querySelector('meta[name=\"keywords\"]');if (existingMetaKeywords) {existingMetaKeywords.content = keywords[browserLang] || '%s';} else {var metaKeywords = document.createElement('meta');metaKeywords.name = 'keywords';metaKeywords.content = keywords[browserLang] || '%s';head.appendChild(metaKeywords);};", websiteKeywordsEN, websiteKeywordsEN)
	// Set default language content
	fmt.Fprintf(w, "document.getElementById('title').innerHTML = titles[browserLang] || '%s';", websiteNameEN)
	fmt.Fprintf(w, "document.getElementById('copyright').innerHTML = \"&copy;\" + currentYear + \" \" + (titles[browserLang] || '%s');", websiteNameEN)
	fmt.Fprintf(w, "document.getElementById('description').innerHTML = descriptions[browserLang] || '%s';", websiteDescEN)
	fmt.Fprintf(w, "document.getElementById('songInput').placeholder = songInputs[browserLang] || 'Enter song name...';")
	fmt.Fprintf(w, "document.getElementById('artistInput').placeholder = singerInputs[browserLang] || 'Singer name(optional)';")
	fmt.Fprintf(w, "document.getElementById('searchBtn').innerHTML = searchBtns[browserLang] || '<i class=\"fa fa-search\"></i> Search';")
	fmt.Fprintf(w, "document.getElementById('getError').innerHTML = getErrors[browserLang] || 'Failed to get data<br>It may be because of network response error or other reasons<br>Please check your network and try again later';")
	fmt.Fprintf(w, "document.getElementById('noEnter').innerHTML = noEnters[browserLang] || 'Please enter song name';")
	fmt.Fprintf(w, "document.getElementById('noResult').innerHTML = noResults[browserLang] || 'No related songs found';")
	fmt.Fprintf(w, "document.getElementById('resultTitle').innerHTML = resultTitles[browserLang] || '<i class=\"fa fa-list-ul\"></i> Search Result';")
	fmt.Fprintf(w, "document.getElementById('cache').innerHTML = caches[browserLang] || 'Cache';")
	fmt.Fprintf(w, "document.getElementById('noLyric').innerHTML = noLyrics[browserLang] || 'No lyrics';")
	fmt.Fprintf(w, "document.getElementById('playBtn').innerHTML = playBtns[browserLang] || '<i class=\"fa fa-play-circle-o\"></i> Play';")
	fmt.Fprintf(w, "document.getElementById('pauseBtn').innerHTML = pauseBtns[browserLang] || '<i class=\"fa fa-pause-circle-o\"></i> Pause';")
	fmt.Fprintf(w, "document.getElementById('streamPcmTitle').innerHTML = streamPcmTitle[browserLang] || '<i class=\"fa fa-info-circle\"></i> stream_pcm response: ';")
	fmt.Fprintf(w, "document.getElementById('streamPcmTypeTitle').innerHTML = streamPcmTypeTitle[browserLang] || 'Response type: ';")
	fmt.Fprintf(w, "document.getElementById('streamPcmTypeValue').innerHTML = streamPcmTypeValue[browserLang] || 'Single song playback message';")
	fmt.Fprintf(w, "document.getElementById('streamPcmContentNumTitle').innerHTML = streamPcmContentNumTitle[browserLang] || 'Response number: ';")
	fmt.Fprintf(w, "document.getElementById('streamPcmContentTimeTitle').innerHTML = streamPcmContentTimeTitle[browserLang] || 'Response time: ';")
	fmt.Fprintf(w, "document.getElementById('streamPcmResponseTitle').innerHTML = streamPcmResponseTitle[browserLang] || 'Full response: ';")
	fmt.Fprintf(w, "document.getElementById('info').innerHTML = info[browserLang] || '<strong><i class=\"fa fa-info-circle\"></i> System Information</strong><br>Embedded Music Search Server | Ver %s<br>Support cloud/local music search, support various music formats, support various languages<br>Based on aggregation API, support local music cache';", websiteVersion)
	fmt.Fprintf(w, "document.getElementById('showStreamPcmBtn').innerHTML = showStreamPcmBtns[browserLang] || '<i class=\"fa fa-eye\"></i> Show stream_pcm response';")
	fmt.Fprintf(w, "document.getElementById('hideStreamPcmBtn').innerHTML = hideStreamPcmBtns[browserLang] || '<i class=\"fa fa-eye-slash\"></i> Hide stream_pcm response';")
	// Listen language selection change and update title
	fmt.Fprintf(w, "document.getElementById('languageSelect').addEventListener('change', function () {")
	fmt.Fprintf(w, "const selectedLang = this.value;")
	// Set HTML lang attribute
	fmt.Fprintf(w, "document.documentElement.lang = selectedLang || \"en\";")
	// Set title
	fmt.Fprintf(w, "document.title = (titles[selectedLang] || '%s') + \" - \" + (titles2[selectedLang] || '%s');", websiteNameEN, websiteTitleEN)
	// Initialize meta description
	fmt.Fprintf(w, "var existingMetaDescription = document.querySelector('meta[name=\"description\"]');if (existingMetaDescription) {existingMetaDescription.content = descriptions[selectedLang] || '%s';} else {var metaDescription = document.createElement('meta');metaDescription.name = 'description';metaDescription.content = descriptions[selectedLang] || '%s';head.appendChild(metaDescription);};", websiteDescEN, websiteDescEN)
	// Initialize meta keywords
	fmt.Fprintf(w, "var existingMetaKeywords = document.querySelector('meta[name=\"keywords\"]');if (existingMetaKeywords) {existingMetaKeywords.content = keywords[selectedLang] || '%s';} else {var metaKeywords = document.createElement('meta');metaKeywords.name = 'keywords';metaKeywords.content = keywords[selectedLang] || '%s';head.appendChild(metaKeywords);};", websiteKeywordsEN, websiteKeywordsEN)
	// Set default language content
	fmt.Fprintf(w, "document.getElementById('title').innerHTML = titles[selectedLang] || '%s';", websiteNameEN)
	fmt.Fprintf(w, "document.getElementById('copyright').innerHTML = \"&copy;\" + currentYear + \" \" + (titles[selectedLang] || '%s');", websiteNameEN)
	fmt.Fprintf(w, "document.getElementById('description').innerHTML = descriptions[selectedLang] || '%s';", websiteDescEN)
	fmt.Fprintf(w, "document.getElementById('songInput').placeholder = songInputs[selectedLang] || 'Enter song name...';")
	fmt.Fprintf(w, "document.getElementById('artistInput').placeholder = singerInputs[selectedLang] || 'Singer name(optional)';")
	fmt.Fprintf(w, "document.getElementById('searchBtn').innerHTML = searchBtns[selectedLang] || '<i class=\"fa fa-search\"></i> Search';")
	fmt.Fprintf(w, "document.getElementById('getError').innerHTML = getErrors[selectedLang] || 'Failed to get data<br>It may be because of network response error or other reasons<br>Please check your network and try again later';")
	fmt.Fprintf(w, "document.getElementById('noEnter').innerHTML = noEnters[selectedLang] || 'Please enter song name';")
	fmt.Fprintf(w, "document.getElementById('noResult').innerHTML = noResults[selectedLang] || 'No related songs found';")
	fmt.Fprintf(w, "document.getElementById('resultTitle').innerHTML = resultTitles[selectedLang] || '<i class=\"fa fa-list-ul\"></i> Search Result';")
	fmt.Fprintf(w, "document.getElementById('cache').innerHTML = caches[selectedLang] || 'Cache';")
	fmt.Fprintf(w, "document.getElementById('noLyric').innerHTML = noLyrics[selectedLang] || 'No lyrics';")
	fmt.Fprintf(w, "document.getElementById('playBtn').innerHTML = playBtns[selectedLang] || '<i class=\"fa fa-play-circle-o\"></i> Play';")
	fmt.Fprintf(w, "document.getElementById('pauseBtn').innerHTML = pauseBtns[selectedLang] || '<i class=\"fa fa-pause-circle-o\"></i> Pause';")
	fmt.Fprintf(w, "document.getElementById('streamPcmTitle').innerHTML = streamPcmTitle[selectedLang] || '<i class=\"fa fa-info-circle\"></i> stream_pcm response: ';")
	fmt.Fprintf(w, "document.getElementById('streamPcmTypeTitle').innerHTML = streamPcmTypeTitle[selectedLang] || 'Response type: ';")
	fmt.Fprintf(w, "document.getElementById('streamPcmTypeValue').innerHTML = streamPcmTypeValue[selectedLang] || 'Single song playback message';")
	fmt.Fprintf(w, "document.getElementById('streamPcmContentNumTitle').innerHTML = streamPcmContentNumTitle[selectedLang] || 'Response number: ';")
	fmt.Fprintf(w, "document.getElementById('streamPcmContentTimeTitle').innerHTML = streamPcmContentTimeTitle[selectedLang] || 'Response time: ';")
	fmt.Fprintf(w, "document.getElementById('streamPcmResponseTitle').innerHTML = streamPcmResponseTitle[selectedLang] || 'Full response: ';")
	fmt.Fprintf(w, "document.getElementById('info').innerHTML = info[selectedLang] || '<strong><i class=\"fa fa-info-circle\"></i> System Information</strong><br>Embedded Music Search Server | Ver %s<br>Support cloud/local music search, support various music formats, support various languages<br>Based on aggregation API, support local music cache';", websiteVersion)
	fmt.Fprintf(w, "document.getElementById('showStreamPcmBtn').innerHTML = showStreamPcmBtns[selectedLang] || '<i class=\"fa fa-eye\"></i> Show stream_pcm response';")
	fmt.Fprintf(w, "document.getElementById('hideStreamPcmBtn').innerHTML = hideStreamPcmBtns[selectedLang] || '<i class=\"fa fa-eye-slash\"></i> Hide stream_pcm response';")
	fmt.Fprintf(w, "});")
	// Getting Elements
	fmt.Fprintf(w, "const songInput = document.getElementById('songInput');")
	fmt.Fprintf(w, "const artistInput = document.getElementById('artistInput');")
	fmt.Fprintf(w, "const searchBtn = document.getElementById('searchBtn');")
	fmt.Fprintf(w, "const getError = document.getElementById('getError');")
	fmt.Fprintf(w, "const noEnter = document.getElementById('noEnter');")
	fmt.Fprintf(w, "const noResult = document.getElementById('noResult');")
	fmt.Fprintf(w, "const loading = document.getElementById('loading');")
	fmt.Fprintf(w, "const result = document.getElementById('result');")
	fmt.Fprintf(w, "const songName = document.getElementById('songName');")
	fmt.Fprintf(w, "const cache = document.getElementById('cache');")
	fmt.Fprintf(w, "const singerName = document.getElementById('singerName');")
	fmt.Fprintf(w, "const noLyric = document.getElementById('noLyric');")
	fmt.Fprintf(w, "const playBtn = document.getElementById('playBtn');")
	fmt.Fprintf(w, "const pauseBtn = document.getElementById('pauseBtn');")
	fmt.Fprintf(w, "const lyric = document.getElementById('lyric');")
	fmt.Fprintf(w, "const streamPcm = document.getElementById('streamPcm');")
	fmt.Fprintf(w, "const streamPcmContentTimeValue = document.getElementById('streamPcmContentTimeValue');")
	fmt.Fprintf(w, "const streamPcmResponseValue = document.getElementById('streamPcmResponseValue');")
	fmt.Fprintf(w, "const showStreamPcmBtn = document.getElementById('showStreamPcmBtn');")
	fmt.Fprintf(w, "const hideStreamPcmBtn = document.getElementById('hideStreamPcmBtn');")
	// Hide content that should not be displayed before searching
	fmt.Fprintf(w, "getError.style.display = 'none';")
	fmt.Fprintf(w, "noEnter.style.display = 'none';")
	fmt.Fprintf(w, "noResult.style.display = 'none';")
	fmt.Fprintf(w, "loading.style.display = 'none';")
	fmt.Fprintf(w, "result.style.display = 'none';")
	fmt.Fprintf(w, "cache.style.display = 'none';")
	fmt.Fprintf(w, "noLyric.style.display = 'none';")
	fmt.Fprintf(w, "playBtn.style.display = 'none';")
	fmt.Fprintf(w, "pauseBtn.style.display = 'none';")
	fmt.Fprintf(w, "streamPcm.style.display = 'none';")
	fmt.Fprintf(w, "showStreamPcmBtn.style.display = 'none';")
	fmt.Fprintf(w, "hideStreamPcmBtn.style.display = 'none';")
	// Empty song name processing
	fmt.Fprintf(w, "searchBtn.addEventListener('click', function () {if (songInput.value.trim() === '') {noEnter.style.display = 'block';} else {noEnter.style.display = 'none';search();}});")
	fmt.Fprintf(w, "songInput.addEventListener('keydown', function (event) {if (event.key === 'Enter') {if (songInput.value.trim() === '') {noEnter.style.display = 'block';} else {noEnter.style.display = 'none';search();}}});")
	fmt.Fprintf(w, "artistInput.addEventListener('keydown', function (event) {if (event.key === 'Enter') {if (songInput.value.trim() === '') {noEnter.style.display = 'block';} else {noEnter.style.display = 'none';search();}}});")
	// Searching for songs
	fmt.Fprintf(w, "function search() {")
	// Show loading
	fmt.Fprintf(w, "loading.style.display = 'block';")
	// Hide error
	fmt.Fprintf(w, "getError.style.display = 'none';")
	// Build request URL, urlencode song name and artist name
	fmt.Fprintf(w, "const song = encodeURIComponent(songInput.value);")
	fmt.Fprintf(w, "const artist = encodeURIComponent(artistInput.value);")
	fmt.Fprintf(w, "const requestUrl = `/stream_pcm?song=${song}&artist=${artist}`;")
	// Send request to server
	fmt.Fprintf(w, "fetch(requestUrl)")
	fmt.Fprintf(w, ".then(response => {if (!response.ok) {getError.style.display = 'block';throw new Error('Network response was not ok');}return response.json();})")
	fmt.Fprintf(w, ".then(data => {")
	// Fill in all the obtained content into streamPcmResponseValue
	fmt.Fprintf(w, "streamPcmResponseValue.innerHTML = JSON.stringify(data, null, 2);")
	// Get the current time and fill in streamPcmCntentTimeValue
	fmt.Fprintf(w, "streamPcmContentTimeValue.innerHTML = new Date().toISOString();")
	// Display result
	fmt.Fprintf(w, "result.style.display = 'block';")
	// Display Play button
	fmt.Fprintf(w, "playBtn.style.display = 'block';")
	// Display showStreamPcmBtn
	fmt.Fprintf(w, "showStreamPcmBtn.style.display = 'block';")
	// Hide hideStreamPcmBtn
	fmt.Fprintf(w, "hideStreamPcmBtn.style.display = 'none';")
	// Fill the title into the songName field
	fmt.Fprintf(w, "if (data.title === \"\") {noResult.style.display = 'block';result.style.display = 'none';} else {noResult.style.display = 'none';songName.textContent = data.title;};")
	// Fill the artist into the singerName field
	fmt.Fprintf(w, "singerName.textContent = data.artist;")
	// Set parsed lyrics to an empty array
	fmt.Fprintf(w, "let parsedLyrics = [];")
	// Check if the link 'lyric_url' is empty
	fmt.Fprintf(w, "if (data.lyric_url) {")
	// Visit lyric_url
	fmt.Fprintf(w, "fetch(data.lyric_url)")
	fmt.Fprintf(w, ".then(response => {if (!response.ok) {throw new Error('Lyrics request error');}return response.text();})")
	fmt.Fprintf(w, ".then(lyricText => {")
	// Show lyric
	fmt.Fprintf(w, "lyric.style.display = 'block';")
	// Parse lyrics
	fmt.Fprintf(w, "parsedLyrics = parseLyrics(lyricText);")
	fmt.Fprintf(w, "})")
	fmt.Fprintf(w, ".catch(error => {")
	// Show noLyric
	fmt.Fprintf(w, "noLyric.style.display = 'block';")
	// Hide lyric
	fmt.Fprintf(w, "lyric.style.display = 'none';")
	fmt.Fprintf(w, "});")
	fmt.Fprintf(w, "} else {")
	// If lyric_url is empty, display noLyric
	fmt.Fprintf(w, "noLyric.style.display = 'block';")
	// Hide lyric
	fmt.Fprintf(w, "lyric.style.display = 'none';")
	fmt.Fprintf(w, "};")
	// Check cache
	fmt.Fprintf(w, "if (data.from_cache === true) {")
	// If the song is obtained from cache, display cache
	fmt.Fprintf(w, "cache.style.display = 'block';")
	fmt.Fprintf(w, "} else {")
	// If the song is not obtained from cache, hide cache
	fmt.Fprintf(w, "cache.style.display = 'none';")
	fmt.Fprintf(w, "};")
	// Create audio player
	fmt.Fprintf(w, "const audioPlayer = document.getElementById('audio');")
	fmt.Fprintf(w, "const customProgress = document.getElementById('progress');")
	fmt.Fprintf(w, "const timeDisplay = document.getElementById('time');")
	// Set audio source
	fmt.Fprintf(w, "audioPlayer.src = data.audio_full_url;")
	fmt.Fprintf(w, "audio.addEventListener('timeupdate', function () {")
	fmt.Fprintf(w, "const currentTime = formatTime(audioPlayer.currentTime);")
	fmt.Fprintf(w, "const duration = formatTime(audioPlayer.duration);")
	fmt.Fprintf(w, "timeDisplay.textContent = `${currentTime}/${duration}`;")
	fmt.Fprintf(w, "const progress = (audioPlayer.currentTime / audioPlayer.duration) * 100;")
	fmt.Fprintf(w, "customProgress.style.width = progress + '%%';")
	// Find current lyric
	fmt.Fprintf(w, "let currentLyric = parsedLyrics.find((lyric, index, arr) =>")
	fmt.Fprintf(w, "lyric.timestamp > audioPlayer.currentTime && (index === 0 || arr[index - 1].timestamp <= audioPlayer.currentTime));")
	fmt.Fprintf(w, "lyric.textContent = currentLyric ? currentLyric.lyricLine : '';")
	fmt.Fprintf(w, "});")
	// Save current time before playing audio
	fmt.Fprintf(w, "let savedCurrentTime = 0;")
	// PlayBtn click event
	fmt.Fprintf(w, "playBtn.addEventListener('click', function () {")
	// Save current time
	fmt.Fprintf(w, "audioPlayer.currentTime = savedCurrentTime;")
	// Play audio
	fmt.Fprintf(w, "audioPlayer.play();")
	// Hide playBtn
	fmt.Fprintf(w, "playBtn.style.display = 'none';")
	// Show pauseBtn
	fmt.Fprintf(w, "pauseBtn.style.display = 'block';")
	fmt.Fprintf(w, "});")
	// PauseBtn click event
	fmt.Fprintf(w, "pauseBtn.addEventListener('click', function () {")
	// Save current time
	fmt.Fprintf(w, "savedCurrentTime = audioPlayer.currentTime;")
	// Pause audio
	fmt.Fprintf(w, "audioPlayer.pause();")
	// Hide pauseBtn
	fmt.Fprintf(w, "pauseBtn.style.display = 'none';")
	// Show playBtn
	fmt.Fprintf(w, "playBtn.style.display = 'block';")
	fmt.Fprintf(w, "});")
	fmt.Fprintf(w, "})")
	fmt.Fprintf(w, ".catch(error => {")
	fmt.Fprintf(w, "console.error('Error requesting song information:', error);")
	// When there is an error in the request, you can also consider displaying a prompt message or other content
	fmt.Fprintf(w, "getError.style.display = 'block';")
	fmt.Fprintf(w, "})")
	fmt.Fprintf(w, ".finally(() => {")
	// Regardless of the request result, the loading display should be turned off in the end
	fmt.Fprintf(w, "loading.style.display = 'none';")
	fmt.Fprintf(w, "});};")
	// Format time
	fmt.Fprintf(w, "function formatTime(seconds) {const minutes = Math.floor(seconds / 60);const secondsRemainder = Math.floor(seconds %% 60);return minutes.toString().padStart(2, '0') + ':' +secondsRemainder.toString().padStart(2, '0');};")
	// Function to parse lyrics
	fmt.Fprintf(w, "function parseLyrics(lyricText) {const lines = lyricText.split('\\n');const lyrics = [];for (let line of lines) {const match = line.match(/\\[(\\d{2}:\\d{2})(?:\\.\\d{2})?\\](.*)/);if (match) {const timestamp = match[1]; const lyricLine = match[2].trim();const [minutes, seconds] = timestamp.split(':');const timeInSeconds = (parseInt(minutes) * 60) + parseInt(seconds);lyrics.push({ timestamp: timeInSeconds, lyricLine });}}return lyrics;};")
	// Show stream_pcm response
	fmt.Fprintf(w, "showStreamPcmBtn.addEventListener('click', function () {streamPcm.style.display = 'block';showStreamPcmBtn.style.display = 'none';hideStreamPcmBtn.style.display = 'block';});")
	// Hide stream_pcm response
	fmt.Fprintf(w, "hideStreamPcmBtn.addEventListener('click', function () {streamPcm.style.display = 'none';showStreamPcmBtn.style.display = 'block';hideStreamPcmBtn.style.display = 'none';});")
	fmt.Fprintf(w, "</script></body></html>")
}
