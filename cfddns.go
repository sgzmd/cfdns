package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

func main() {
	var apiToken, zoneName, recordName, recordContent string

	flag.StringVar(&apiToken, "token", "", "Cloudflare API token")
	flag.StringVar(&zoneName, "zone", "", "Zone name (e.g., example.com)")
	flag.StringVar(&recordName, "record", "", "Record name to update (e.g., sub.example.com)")
	flag.StringVar(&recordContent, "content", "", "New content of the record (e.g., IP address)")
	flag.Parse()

	// Validate input
	if apiToken == "" || zoneName == "" || recordName == "" || recordContent == "" {
		log.Printf("token=%s, zone=%s, record=%s, content=%s", apiToken, zoneName, recordName, recordContent)
		fmt.Println("Missing required parameters.")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		fmt.Println("All parameters are required.")
		flag.Usage()
		os.Exit(1)
	}

	// Create a new API object
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.Println("Error creating Cloudflare API object:", err)
		os.Exit(1)
	}

	// Fetch the zone ID
	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		log.Println("Error fetching zone ID:", err)
		os.Exit(1)
	}

	// Fetch the DNS record ID
	zi := cloudflare.ZoneIdentifier(zoneID)

	params := cloudflare.ListDNSRecordsParams{Name: recordName}
	records, _, err := api.ListDNSRecords(context.Background(), zi, params)
	if err != nil {
		log.Println("Error fetching DNS records:", err)
		os.Exit(1)
	}
	if len(records) == 0 {
		log.Println("No DNS records found for", recordName)
		os.Exit(1)
	}

	for _, record := range records {
		log.Printf("Found DNS record: %+v\n", record)
	}

	// Update the DNS record
	record := records[0]
	record.Content = recordContent
	upd := cloudflare.UpdateDNSRecordParams{
		Type:    record.Type,
		Name:    record.Name,
		ID:      record.ID,
		Content: recordContent,
		TTL:     record.TTL,
	}
	log.Printf(
		"Updating DNS record %+v", upd)
	rec, err := api.UpdateDNSRecord(context.Background(), zi, upd)

	if err != nil {
		log.Println("Error updating DNS record:", err)
		os.Exit(1)
	}

	j, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		log.Println("Error marshalling DNS record:", err)
		os.Exit(1)
	}
	log.Printf("DNS record updated successfully: \n %s\n", j)
}
