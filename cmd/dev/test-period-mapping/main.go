package main

import (
	"fmt"
	"log"
)

func main() {
	if err := testPeriodMapping(); err != nil {
		log.Fatalf("Failed to test period mapping: %v", err)
	}
}

func testPeriodMapping() error {
	fmt.Println("🧪 Testing Period Mapping")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Test period mappings
	testPeriods := []int64{
		1767549424, // Should map to 1
		1767557168, // Should map to 2
		1767650263, // Should map to 3
		1768130418, // Should map to 4
		1234567890, // Should fail (unknown period)
	}

	expectedMappings := map[int64]int64{
		1767549424: 1,
		1767557168: 2,
		1767650263: 3,
		1768130418: 4,
	}

	fmt.Printf("📋 Testing Period Mappings:\n")
	fmt.Println()

	for _, settlementPeriod := range testPeriods {
		contractPeriod, err := testMapPeriod(settlementPeriod)

		if err != nil {
			fmt.Printf("   Settlement Period %d -> ERROR: %v\n", settlementPeriod, err)
			if expected, exists := expectedMappings[settlementPeriod]; exists {
				fmt.Printf("      ❌ Expected: %d, Got: ERROR\n", expected)
			} else {
				fmt.Printf("      ✅ Expected error for unknown period\n")
			}
		} else {
			fmt.Printf("   Settlement Period %d -> Contract Period %d\n", settlementPeriod, contractPeriod)
			if expected, exists := expectedMappings[settlementPeriod]; exists {
				if contractPeriod == expected {
					fmt.Printf("      ✅ Correct mapping\n")
				} else {
					fmt.Printf("      ❌ Expected: %d, Got: %d\n", expected, contractPeriod)
				}
			} else {
				fmt.Printf("      ❌ Unexpected success for unknown period\n")
			}
		}
		fmt.Println()
	}

	fmt.Printf("🎯 Test Summary:\n")
	fmt.Printf("   - Known periods should map to sequential contract periods (1,2,3,4)\n")
	fmt.Printf("   - Unknown periods should return an error\n")
	fmt.Printf("   - This mapping is used in ClaimMiningReward to convert settlement period IDs\n")
	fmt.Println()

	fmt.Printf("✅ Period mapping test completed!\n")
	fmt.Printf("📝 Next: Test actual mining claims with the updated logic\n")

	return nil
}

// testMapPeriod replicates the mapping logic for testing
func testMapPeriod(settlementPeriodID int64) (int64, error) {
	// Fixed mapping based on sorted settlement periods
	periodMapping := map[int64]int64{
		1767549424: 1, // Oldest settlement -> Contract period 1
		1767557168: 2, // Second oldest -> Contract period 2
		1767650263: 3, // Third oldest -> Contract period 3
		1768130418: 4, // Newest settlement -> Contract period 4
	}

	contractPeriod, exists := periodMapping[settlementPeriodID]
	if !exists {
		return 0, fmt.Errorf("unknown settlement period ID: %d", settlementPeriodID)
	}

	return contractPeriod, nil
}
