package services

import (
	"fmt"
	"path/filepath"
	"time"
	"trader/utils"
)

func SimulateWeekly() {

	fmt.Println("\n--- Starting Weekly DCA Simulations for a Year ---")
	jsonInputPath := "./data/btc/merged_eth_data.json"                       // Use the merged data
	investmentAmount := 50.0                                                 // Invest $50 each interval (daily)
	yearsToSimulate := []int{2018, 2019, 2020, 2021, 2022, 2023, 2024, 2025} // Define the years to simulate

	fmt.Println("\n--- Starting Weekly DCA Simulations for Years:", yearsToSimulate, "---")

	for _, yearToSimulate := range yearsToSimulate {
		fmt.Printf("\n=== Simulating Year: %d ===\n", yearToSimulate)
		for month := 1; month <= 12; month++ {
			fmt.Printf("\n--- Simulating Month: %d ---\n", month)
			monthStart := time.Date(yearToSimulate, time.Month(month), 1, 6, 0, 0, 0, time.UTC)
			monthEnd := monthStart.AddDate(0, 1, -1) // Last day of the month

			currentDate := monthStart
			weekNumber := 1

			for currentDate.Before(monthEnd) || currentDate.Equal(monthEnd) {
				// Calculate week start (Monday) and end (Sunday)
				weekStart := currentDate
				// Find the end of the current week (Sunday) or month end, whichever comes first
				weekEnd := weekStart.AddDate(0, 0, 6) // Tentative end of week
				if weekEnd.After(monthEnd) {
					weekEnd = monthEnd
				}

				simStartDateStr := weekStart.Format("2006-01-02")
				simEndDateStr := weekEnd.Format("2006-01-02")

				// Construct report path
				reportDir := fmt.Sprintf("./reports/weekly")
				reportFileName := fmt.Sprintf("dca_report_week%d_%02d_%d.txt", weekNumber, month, yearToSimulate)
				dailyReportOutputPath := filepath.Join(reportDir, reportFileName)

				fmt.Printf("Running simulation for Week %d (%s to %s)\n", weekNumber, simStartDateStr, simEndDateStr)

				err := utils.SimulateDCA(jsonInputPath, simStartDateStr, simEndDateStr, utils.Daily, investmentAmount, dailyReportOutputPath)
				if err != nil {
					// Log the error but continue with the next week/month
					fmt.Printf("Error running DCA simulation for %s to %s: %v\n", simStartDateStr, simEndDateStr, err)
				} else {
					fmt.Printf("Successfully completed simulation for Week %d. Report: %s\n", weekNumber, dailyReportOutputPath)
				}

				// Move to the next week
				currentDate = weekEnd.AddDate(0, 0, 1)
				weekNumber++
			}
		}
	}
	fmt.Println("\n--- All Weekly DCA Simulations Completed ---")
}
