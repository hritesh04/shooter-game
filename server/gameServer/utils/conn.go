package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func RegisterServer(wg *sync.WaitGroup, address, serverType string) {
	defer wg.Done()
	server := struct {
		Address string `json:"address"`
		Type    string `json:"type"`
	}{
		Address: address,
		Type:    serverType,
	}
	out, err := json.Marshal(server)
	if err != nil {
		log.Fatalf("error binding request payload : %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, os.Getenv("ROOT_SERVER_URL")+"/registerServer", bytes.NewBuffer(out))
	if err != nil {
		log.Printf("error creating request : %v", err)
	}
	tlsVerification := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tlsVerification}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("error making http call : %v", err)
	}
	if res != nil {
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				log.Printf("error parsing response body :%v", err)
				return
			}
			log.Printf("error in registrating server :%s", string(body))
			return
		}
		log.Printf("Successfully registered %s server with root server", serverType)
	}
}
