// Package bucketize groups drift results into time-based buckets for
// trend analysis and visualisation.
package bucketize

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/driftwatch/internal/drift"
)

// Granularity controls the width of each bucket.
type Granularity string

const (
	Hourly  Granularity = "hourly"
	Daily   Granularity = "daily"
	Weekly  Granularity = "weekly"
)

// Bucket holds aggregated counts for a single time window.
type Bucket struct {
	Label    string
	Start    time.Time
	Total    int
	Drifted  int
}

// DriftRate returns the fraction of drifted results in the bucket.
func (b Bucket) DriftRate() float64 {
	if b.Total == 0 {
		return 0
	}
	return float64(b.Drifted) / float64(b.Total)
}

// Apply groups results by the chosen granularity and returns sorted buckets.
func Apply(results []drift.Result, g Granularity) []Bucket {
	index := map[string]*Bucket{}
	for _, r := range results {
		label, start := key(r.DetectedAt, g)
		b, ok := index[label]
		if !ok {
			b = &Bucket{Label: label, Start: start}
			index[label] = b
		}
		b.Total++
		if r.Drifted {
			b.Drifted++
		}
	}
	buckets := make([]Bucket, 0, len(index))
	for _, b := range index {
		buckets = append(buckets, *b)
	}
	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].Start.Before(buckets[j].Start)
	})
	return buckets
}

func key(t time.Time, g Granularity) (string, time.Time) {
	switch g {
	case Hourly:
		start := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
		return start.Format("2006-01-02T15"), start
	case Weekly:
		year, week := t.ISOWeek()
		label := fmt.Sprintf("%d-W%02d", year, week)
		start := weekStart(t)
		return label, start
	default: // daily
		start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		return start.Format("2006-01-02"), start
	}
}

func weekStart(t time.Time) time.Time {
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	delta := time.Duration(wd-1) * 24 * time.Hour
	d := t.Add(-delta)
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, t.Location())
}

// Write renders the buckets as a text table to w.
func Write(w io.Writer, buckets []Bucket) {
	fmt.Fprintf(w, "%-20s  %6s  %7s  %8s\n", "BUCKET", "TOTAL", "DRIFTED", "RATE")
	for _, b := range buckets {
		fmt.Fprintf(w, "%-20s  %6d  %7d  %7.1f%%\n",
			b.Label, b.Total, b.Drifted, b.DriftRate()*100)
	}
}
