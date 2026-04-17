package bucketize_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/driftwatch/internal/bucketize"
	"github.com/driftwatch/internal/drift"
)

func makeResult(drifted bool, at time.Time) drift.Result {
	return drift.Result{Drifted: drifted, DetectedAt: at}
}

var base = time.Date(2024, 6, 10, 14, 30, 0, 0, time.UTC)

func TestApply_Daily_GroupsCorrectly(t *testing.T) {
	results := []drift.Result{
		makeResult(true, base),
		makeResult(false, base.Add(1*time.Hour)),
		makeResult(true, base.Add(25*time.Hour)),
	}
	buckets := bucketize.Apply(results, bucketize.Daily)
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
}

func TestApply_Daily_Totals(t *testing.T) {
	results := []drift.Result{
		makeResult(true, base),
		makeResult(false, base.Add(2*time.Hour)),
	}
	buckets := bucketize.Apply(results, bucketize.Daily)
	if buckets[0].Total != 2 {
		t.Errorf("expected total 2, got %d", buckets[0].Total)
	}
	if buckets[0].Drifted != 1 {
		t.Errorf("expected drifted 1, got %d", buckets[0].Drifted)
	}
}

func TestApply_Hourly_SplitsOnHour(t *testing.T) {
	results := []drift.Result{
		makeResult(true, base),
		makeResult(true, base.Add(90*time.Minute)),
	}
	buckets := bucketize.Apply(results, bucketize.Hourly)
	if len(buckets) != 2 {
		t.Fatalf("expected 2 hourly buckets, got %d", len(buckets))
	}
}

func TestApply_Weekly_GroupsSameWeek(t *testing.T) {
	mon := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	fri := time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC)
	nextMon := time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC)
	results := []drift.Result{
		makeResult(false, mon),
		makeResult(true, fri),
		makeResult(true, nextMon),
	}
	buckets := bucketize.Apply(results, bucketize.Weekly)
	if len(buckets) != 2 {
		t.Fatalf("expected 2 weekly buckets, got %d", len(buckets))
	}
}

func TestApply_EmptyInput_ReturnsNil(t *testing.T) {
	buckets := bucketize.Apply(nil, bucketize.Daily)
	if len(buckets) != 0 {
		t.Errorf("expected empty, got %d", len(buckets))
	}
}

func TestApply_SortedOldestFirst(t *testing.T) {
	results := []drift.Result{
		makeResult(false, base.Add(48*time.Hour)),
		makeResult(true, base),
	}
	buckets := bucketize.Apply(results, bucketize.Daily)
	if !buckets[0].Start.Before(buckets[1].Start) {
		t.Error("expected buckets sorted oldest first")
	}
}

func TestDriftRate_ZeroTotal(t *testing.T) {
	b := bucketize.Bucket{}
	if b.DriftRate() != 0 {
		t.Error("expected zero drift rate for empty bucket")
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	buckets := []bucketize.Bucket{
		{Label: "2024-06-10", Total: 4, Drifted: 2},
	}
	bucketize.Write(&buf, buckets)
	out := buf.String()
	for _, want := range []string{"BUCKET", "TOTAL", "DRIFTED", "RATE"} {
		if !bytes.Contains([]byte(out), []byte(want)) {
			t.Errorf("expected header %q in output", want)
		}
	}
}
