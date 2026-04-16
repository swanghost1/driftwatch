// Package batch splits a slice of drift results into fixed-size chunks.
//
// This is useful when downstream consumers — such as notification webhooks
// or audit endpoints — impose limits on the number of records that can be
// sent in a single request.
//
// Basic usage:
//
//	batches := batch.Apply(results, batch.Options{Size: 25})
//	for _, b := range batches {
//		sendToWebhook(b)
//	}
//
// A Size of zero or less disables batching and returns all results in a
// single batch.
package batch
