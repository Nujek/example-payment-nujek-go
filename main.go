package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	merchantapi "github.com/Nujek/sdk-payment-nujek-go"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	loadDotEnv(".env")
	if len(os.Args) < 2 {
		return errors.New("usage: go run . <create-bill|bill|list-bills|balance|create-payout|qris-static>")
	}
	client, err := newClient()
	if err != nil {
		return err
	}

	switch os.Args[1] {
	case "create-bill":
		f := flag.NewFlagSet("create-bill", flag.ContinueOnError)
		externalID := f.String("external-id", "", "merchant's order/reference ID")
		channel := f.String("channel", "NOBU_QRIS", "payment channel ID")
		total := f.String("total", "", "amount, preferably as a decimal string")
		currency := f.String("currency", "IDR", "currency")
		if err := f.Parse(os.Args[2:]); err != nil {
			return err
		}
		if *externalID == "" || *total == "" {
			return errors.New("create-bill requires -external-id and -total")
		}
		result, err := client.CreateBill(ctx, merchantapi.CreateBillRequest{ExternalID: *externalID, ChannelID: *channel, Total: *total, Currency: *currency})
		return printResult(result, err)

	case "bill":
		if len(os.Args) != 3 {
			return errors.New("usage: go run . bill <uuid>")
		}
		result, err := client.GetBill(ctx, os.Args[2])
		return printResult(result, err)

	case "list-bills":
		f := flag.NewFlagSet("list-bills", flag.ContinueOnError)
		page := f.Int("page", 1, "page number")
		perPage := f.Int("per-page", 20, "items per page")
		status := f.String("status", "", "optional bill status")
		if err := f.Parse(os.Args[2:]); err != nil {
			return err
		}
		result, err := client.ListBills(ctx, merchantapi.ListBillsQuery{Page: *page, PerPage: *perPage, Status: *status})
		return printResult(result, err)

	case "balance":
		result, err := client.GetBalance(ctx)
		return printResult(result, err)

	case "create-payout":
		f := flag.NewFlagSet("create-payout", flag.ContinueOnError)
		externalID := f.String("external-id", "", "merchant's payout/reference ID")
		bank := f.String("bank", "", "destination bank")
		account := f.String("account", "", "destination account")
		name := f.String("name", "", "destination account name")
		amount := f.String("amount", "", "amount, preferably as a decimal string")
		currency := f.String("currency", "IDR", "currency")
		if err := f.Parse(os.Args[2:]); err != nil {
			return err
		}
		if *externalID == "" || *bank == "" || *account == "" || *name == "" || *amount == "" {
			return errors.New("create-payout requires -external-id, -bank, -account, -name, and -amount")
		}
		result, err := client.CreatePayout(ctx, merchantapi.CreatePayoutRequest{ExternalID: *externalID, DestinationBank: *bank, DestinationAccount: *account, DestinationName: *name, Amount: *amount, Currency: *currency})
		return printResult(result, err)

	case "qris-static":
		if len(os.Args) == 2 {
			result, err := client.ListQRISStatic(ctx)
			return printResult(result, err)
		}
		result, err := client.GetQRISStatic(ctx, os.Args[2])
		return printResult(result, err)
	default:
		return fmt.Errorf("unknown command %q", os.Args[1])
	}
}

// loadDotEnv loads simple KEY=VALUE entries without overwriting environment
// variables already supplied by the shell or deployment platform.
func loadDotEnv(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), "\"'")
		if key != "" {
			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, value)
			}
		}
	}
}

func newClient() (*merchantapi.Client, error) {
	baseURL := os.Getenv("PAYMENT_BASE_URL")
	if baseURL == "" {
		return nil, errors.New("PAYMENT_BASE_URL is required")
	}
	apiKey := os.Getenv("PAYMENT_API_KEY")
	apiSecret := os.Getenv("PAYMENT_API_SECRET")
	if apiKey == "" || apiKey == "your-api-key" {
		return nil, errors.New("PAYMENT_API_KEY belum diisi; gunakan API key merchant yang valid")
	}
	if apiSecret == "" || apiSecret == "your-api-secret" {
		return nil, errors.New("PAYMENT_API_SECRET belum diisi; gunakan API secret merchant yang valid")
	}
	return merchantapi.NewClient(baseURL, apiKey, apiSecret)
}

func printResult(value any, err error) error {
	if err != nil {
		return err
	}
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(encoded))
	return nil
}
