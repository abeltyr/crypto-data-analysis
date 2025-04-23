package main

import (
	"fmt"
	"os"
	"trader/utils"
)

func main() {

	fmt.Println("\n--- Starting DCA Simulation ---")
	jsonInputPath := "./data/btc/merged_eth_data.json" // Use the merged data
	// "01", "08" "22" "31"
	simStartDate := "2017-01-01"
	simEndDate := "2023-12-31"

	investmentAmount := 13.0 // Invest $100 each interval
	dailyReportOutputPath := "./reports/sum/year/7/2017/dca_report_daily.txt"
	err := utils.SimulateDCA(jsonInputPath, simStartDate, simEndDate, utils.Daily, investmentAmount, dailyReportOutputPath)
	if err != nil {
		fmt.Printf("Error running DCA simulation: %v\n", err)
		os.Exit(1)
	}

	weeklyInvestmentAmount := 92.3 // Invest $100 each interval
	weeklyReportOutputPath := "./reports/sum/year/7/2017/dca_report_weekly.txt"
	err = utils.SimulateDCA(jsonInputPath, simStartDate, simEndDate, utils.Weekly, weeklyInvestmentAmount, weeklyReportOutputPath)
	if err != nil {
		fmt.Printf("Error running DCA simulation: %v\n", err)
		os.Exit(1)
	}

	monthlyInvestmentAmount := 400.0 // Invest $100 each interval
	monthlyReportOutputPath := "./reports/sum/year/7/2017/dca_report_monthly.txt"
	errs := utils.SimulateDCA(jsonInputPath, simStartDate, simEndDate, utils.Monthly, monthlyInvestmentAmount, monthlyReportOutputPath)
	if errs != nil {
		fmt.Printf("errsor running DCA simulation: %v\n", errs)
		os.Exit(1)
	}

	fmt.Println("\n--- All Processes Completed Successfully ---")
}
