package services

import (
	"math"
	"sort"
	"sync"
)

var (
	packSizes   = []int{5000, 2000, 1000, 500, 250} // necessarily in descending order
	packSizesMu sync.RWMutex
)

type PackResult struct {
	Pack     int `json:"pack"`
	Quantity int `json:"quantity"`
}

type PackService struct{}

func (s *PackService) GetPackSizes() []int {
	packSizesMu.RLock()
	defer packSizesMu.RUnlock()
	return append([]int{}, packSizes...)
}

func (s *PackService) SetPackSizes(newSizes []int) {
	packSizesMu.Lock()
	defer packSizesMu.Unlock()

	validSizes := []int{}
	for _, size := range newSizes {
		if size > 0 {
			validSizes = append(validSizes, size)
		}
	}

	if len(validSizes) > 0 {
		packSizes = validSizes
		sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))
	}
}

func (s *PackService) CalculatePacks(itemsNo int) []PackResult {
	sizes := s.GetPackSizes()
	if len(sizes) == 0 {
		return nil
	}

	// We'll search up to itemsNo + largest pack size to find the minimal solution
	maxSearch := itemsNo + sizes[0]
	minItems := math.MaxInt32
	minPacks := math.MaxInt32
	var bestSolution []PackResult

	// DP table where dp[i] represents the best solution for i items
	dp := make([]struct {
		totalItems int          // total items in solution
		totalPacks int          // total number of packs
		packs      []PackResult // detailed pack breakdown
	}, maxSearch+1)

	// Initialize DP table
	for i := 1; i <= maxSearch; i++ {
		dp[i].totalItems = math.MaxInt32
		dp[i].totalPacks = math.MaxInt32
	}

	// Base case: 0 items requires 0 packs
	dp[0].totalItems = 0
	dp[0].totalPacks = 0
	dp[0].packs = []PackResult{}

	for i := 1; i <= maxSearch; i++ {
		for _, size := range sizes {
			if size <= i {
				candidateItems := dp[i-size].totalItems + size
				candidatePacks := dp[i-size].totalPacks + 1

				// Check if this candidate is better than current solution
				if candidateItems < dp[i].totalItems ||
					(candidateItems == dp[i].totalItems && candidatePacks < dp[i].totalPacks) {
					dp[i].totalItems = candidateItems
					dp[i].totalPacks = candidatePacks

					// Copy the previous packs
					dp[i].packs = make([]PackResult, len(dp[i-size].packs))
					copy(dp[i].packs, dp[i-size].packs)

					// Add current pack
					found := false
					for j := range dp[i].packs {
						if dp[i].packs[j].Pack == size {
							dp[i].packs[j].Quantity++
							found = true
							break
						}
					}
					if !found {
						dp[i].packs = append(dp[i].packs, PackResult{Pack: size, Quantity: 1})
					}
				}
			}
		}

		// Check if this is a potential solution (i >= itemsNo)
		if i >= itemsNo {
			// Rule 2: Prefer solution with least items
			if dp[i].totalItems < minItems ||
				// Rule 3: If same items, prefer fewer packs
				(dp[i].totalItems == minItems && dp[i].totalPacks < minPacks) {
				minItems = dp[i].totalItems
				minPacks = dp[i].totalPacks
				bestSolution = dp[i].packs
			}
		}
	}

	return bestSolution
}
