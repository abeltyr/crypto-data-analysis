package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Investment represents a single investment action during the simulation.
type Investment struct {
	Date       time.Time `json:"date"`
	Price      float64   `json:"price"`
	AmountUSD  float64   `json:"amountUSD"`
	AmountCoin float64   `json:"amountCoin"`
}

// SimulationResult holds the summary of the DCA simulation.
type SimulationResult struct {
	StartDate        string       `json:"startDate"`
	EndDate          string       `json:"endDate"`
	Interval         string       `json:"interval"`
	InvestmentAmount float64      `json:"investmentAmount"`
	TotalInvestedUSD float64      `json:"totalInvestedUSD"`
	TotalCoinBought  float64      `json:"totalCoinBought"`
	FinalCoinPrice   float64      `json:"finalCoinPrice"`
	FinalValueUSD    float64      `json:"finalValueUSD"`
	GainLossUSD      float64      `json:"gainLossUSD"`
	GainLossPercent  float64      `json:"gainLossPercent"`
	Investments      []Investment `json:"investments"`
}

// Interval type for clarity
type Interval string

const (
	Daily   Interval = "daily"
	Weekly  Interval = "weekly"
	Monthly Interval = "monthly"
)

// SimulateDCA performs a Dollar-Cost Averaging simulation.
func SimulateDCA(jsonPath, startDateStr, endDateStr string, interval Interval, amount float64, reportPath string) error {
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return fmt.Errorf("error parsing start date '%s': %v", startDateStr, err)
	}
	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return fmt.Errorf("error parsing end date '%s': %v", endDateStr, err)
	}

	if endDate.Before(startDate) {
		return fmt.Errorf("end date must be on or after start date")
	}

	// Read the JSON data file
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("error reading JSON file '%s': %v", jsonPath, err)
	}

	// Parse the JSON data: map[string][]*Kline (date -> hourly klines)
	var klinesData map[string][]*Kline
	err = json.Unmarshal(jsonData, &klinesData)
	if err != nil {
		return fmt.Errorf("error parsing JSON data from '%s': %v", jsonPath, err)
	}

	// Generate investment dates based on interval
	investmentDates := generateInvestmentDates(startDate, endDate, interval)
	if len(investmentDates) == 0 {
		return fmt.Errorf("no investment dates generated for the given range and interval")
	}

	var investments []Investment
	totalInvestedUSD := 0.0
	totalCoinBought := 0.0

	fmt.Println("Starting DCA Simulation...")
	// Simulate investments
	for _, investDate := range investmentDates {
		dateStr := investDate.Format(layout)
		hourlyKlines, dateExists := klinesData[dateStr]

		if !dateExists || len(hourlyKlines) == 0 {
			fmt.Printf("Warning: No data found for investment date %s. Skipping investment.\n", dateStr)
			continue
		}

		// Find the first available non-nil Kline for the day to get the price
		var priceStr string
		foundPrice := false
		for _, kline := range hourlyKlines {
			if kline != nil && kline.Close != "" {
				priceStr = kline.Close // Use the closing price
				foundPrice = true
				break
			}
		}

		if !foundPrice {
			fmt.Printf("Warning: No valid price found for investment date %s. Skipping investment.\n", dateStr)
			continue
		}

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			fmt.Printf("Warning: Could not parse price '%s' for date %s. Skipping investment. Error: %v\n", priceStr, dateStr, err)
			continue
		}

		if price <= 0 {
			fmt.Printf("Warning: Invalid price %f found for date %s. Skipping investment.\n", price, dateStr)
			continue
		}

		coinBought := amount / price
		totalInvestedUSD += amount
		totalCoinBought += coinBought

		investments = append(investments, Investment{
			Date:       investDate,
			Price:      price,
			AmountUSD:  amount,
			AmountCoin: coinBought,
		})
		fmt.Printf("Invested %.2f USD on %s at price %.4f, bought %.8f coin\n", amount, dateStr, price, coinBought)
	}

	if len(investments) == 0 {
		return fmt.Errorf("no investments were made during the simulation period")
	}

	// Find the final price on the end date or the last available date
	finalPrice, lastDateWithPrice := findFinalPrice(endDate, startDate, klinesData)
	if finalPrice <= 0 {
		return fmt.Errorf("could not determine a valid final price for the portfolio")
	}

	finalValueUSD := totalCoinBought * finalPrice
	gainLossUSD := finalValueUSD - totalInvestedUSD
	gainLossPercent := 0.0
	if totalInvestedUSD > 0 {
		gainLossPercent = (gainLossUSD / totalInvestedUSD) * 100
	}

	// Prepare result
	result := SimulationResult{
		StartDate:        startDateStr,
		EndDate:          endDateStr,
		Interval:         string(interval),
		InvestmentAmount: amount,
		TotalInvestedUSD: totalInvestedUSD,
		TotalCoinBought:  totalCoinBought,
		FinalCoinPrice:   finalPrice,
		FinalValueUSD:    finalValueUSD,
		GainLossUSD:      gainLossUSD,
		GainLossPercent:  gainLossPercent,
		Investments:      investments,
	}

	// Generate report content
	reportContent := generateReport(result, lastDateWithPrice)

	// Ensure report directory exists
	reportDir := filepath.Dir(reportPath)
	if err := os.MkdirAll(reportDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating report directory '%s': %v", reportDir, err)
	}

	// Write report to file
	err = os.WriteFile(reportPath, []byte(reportContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing report to file '%s': %v", reportPath, err)
	}

	fmt.Printf("\nDCA Simulation Complete. Text report saved to %s\n", reportPath)

	result.Investments = []Investment{} // Clear the Investments slice for JSON output

	// Save results as JSON
	jsonReportPath := strings.TrimSuffix(reportPath, filepath.Ext(reportPath)) + ".json"
	jsonDataBytes, err := json.MarshalIndent(result, "", "  ") // Use MarshalIndent for pretty printing
	if err != nil {
		return fmt.Errorf("error marshaling simulation result to JSON: %v", err)
	}

	err = os.WriteFile(jsonReportPath, jsonDataBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON report to file '%s': %v", jsonReportPath, err)
	}

	fmt.Printf("JSON report saved to %s\n", jsonReportPath)

	return nil
}

