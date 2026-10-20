// Validate an e-invoice (Factur-X / ZUGFeRD PDF, XRechnung, Peppol BIS or other EN 16931 XML).
// Go standard library only.
// Usage: INVOIR_API_KEY=inv_... go run ./go/validate [file]   (from the repository root)
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type report struct {
	Status   string `json:"status"`
	Format   string `json:"format"`
	Errors   int    `json:"errors"`
	Warnings int    `json:"warnings"`
	Error    string `json:"error"`
	Message  string `json:"message"`
	Checks   []struct {
		Findings []struct {
			Severity string `json:"severity"`
			Rule     string `json:"rule"`
			Message  string `json:"message"`
			Location string `json:"location"`
		} `json:"findings"`
	} `json:"checks"`
}

func main() {
	key := os.Getenv("INVOIR_API_KEY")
	if key == "" {
		log.Fatal("Set INVOIR_API_KEY (get a free key at https://validate.invoir.dev/signup)")
	}
	path := "RE-2026-0042.pdf"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, _ := w.CreateFormFile("file", filepath.Base(path))
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(part, f)
	f.Close()
	w.Close()

	req, _ := http.NewRequest("POST", "https://validate.invoir.dev/api/v1/validate", &body)
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{Timeout: 60 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var r report
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatal(err)
	}
	if res.StatusCode >= 400 {
		log.Fatal(r.Error + ": " + r.Message)
	}
	fmt.Println(r.Status, r.Format, r.Errors, "errors", r.Warnings, "warnings")
	for _, c := range r.Checks {
		for _, f := range c.Findings {
			fmt.Println(" ", f.Severity, f.Rule, f.Message, f.Location)
		}
	}
}
