package bucketize

import (
	"encoding/json"
	"fmt"
	"io"
)

type jsonBucket struct {
	Label     string  `json:"label"`
	Start     string  `json:"start"`
	Total     int     `json:"total"`
	Drifted   int     `json:"drifted"`
	DriftRate float64 `json:"drift_rate"`
}

// WriteJSON encodes buckets as a JSON array to w.
func WriteJSON(w io.Writer, buckets []Bucket) error {
	out := make([]jsonBucket, len(buckets))
	for i, b := range buckets {
		out[i] = jsonBucket{
			Label:     b.Label,
			Start:     b.Start.Format("2006-01-02T15:04:05Z"),
			Total:     b.Total,
			Drifted:   b.Drifted,
			DriftRate: b.DriftRate(),
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("bucketize: encode json: %w", err)
	}
	return nil
}
