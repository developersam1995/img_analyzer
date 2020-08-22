package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Service that implements faceAnalyzer
// Having a single underlying http.Client for optimization
type microsoftCv struct {
	subKey     string
	endPoint   string
	analyzeURI string
	hClient    *http.Client
}

// result of the analysis
type imageAnalysis struct {
	numberOfFaces int
}

// Microsoft CV response
type analyzedResponse struct {
	Faces []face `json:"faces"`
}
type face struct {
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

func newMsCvAnlayzer(subKey, endPoint string) *microsoftCv {
	c := &http.Client{}
	if endPoint[0:4] != "http" {
		endPoint = "https://" + endPoint
	}
	return &microsoftCv{
		subKey:     subKey,
		endPoint:   endPoint,
		analyzeURI: "/vision/v3.0/analyze",
		hClient:    c,
	}
}

func (m *microsoftCv) analyzeFaces(i io.Reader) chan imageAnalysis {
	ia := make(chan imageAnalysis)

	req, newReqErr := http.NewRequest("POST", m.endPoint+m.analyzeURI+"?visualFeatures=Faces", i)

	if newReqErr != nil {
		fmt.Println(newReqErr)
		close(ia)
	}

	req.Header.Add("content-type", "application/octet-stream")
	req.Header.Add("Ocp-Apim-Subscription-Key", m.subKey)

	go func() {
		res, reqErr := m.hClient.Do(req)

		if reqErr != nil {
			fmt.Println(reqErr)
			close(ia)
		}

		respBody, bodyErr := ioutil.ReadAll(res.Body)
		if bodyErr != nil {
			fmt.Println(bodyErr)
			close(ia)
		}

		resp := analyzedResponse{}
		jsonErr := json.Unmarshal(respBody, &resp)
		if jsonErr != nil {
			log.Println(jsonErr)
			close(ia)
		}

		ia <- imageAnalysis{
			numberOfFaces: len(resp.Faces),
		}
	}()
	return ia
}
