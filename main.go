package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func unmarshalSamples(data []byte) (samples, error) {
	var m samples
	err := json.Unmarshal(data, &m)
	return m, err
}

type samples struct {
	Malware []malware `json:"malware"`
}

type malware struct {
	Name       string  `json:"name"`
	Type       *string `json:"type,omitempty"`
	VirusTotal string  `json:"virusTotal"`
	MoreInfo   string  `json:"moreInfo"`
	Download   string  `json:"download"`
}

var outputDir string

const (
	defaultDir = "./malware/"
	jsonURL    = "https://objective-see.com/malware.json"
)

func main() {
	flag.StringVar(&outputDir, "outputdir", defaultDir, "output directory for the malware samples")

	flag.Parse()

	if outputDir[len(outputDir)-1] != '/' {
		outputDir += "/"
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, 0644)
		if err != nil {
			panic(err)
		}
	}

	buf, err := getObjectiveSeeSamples()
	if err != nil {
		panic(err)
	}

	samples, err := unmarshalSamples(buf)
	if err != nil {
		panic(err)
	}

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	for _, s := range samples.Malware {
		req, err := http.NewRequest("GET", s.Download, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0")

		resp, err := netClient.Do(req)
		if err != nil {
			fmt.Printf("failed to get sample: %s\n", err)
			continue
		}

		cutName := cleanString(s.Name)

		fmt.Println(cutName)

		out, err := os.Create(outputDir + "/" + cutName)
		if err != nil {
			fmt.Printf("failed to create file: %s\n", err)
			continue
		}

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Printf("failed to write file: %s\n", err)
			continue
		}
		out.Close()
	}

}

func getObjectiveSeeSamples() ([]byte, error) {
	resp, err := http.Get(jsonURL)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(resp.Body)

	return buf, err
}

func cleanString(s string) string {
	b := &strings.Builder{}
	for _, c := range s {
		if c == ' ' {
			break
		}
		b.WriteRune(c)
	}
	return b.String()
}
