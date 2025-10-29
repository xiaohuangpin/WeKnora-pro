package metric

import (
	"github.com/Tencent/WeKnora/internal/types"
)

// RecallMetric calculates recall for retrieval evaluation
type RecallMetric struct{}

// NewRecallMetric creates a new RecallMetric instance
func NewRecallMetric() *RecallMetric {
	return &RecallMetric{}
}

// Compute calculates the recall score
func (r *RecallMetric) Compute(metricInput *types.MetricInput) float64 {
	// Get ground truth and predicted IDs
	gts := metricInput.RetrievalGT
	ids := metricInput.RetrievalIDs

	// Convert ground truth to sets for efficient lookup
	gtSets := SliceMap(gts, ToSet)
	// Count total hits across all relevant documents
	ahit := Fold(gtSets, 0, func(a int, b map[int]struct{}) int { return a + Hit(ids, b) })

	// Handle case with no ground truth
	if len(gtSets) == 0 {
		return 0.0
	}

	// Recall = total hits / total relevant documents
	return float64(ahit) / float64(len(gtSets))
}
