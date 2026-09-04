package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CarbonResponse struct {
	Zone            string  `json:"zone"`
	CarbonIntensity float64 `json:"carbonIntensity"`
	Datetime        string  `json:"datetime"`
}

func main() {
	apiKey := os.Getenv("ELECTRICITYMAPS_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: set ELECTRICITYMAPS_API_KEY environment variable")
		os.Exit(1)
	}

	url := "https://api.electricitymap.org/v3/carbon-intensity/latest?zone=IT"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("auth-token", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result CarbonResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing response:", err)
		os.Exit(1)
	}

	fmt.Printf("Zone: %s\nCarbon Intensity: %.2f gCO2eq/kWh\nDatetime: %s\n",
		result.Zone, result.CarbonIntensity, result.Datetime)

	const carbonThreshold = 250.0 // gCO2eq/kWh — based on Italy's typical daily low (~210) vs high (~320)

	if result.CarbonIntensity <= carbonThreshold {
		fmt.Println(" Carbon intensity is LOW — safe to run workload now")
	} else {
		fmt.Println(" Carbon intensity is HIGH — workload should wait")
	}
}
