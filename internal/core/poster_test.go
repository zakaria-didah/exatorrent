package core

import "testing"

func TestParseTorrentName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTitle string
		wantYear  string
		wantIsTv  bool
	}{
		{
			name:      "standard movie with year",
			input:     "The.Matrix.1999.1080p.BluRay.x264-GROUP",
			wantTitle: "The Matrix",
			wantYear:  "1999",
			wantIsTv:  false,
		},
		{
			name:      "movie with dashes instead of dots",
			input:     "The-Matrix-1999-1080p-BluRay",
			wantTitle: "The Matrix",
			wantYear:  "1999",
			wantIsTv:  false,
		},
		{
			name:      "movie with underscores",
			input:     "The_Matrix_1999_1080p_BluRay",
			wantTitle: "The Matrix",
			wantYear:  "1999",
			wantIsTv:  false,
		},
		{
			name:      "TV show with season episode",
			input:     "Breaking.Bad.S01E01.720p.BluRay.x264",
			wantTitle: "Breaking Bad",
			wantYear:  "",
			wantIsTv:  true,
		},
		{
			name:      "TV show with year and season (year after season strip)",
			input:     "The.Mandalorian.2019.S01E01.1080p.WEB-DL",
			wantTitle: "The Mandalorian 2019",
			wantYear:  "",
			wantIsTv:  true,
		},
		{
			name:      "no year present",
			input:     "Ubuntu.Desktop.iso",
			wantTitle: "Ubuntu Desktop iso",
			wantYear:  "",
			wantIsTv:  false,
		},
		{
			name:      "year in title without quality tags",
			input:     "2001.A.Space.Odyssey.1968.1080p.BluRay",
			wantTitle: "2001 A Space Odyssey",
			wantYear:  "1968",
			wantIsTv:  false,
		},
		{
			name:      "empty string",
			input:     "",
			wantTitle: "",
			wantYear:  "",
			wantIsTv:  false,
		},
		{
			name:      "whitespace only",
			input:     "   ",
			wantTitle: "",
			wantYear:  "",
			wantIsTv:  false,
		},
		{
			name:      "brackets in name",
			input:     "[YTS.MX] The.Matrix.1999.1080p.BluRay",
			wantTitle: "The Matrix",
			wantYear:  "1999",
			wantIsTv:  false,
		},
		{
			name:      "parentheses in name",
			input:     "(Official) The.Matrix.1999.1080p",
			wantTitle: "The Matrix",
			wantYear:  "1999",
			wantIsTv:  false,
		},
		{
			name:      "multiple quality tags",
			input:     "Inception.2010.2160p.4K.HDR.BluRay.x265-GROUP",
			wantTitle: "Inception",
			wantYear:  "2010",
			wantIsTv:  false,
		},
		{
			name:      "TV show uppercase S",
			input:     "Game.of.Thrones.S08E06.1080p.WEB",
			wantTitle: "Game of Thrones",
			wantYear:  "",
			wantIsTv:  true,
		},
		{
			name:      "REMUX tag",
			input:     "Dune.2021.REMUX.2160p.BluRay",
			wantTitle: "Dune",
			wantYear:  "2021",
			wantIsTv:  false,
		},
		{
			name:      "streaming service tags",
			input:     "Wednesday.S01E01.1080p.AMZN.WEB-DL",
			wantTitle: "Wednesday",
			wantYear:  "",
			wantIsTv:  true,
		},
		{
			name:      "single word title",
			input:     "Tenet.2020.1080p.WEB-DL",
			wantTitle: "Tenet",
			wantYear:  "2020",
			wantIsTv:  false,
		},
		{
			name:      "very long title",
			input:     "The.Lord.of.the.Rings.The.Return.of.the.King.2003.EXTENDED.1080p.BluRay.x264",
			wantTitle: "The Lord of the Rings The Return of the King",
			wantYear:  "2003",
			wantIsTv:  false,
		},
		{
			name:      "no separator after year",
			input:     "Movie.2023BluRay",
			wantTitle: "Movie 2023BluRay",
			wantYear:  "",
			wantIsTv:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, year, isTv := parseTorrentName(tt.input)
			if title != tt.wantTitle {
				t.Errorf("parseTorrentName(%q) title = %q, want %q", tt.input, title, tt.wantTitle)
			}
			if year != tt.wantYear {
				t.Errorf("parseTorrentName(%q) year = %q, want %q", tt.input, year, tt.wantYear)
			}
			if isTv != tt.wantIsTv {
				t.Errorf("parseTorrentName(%q) isTv = %v, want %v", tt.input, isTv, tt.wantIsTv)
			}
		})
	}
}
