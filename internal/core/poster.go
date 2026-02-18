package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

var TMDBAPIKey string

type posterCacheEntry struct {
	result    *posterResult
	cachedAt  time.Time
}

var posterCache sync.Map

var (
	reYear    = regexp.MustCompile(`[\.\s\-_](\d{4})[\.\s\-_]`)
	reSeason  = regexp.MustCompile(`(?i)[\.\s\-_]S\d{1,2}`)
	reCleanup = regexp.MustCompile(`(?i)[\.\s\-_](1080[pi]|720[pi]|2160[pi]|4[kK]|480[pi]|BDRip|BRRip|BluRay|WEB[\.\-]?DL|WEB[\.\-]?Rip|WEBRip|HDRip|DVDRip|HDTV|PDTV|XviD|x264|x265|H\.?264|H\.?265|HEVC|AAC|AC3|DTS|FLAC|REMUX|PROPER|REPACK|EXTENDED|UNRATED|DIRECTORS\.?CUT|10bit|HDR|SDR|AMZN|NF|HULU|DSNP|ATVP|iNTERNAL)`)
	reGroup   = regexp.MustCompile(`-[A-Za-z0-9]+$`)
	reBracket = regexp.MustCompile(`[\[\(].*?[\]\)]`)
)

type posterResult struct {
	PosterURL   string `json:"poster_url"`
	BackdropURL string `json:"backdrop_url"`
}

func parseTorrentName(raw string) (title string, year string, isTv bool) {
	name := strings.TrimSpace(raw)
	name = reBracket.ReplaceAllString(name, "")
	name = reGroup.ReplaceAllString(name, "")

	isTv = reSeason.MatchString(name)

	if loc := reSeason.FindStringIndex(name); loc != nil {
		name = name[:loc[0]]
	}

	if m := reYear.FindStringSubmatchIndex(name); m != nil {
		year = name[m[2]:m[3]]
		name = name[:m[0]]
	} else {
		name = reCleanup.Split(name, 2)[0]
	}

	name = strings.NewReplacer(".", " ", "_", " ", "-", " ").Replace(name)

	fields := strings.Fields(name)
	title = strings.Join(fields, " ")
	return
}

func searchTMDB(title, year string, isTv bool) (*posterResult, error) {
	endpoint := "movie"
	if isTv {
		endpoint = "tv"
	}

	params := url.Values{
		"api_key": {TMDBAPIKey},
		"query":   {title},
	}
	if year != "" && !isTv {
		params.Set("year", year)
	}
	if year != "" && isTv {
		params.Set("first_air_date_year", year)
	}

	u := fmt.Sprintf("https://api.themoviedb.org/3/search/%s?%s", endpoint, params.Encode())

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Results []struct {
			PosterPath   string `json:"poster_path"`
			BackdropPath string `json:"backdrop_path"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	result := &posterResult{}
	if len(data.Results) > 0 {
		r := data.Results[0]
		if r.PosterPath != "" {
			result.PosterURL = "https://image.tmdb.org/t/p/w300" + r.PosterPath
		}
		if r.BackdropPath != "" {
			result.BackdropURL = "https://image.tmdb.org/t/p/w780" + r.BackdropPath
		}
	}
	return result, nil
}

func PosterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")
	if name == "" {
		json.NewEncoder(w).Encode(posterResult{})
		return
	}

	if TMDBAPIKey == "" {
		json.NewEncoder(w).Encode(posterResult{})
		return
	}

	title, year, isTv := parseTorrentName(name)
	if title == "" {
		json.NewEncoder(w).Encode(posterResult{})
		return
	}

	cacheKey := strings.ToLower(title + "|" + year)
	if cached, ok := posterCache.Load(cacheKey); ok {
		entry := cached.(*posterCacheEntry)
		json.NewEncoder(w).Encode(entry.result)
		return
	}

	result, err := searchTMDB(title, year, isTv)
	if err != nil {
		slog.Warn("TMDB search failed", "title", title, "error", err)
		json.NewEncoder(w).Encode(posterResult{})
		return
	}

	posterCache.Store(cacheKey, &posterCacheEntry{result: result, cachedAt: time.Now()})

	w.Header().Set("Cache-Control", "public, max-age=86400")
	json.NewEncoder(w).Encode(result)
}
