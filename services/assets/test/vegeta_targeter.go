package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	vegeta "github.com/tsenart/vegeta/lib"
)

// Config holds the configuration for the load test.
type Config struct {
	Duration time.Duration
	RateFreq int
	RatePer  time.Duration
}

// TargetConfig represents the configuration for a single target.
type TargetConfig struct {
	URLTemplate string
	Method      string
	Header      map[string][]string
	Body        []byte
	ReportPath  string
}

// newDynamicTargeter creates a vegeta.Targeter for a single target configuration.
func newDynamicTargeter(tc TargetConfig) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		// Generate a new UUID for dynamic URL parameter.
		uuid := uuid.New()
		url := fmt.Sprintf(tc.URLTemplate, uuid.String())

		tgt.Method = tc.Method
		tgt.URL = url
		tgt.Header = tc.Header
		tgt.Body = tc.Body

		return nil
	}
}

// executeAttack runs the attack using the provided configuration, targeter, and report path.
func executeAttack(config Config, targeter vegeta.Targeter, reportPath string) error {
	attacker := vegeta.NewAttacker()
	rate := vegeta.Rate{Freq: config.RateFreq, Per: config.RatePer}
	metrics := &vegeta.Metrics{}

	// Start the attack and collect metrics.
	for res := range attacker.Attack(targeter, rate, config.Duration, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	// Generate and save the report.
	reporter := vegeta.NewHDRHistogramPlotReporter(metrics)
	f, err := os.Create(reportPath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer f.Close()

	if err = reporter.Report(f); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	log.Printf("Report saved to %s", f.Name())
	return nil
}

func main() {
	// Configure the list of target URLs and their settings.
	targetConfigs := []TargetConfig{
		{
			URLTemplate: "http://localhost:6000/api/v1/assets/%s",
			Method:      "POST",
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body: []byte(`{
				"assetMoneyAmount": 1000,
				"assetMoneyCurrency": "USD",
				"assetName": "My Asset",
				"assetType": "cash"
			}`),
			ReportPath: "./reports/vegeta-report-create-asset.hgrm",
		},
		// Additional targets can be added here.
	}

	// Set up the load test configuration.
	config := Config{
		Duration: 5 * time.Second,
		RateFreq: 100,
		RatePer:  time.Second,
	}

	// For each target configuration, execute the attack.
	for _, tc := range targetConfigs {
		// Initialize the targeter with the target configuration.
		targeter := newDynamicTargeter(tc)

		// Execute the load test attack for this target.
		if err := executeAttack(config, targeter, tc.ReportPath); err != nil {
			log.Fatalf("Error executing attack for target %s: %v", tc.URLTemplate, err)
		}
	}
}
