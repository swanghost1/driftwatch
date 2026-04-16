package stream_test

import (
	"testing"

	"github.com/example/driftwatch/internal/drift"
	"github.com/example/driftwatch/internal/stream"
)

func TestFanout_SingleSubscriber_ReceivesAll(t *testing.T) {
	fo := stream.NewFanout(10)
	sub := fo.Subscribe()

	src := make(chan drift.Result, 3)
	src <- drift.Result{Service: "a"}
	src <- drift.Result{Service: "b"}
	src <- drift.Result{Service: "c"}
	close(src)

	fo.Run(src)

	var got []string
	for r := range sub {
		got = append(got, r.Service)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}
}

func TestFanout_MultipleSubscribers_EachReceiveAll(t *testing.T) {
	fo := stream.NewFanout(10)
	sub1 := fo.Subscribe()
	sub2 := fo.Subscribe()

	src := make(chan drift.Result, 2)
	src <- drift.Result{Service: "x"}
	src <- drift.Result{Service: "y"}
	close(src)

	fo.Run(src)

	for _, sub := range []<-chan drift.Result{sub1, sub2} {
		count := 0
		for range sub {
			count++
		}
		if count != 2 {
			t.Errorf("expected 2 results per subscriber, got %d", count)
		}
	}
}

func TestFanout_Close_ClosesAllSubscribers(t *testing.T) {
	fo := stream.NewFanout(5)
	sub := fo.Subscribe()
	fo.Close()
	_, open := <-sub
	if open {
		t.Error("expected subscriber channel to be closed")
	}
}

func TestFanout_Publish_AfterClose_DoesNotPanic(t *testing.T) {
	fo := stream.NewFanout(5)
	fo.Close()
	// Should not panic; no subscribers remain.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	fo.Publish(drift.Result{Service: "z"})
}

func TestFanout_NoSubscribers_PublishIsNoop(t *testing.T) {
	fo := stream.NewFanout(5)
	// Should not block or panic with no subscribers.
	fo.Publish(drift.Result{Service: "solo"})
	fo.Close()
}
