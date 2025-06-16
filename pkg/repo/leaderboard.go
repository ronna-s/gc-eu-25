package repo

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
)

type ScoreEntry struct {
	BandName string
	Score    int
}

func (e ScoreEntry) Less(e2 ScoreEntry) bool {
	return e.Score > e2.Score
}

func GetTop(n int) ([]ScoreEntry, error) {
	scores, err := GetAll()
	if err != nil {
		return nil, err
	}
	if n > len(scores) {
		n = len(scores)
	}
	return scores[:n], nil
}

func GetAll() ([]ScoreEntry, error) {
	// Find the CSV file relative to this file
	csvPath := filepath.Join("pkg", "repo", "allscores.csv")
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var leaderboard []ScoreEntry
	for _, rec := range records {
		if len(rec) < 2 {
			continue
		}
		score, err := strconv.Atoi(rec[1])
		if err != nil {
			continue
		}
		leaderboard = append(leaderboard, ScoreEntry{
			BandName: rec[0],
			Score:    score,
		})
	}
	return leaderboard, nil
}

// Persist appends a new leaderboard entry to the CSV, allowing duplicates and not checking for existing entries.
func Persist(entry ScoreEntry) error {
	csvPath := filepath.Join("pkg", "pnp", "repo", "allscores.csv")
	file, err := os.OpenFile(csvPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err == nil && stat.Size() > 0 {
		// Check if the last byte is a newline
		if _, err := file.Seek(-1, 2); err == nil { // 2 is io.SeekEnd, use constant for clarity
			var lastByte [1]byte
			if _, err := file.Read(lastByte[:]); err == nil {
				if lastByte[0] != '\n' {
					if _, err := file.Seek(0, 2); err == nil {
						if _, err := file.Write([]byte("\n")); err != nil {
							return err
						}
						if err := file.Sync(); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	// Always seek to end before writing
	if _, err := file.Seek(0, 2); err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	record := []string{entry.BandName, strconv.Itoa(entry.Score)}
	if err := writer.Write(record); err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}
