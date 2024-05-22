package cli

import (
	"log"
	"net/http"
	"sync"

	"github.com/Jeffail/gabs"
)

type RequestBody struct {
	SourceLang string
	TargetLang string
	SourceText string
}

const translateUrl = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(body RequestBody, str chan string, wg *sync.WaitGroup) {
	var client = &http.Client{}
	var req, err1 = http.NewRequest("GET", translateUrl, nil)

	var query = req.URL.Query()
	query.Add("client", "gtx")
	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)

	req.URL.RawQuery = query.Encode()

	if err1 != nil {
		log.Fatalf("Masalah: %s", err1)
	}

	var res, err2 = client.Do(req)
	if err2 != nil {
		log.Fatalf("Masalah: %s", err2)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		str <- "Kamu sudah mengakses lebih dari batasan maksimal. Silakan coba lain waktu."
		wg.Done()
		return
	}

	var parsedJson, err3 = gabs.ParseJSONBuffer(res.Body)
	if err3 != nil {
		log.Fatalf("Masalah: %s", err3)
	}

	var nestOne, err4 = parsedJson.ArrayElement(0)
	if err4 != nil {
		log.Fatalf("Masalah: %s", err4)
	}

	var nestTwo, err5 = nestOne.ArrayElement(0)
	if err5 != nil {
		log.Fatalf("Masalah: %s", err5)
	}

	var translatedStr, err6 = nestTwo.ArrayElement(0)
	if err6 != nil {
		log.Fatalf("Masalah: %s", err6)
	}

	str <- translatedStr.Data().(string)
	wg.Done()
}
