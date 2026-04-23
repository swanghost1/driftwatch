package horizon_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/history"
	"github.com/example/driftwatch/internal/horizon"
)

func makeEntry(recordedAt time.Time, drifted, total int) history.Entry {
	return history.Entry{
		RecordedAt: recordedAt,
		Drifted:    drifted,
		Total:      total,
	}
}

func TestEvaluate_NoEntries_ReturnsNil(t *testing.T) {
	opts := horizon.DefaultOptions()
	got := horizon.Evaluate(nil, time.Now(), opts)
	if got != nil {
		t.Fatalf("expected nil, got %+v", got)
	}
}

func TestEvaluate_AllOutsideWindow_ReturnsNil(t *testing.T) {
	now := time.Now()
	entries := []history.Entry{
		makeEntry(now.Add(-30*24*time.Hour), 2, 10),
	}
	opts := horizon.DefaultOptions() // window = 7 days
	got := horizon.Evaluate(entries, now, opts)
	if got != nil {
		t.Fatalf("expected nil for out-of-window entries, got %+v", got)
	}
}

func TestEvaluate_NoDrift_WillNotExceed(t *testing.T) {
	now := time.Now()
	entries := []history.Entry{
		makeEntry(now.Add(-1*time.Hour), 0, 10),
		makeEntry(now.Add(-2*time.Hour), 0, 10),
	}
	opts := horizon.DefaultOptions()
	got := horizon.Evaluate(entries, now, opts)
	if got == nil {
		t.Fatal("expected prediction, got nil")
	}
	if got.WillExceed {
		t.Errorf("expected WillExceed=false, got true")
	}
	if got.CurrentRate != 0 {
		t.Errorf("expected CurrentRate=0, got %f", got.CurrentRate)
	}
}

func TestEvaluate_HighDriftRate_WillExceed(t *testing.T) {
	now := time.Now()
	entries := []history.Entry{
		makeEntry(now.Add(-1*time.Hour), 9, 10),
		makeEntry(now.Add(-2*time.Hour), 8, 10),
	}
	opts := horizon.DefaultOptions()
	got := horizon.Evaluate(entries, now, opts)
	if got == nil {
		t.Fatal("expected prediction, got nil")
	}
	if !got.WillExceed {
		t.Errorf("expected WillExceed=true for high drift rate")
	}
}

func TestEvaluate_ProjectedRateCappedAtOne(t *testing.T) {
	now := time.Now()
	entries := []history.Entry{
		makeEntry(now.Add(-1*time.Hour), 10, 10),
	}
	opts := horizon.DefaultOptions()
	got := horizon.Evaluate(entries, now, opts)
	if got == nil {
		t.Fatal("expected prediction, got nil")
	}
	if got.ProjectedRate > 1.0 {
		t.Errorf("projected rate should be capped at 1.0, got %f", got.ProjectedRate)
	}
}

func TestWrite_NilPrediction_ShowsNoData(t *testing.T) {
	var buf bytes.Buffer
	horizon.Write(&buf, nil)
	if !bytes.Contains(buf.Bytes(), []byte("no data")) {
		t.Errorf("expected 'no data' in output, got: %s", buf.String())
	}
}

func TestWrite_WillExceed_ShowsWarning(t *testing.T) {
	var buf bytes.Buffer
	p := &horizon.Prediction{
		CurrentRate:   0.8,
		ProjectedRate: 0.95,
		Horizon:       24 * time.Hour,
		WillExceed:    true,
	}
	horizon.Write(&buf, p)
	if !bytes.Contains(buf.Bytes(), []byte("WARNING")) {
		t.Errorf("expected WARNING in output, got: %s", buf.String())
	}
}

func TestWrite_WillNotExceed_ShowsOK(t *testing.T) {
	var buf bytes.Buffer
	p := &horizon.Prediction{
		CurrentRate:   0.1,
		ProjectedRate: 0.15,
		Horizon:       24 * time.Hour,
		WillExceed:    false,
	}
	horizon.Write(&buf, p)
	if !bytes.Contains(buf.Bytes(), []byte("OK")) {
		t.Errorf("expected OK in output, got: %s", buf.String())
	}
}
