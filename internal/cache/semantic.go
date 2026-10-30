package cache

import "context"

type SemanticMatcher struct {
	threshold float64
}

func NewSemanticMatcher(threshold float64) *SemanticMatcher {
	if threshold <= 0 {
		threshold = 0.92
	}
	return &SemanticMatcher{threshold: threshold}
}

func (m *SemanticMatcher) Match(ctx context.Context, query []float32, candidates map[string][]float32) (key string, score float64, ok bool) {
	bestKey := ""
	best := 0.0
	for k, emb := range candidates {
		s := CosineSimilarity(query, emb)
		if s > best {
			best = s
			bestKey = k
		}
	}
	if best >= m.threshold {
		return bestKey, best, true
	}
	return "", best, false
}
