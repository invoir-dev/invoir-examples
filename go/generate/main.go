// Generate an e-invoice from invoice.json and save the PDF and XML. Go standard library only.
// Usage: INVOIR_API_KEY=inv_... go run ./go/generate [invoice.json]   (from the repository root)
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type file struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

type document struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
	PDF     *file  `json:"pdf"`
	XML     *file  `json:"xml"`
	Totals  struct {
		Due      float64 `json:"due"`
		Currency string  `json:"currency"`
	} `json:"totals"`
	Errors []struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	} `json:"errors"`
	Validation *struct {
		Checks []struct {
			Findings []struct {
				Severity string `json:"severity"`
				Rule     string `json:"rule"`
				Message  string `json:"message"`
				Location string `json:"location"`
			} `json:"findings"`
		} `json:"checks"`
	} `json:"validation"`
}

func main() {
	key := os.Getenv("INVOIR_API_KEY")
	if key == "" {
		log.Fatal("Set INVOIR_API_KEY (get a free key at https://validate.invoir.dev/signup)")
	}
	path := "invoice.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	body, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var invoice struct {
		Number string `json:"number"`
	}
	if err := json.Unmarshal(body, &invoice); err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	req, _ := http.NewRequest("POST", "https://validate.invoir.dev/api/v1/invoices", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	// Idempotency-Key: a retry with the same key and body returns the first answer instead of a second invoice.
	req.Header.Set("Idempotency-Key", invoice.Number)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var doc document
	if err := json.NewDecoder(res.Body).Decode(&doc); err != nil {
		log.Fatal(err)
	}
	if res.StatusCode >= 400 {
		for _, e := range doc.Errors { // 400: what is wrong with the request
			fmt.Println(" ", e.Field, e.Message)
		}
		if doc.Validation != nil { // 422: what the official rules rejected
			for _, c := range doc.Validation.Checks {
				for _, f := range c.Findings {
					fmt.Println(" ", f.Severity, f.Rule, f.Message, f.Location)
				}
			}
		}
		log.Fatal(doc.Error + ": " + doc.Message)
	}
	fmt.Println(doc.ID, doc.Status, doc.Totals.Due, doc.Totals.Currency)

	for _, f := range []*file{doc.PDF, doc.XML} {
		if f == nil {
			continue
		}
		if err := download(client, key, f); err != nil {
			log.Fatal(err)
		}
		fmt.Println("saved", f.Filename)
	}
}

func download(client *http.Client, key string, f *file) error {
	req, _ := http.NewRequest("GET", f.URL, nil)
	req.Header.Set("Authorization", "Bearer "+key)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	out, err := os.Create(f.Filename)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, res.Body)
	return err
}