// generateInvestmentDates creates a list of dates based on the interval.
func generateInvestmentDates(start, end time.Time, interval Interval) []time.Time {
	var dates []time.Time
	current := start

	for !current.After(end) {
		dates = append(dates, current)
		switch interval {
		case Daily:
			current = current.AddDate(0, 0, 1)
		case Weekly:
			current = current.AddDate(0, 0, 7)
		case Monthly:
			current = current.AddDate(0, 1, 0)
		default:
			// Should not happen if interval is validated
			return dates // Return what we have if interval is unknown
		}
	}
	return dates
}

// findFinalPrice attempts to find the closing price on the end date, or the latest available date before it.
func findFinalPrice(endDate, startDate time.Time, klinesData map[string][]*Kline) (float64, string) {
	layout := "2006-01-02"
	// Sort the dates available in the data
	var availableDates []string
	for dateStr := range klinesData {
		availableDates = append(availableDates, dateStr)
	}
	sort.Strings(availableDates)

	// Iterate backwards from the end date within the simulation range
	currentDate := endDate
	for !currentDate.Before(startDate) {
		dateStr := currentDate.Format(layout)
		if hourlyKlines, exists := klinesData[dateStr]; exists && len(hourlyKlines) > 0 {
			// Iterate backwards through the hours of the day to find the last price
			for i := len(hourlyKlines) - 1; i >= 0; i-- {
				if hourlyKlines[i] != nil && hourlyKlines[i].Close != "" {
					price, err := strconv.ParseFloat(hourlyKlines[i].Close, 64)
					if err == nil && price > 0 {
						fmt.Printf("Using final price %.4f from %s (Hour %d)\n", price, dateStr, i)
						return price, dateStr
					}
				}
			}
		}
		// Move to the previous day
		currentDate = currentDate.AddDate(0, 0, -1)
	}

	fmt.Println("Warning: Could not find final price within the specified date range. Trying last available date in data.")
	// If no price found within range, try the very last date in the dataset
	if len(availableDates) > 0 {
		lastDateStr := availableDates[len(availableDates)-1]
		if hourlyKlines, exists := klinesData[lastDateStr]; exists && len(hourlyKlines) > 0 {
			for i := len(hourlyKlines) - 1; i >= 0; i-- {
				if hourlyKlines[i] != nil && hourlyKlines[i].Close != "" {
					price, err := strconv.ParseFloat(hourlyKlines[i].Close, 64)
					if err == nil && price > 0 {
						fmt.Printf("Using final price %.4f from last available date %s (Hour %d)\n", price, lastDateStr, i)
						return price, lastDateStr
					}
				}
			}
		}
	}

	return 0.0, "N/A" // Indicate failure
}

// generateReport creates a formatted string summary of the simulation.
func generateReport(result SimulationResult, lastPriceDate string) string {
	var sb strings.Builder

	sb.WriteString("--- Dollar-Cost Averaging (DCA) Simulation Report ---\n\n")
	sb.WriteString(fmt.Sprintf("Simulation Period: %s to %s\n", result.StartDate, result.EndDate))
	sb.WriteString(fmt.Sprintf("Investment Interval: %s\n", result.Interval))
	sb.WriteString(fmt.Sprintf("Investment Amount per Interval: $%.2f\n", result.InvestmentAmount))
	sb.WriteString("\n--- Summary ---\n")
	sb.WriteString(fmt.Sprintf("Total Amount Invested: $%.2f\n", result.TotalInvestedUSD))
	sb.WriteString(fmt.Sprintf("Total Coin Acquired: %.8f\n", result.TotalCoinBought))
	sb.WriteString(fmt.Sprintf("Final Coin Price (as of %s): $%.4f\n", lastPriceDate, result.FinalCoinPrice))
	sb.WriteString(fmt.Sprintf("Final Portfolio Value: $%.2f\n", result.FinalValueUSD))
	sb.WriteString(fmt.Sprintf("Total Gain/Loss: $%.2f\n", result.GainLossUSD))
	sb.WriteString(fmt.Sprintf("Total Gain/Loss Percentage: %.2f%%\n", result.GainLossPercent))

	sb.WriteString("\n--- Investment Details ---\n")
	sb.WriteString("Date       | Price      | Amount (USD) | Amount (Coin)\n")
	sb.WriteString("-----------|------------|--------------|--------------\n")
	layout := "2006-01-02"
	for _, inv := range result.Investments {
		sb.WriteString(fmt.Sprintf("%s | %10.4f | %12.2f | %13.8f\n",
			inv.Date.Format(layout),
			inv.Price,
			inv.AmountUSD,
			inv.AmountCoin))
	}

	return sb.String()
}
