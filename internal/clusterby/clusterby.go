// Package clusterby groups drift results into clusters based on a shared
// attribute such as service prefix, field name, or drift status.
package clusterby

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/example/driftwatch/internal/drift"
)

// Key is the dimension used for clustering.
type Key string

const (
	KeyField   Key = "field"
	KeyStatus  Key = "status"
	KeyPrefix  Key = "prefix"
)

// Cluster holds a named group of drift results.
type Cluster struct {
	Name    string
	Results []drift.Result
}

// Apply partitions results into clusters by the given key.
// For KeyPrefix the separator is "-"; the first segment of the service name is used.
func Apply(results []drift.Result, by Key) []Cluster {
	if len(results) == 0 {
		return nil
	}

	buckets := make(map[string][]drift.Result)
	for _, r := range results {
		k := clusterKey(r, by)
		buckets[k] = append(buckets[k], r)
	}

	clusters := make([]Cluster, 0, len(buckets))
	for name, rs := range buckets {
		clusters = append(clusters, Cluster{Name: name, Results: rs})
	}
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Name < clusters[j].Name
	})
	return clusters
}

func clusterKey(r drift.Result, by Key) string {
	switch by {
	case KeyField:
		if r.Field == "" {
			return "(none)"
		}
		return r.Field
	case KeyStatus:
		if r.Drifted {
			return "drifted"
		}
		return "ok"
	case KeyPrefix:
		parts := strings.SplitN(r.Service, "-", 2)
		return parts[0]
	default:
		return r.Service
	}
}

// Write renders clusters as a text table to w.
func Write(w io.Writer, clusters []Cluster) {
	if len(clusters) == 0 {
		fmt.Fprintln(w, "no clusters")
		return
	}
	fmt.Fprintf(w, "%-24s %8s %8s\n", "CLUSTER", "TOTAL", "DRIFTED")
	for _, c := range clusters {
		drifted := 0
		for _, r := range c.Results {
			if r.Drifted {
				drifted++
			}
		}
		fmt.Fprintf(w, "%-24s %8d %8d\n", c.Name, len(c.Results), drifted)
	}
}

// WriteJSON encodes clusters as JSON.
func WriteJSON(w io.Writer, clusters []Cluster) error {
	return json.NewEncoder(w).Encode(clusters)
}
