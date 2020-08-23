package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// Service that provides the face analysis
// takes an io.Reader to the image
type faceAnalyzer interface {
	analyzeFaces(image io.Reader) chan imageAnalysis
}

func main() {
	subKey := os.Getenv("MSCV_SUBKEY")
	endPoint := os.Getenv("MSCV_ENDPOINT")
	if subKey == "" || endPoint == "" {
		log.Fatal("Set subscription key and endpoint env variables")
	}

	if len(os.Args) < 2 {
		log.Fatal("Input the path to the directory of images as first arg")
	}

	pathToImages := os.Args[1]
	validateDir(pathToImages)

	dir, err := os.Open(pathToImages)
	defer dir.Close()

	if err != nil {
		log.Fatal(err.Error())
	}

	fileInfos, err := dir.Readdir(-1)

	if err != nil {
		log.Fatal(err.Error())
	}

	var cv faceAnalyzer = newMsCvAnlayzer(subKey, endPoint)

	result := make(map[string]int)
	responses := make([]chan imageAnalysis, 0, len(fileInfos))

	// Two for loops to use make parallel requests
	for _, fi := range fileInfos {
		f, _ := os.Open(pathToImages + "/" + fi.Name())
		resp := cv.analyzeFaces(f)
		responses = append(responses, resp)
		<-time.After(time.Second * 5)
	}

	for i, resp := range responses {
		fname := fileInfos[i].Name()
		analysis, ok := <-resp
		if !ok {
			fmt.Printf("Could not analyze %s\n", fname)
			continue
		}
		result[fname] = analysis.numberOfFaces
	}

	printFormatted(result)

}

// Makes sure the argument passed is a directory,
// if not exits the process
func validateDir(path string) {
	pathInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	if !pathInfo.IsDir() {
		log.Fatalf("Argument path %s is not a directory", path)
	}
}

//Formats according to the requirement
func printFormatted(r map[string]int) {
	total := 0
	fmt.Println("filename,faces")
	for f, n := range r {
		fmt.Println(f + "," + strconv.Itoa(n))
		total += n
	}
	fmt.Println("Total," + strconv.Itoa(total))
}
