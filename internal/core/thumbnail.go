package core

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func ThumbnailHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	infohash := parts[3]
	if len(infohash) != 40 {
		http.Error(w, "invalid infohash", http.StatusBadRequest)
		return
	}

	thumbDir := filepath.Join(Dirconfig.DataDir, "thumbnails")
	thumbPath := filepath.Join(thumbDir, infohash+".jpg")

	if fi, err := os.Stat(thumbPath); err == nil && fi.Size() > 0 {
		w.Header().Set("Cache-Control", "public, max-age=604800")
		http.ServeFile(w, r, thumbPath)
		return
	}

	torrentDir := filepath.Join(Dirconfig.TrntDir, infohash)
	videoFile := findLargestVideo(torrentDir)
	if videoFile == "" {
		http.Error(w, "no video file found", http.StatusNotFound)
		return
	}

	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		slog.Error("failed to create thumbnail dir", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-ss", "30",
		"-i", videoFile,
		"-vframes", "1",
		"-q:v", "5",
		"-vf", "scale=400:-1",
		"-y",
		thumbPath,
	)
	if err := cmd.Run(); err != nil {
		slog.Debug("ffmpeg thumbnail extraction failed", "error", err, "infohash", infohash)
		http.Error(w, "thumbnail generation failed", http.StatusNotFound)
		return
	}

	if fi, err := os.Stat(thumbPath); err != nil || fi.Size() == 0 {
		http.Error(w, "thumbnail generation failed", http.StatusNotFound)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=604800")
	http.ServeFile(w, r, thumbPath)
}

var videoExtensions = map[string]bool{
	".mkv": true, ".mp4": true, ".avi": true, ".mov": true, ".wmv": true,
	".flv": true, ".webm": true, ".m4v": true, ".mpg": true, ".mpeg": true,
	".3gp": true, ".ts": true, ".vob": true,
}

func findLargestVideo(dir string) string {
	var largest string
	var maxSize int64
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if videoExtensions[ext] && info.Size() > maxSize {
			maxSize = info.Size()
			largest = path
		}
		return nil
	})
	return largest
}
