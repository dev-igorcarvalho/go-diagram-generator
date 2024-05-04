package processor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

func CreateDiagram(diagramProcessor func()) {
	clearTempFiles()

	var wg sync.WaitGroup
	wg.Add(1)

	funcName := getFunctionName(diagramProcessor)
	snakeCaseName := toSnakeCase(funcName)
	log.Printf("Creating diagram: %s\n", snakeCaseName)

	go func() {
		defer wg.Done()
		diagramProcessor()
	}()

	wg.Wait()

	// Call saveDiagramAsPNG after waiting for myDancingApi
	saveDiagramAsPNG(snakeCaseName)
	clearTempFiles()
}

func getFunctionName(f interface{}) string {
	fullQualifiedName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	parts := strings.Split(fullQualifiedName, ".")
	return parts[len(parts)-1]
}

func toSnakeCase(s string) string {
	var re = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}

func saveDiagramAsPNG(name string) {
	log.Printf("Saving diagram")
	// Construct the output PNG file name
	outputFile := fmt.Sprintf("%s.png", name)

	rootFolder := getWorkingDirectory()

	goToDiagramsFolder(rootFolder)

	// Run the shell command to convert diagram.dot to diagram.png
	runGraphViz(outputFile)

	// Move the PNG file to the output folder
	createOutputFile(rootFolder, outputFile)

	goToRootFolder()
}

func createOutputFile(folder string, outputFile string) {
	// Move the PNG file to the output folder
	outputFolder := filepath.Join(folder, "output")
	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	outputFilePath := filepath.Join(outputFolder, outputFile)
	if err := os.Rename(outputFile, outputFilePath); err != nil {
		log.Fatal(err)
	}
	log.Printf("Diagram saved as: %s", outputFilePath)
}

func runGraphViz(outputFile string) {
	// Run the shell command to convert diagram.dot to diagram.png
	cmd := exec.Command("dot", "-Tpng", "go-diagram.dot", "-o", outputFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func goToRootFolder() {
	// Change back to the previous working directory
	if err := os.Chdir(".."); err != nil {
		log.Fatal(err)
	}
}

func goToDiagramsFolder(folder string) {
	// Change the working directory to go-diagrams folder
	if err := os.Chdir(filepath.Join(folder, "go-diagrams")); err != nil {
		log.Fatal(err)
	}
}

func getWorkingDirectory() string {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return wd
}

func clearTempFiles() {
	log.Printf("Clearing temporary files.")
	// Get the working directory
	wd := getWorkingDirectory()

	// Check if the go-diagrams folder exists
	goDiagramsFolder := filepath.Join(wd, "go-diagrams")
	if _, err := os.Stat(goDiagramsFolder); os.IsNotExist(err) {
		log.Printf("Temporary does not exist.")
		return
	}

	// Remove the go-diagrams folder and its contents
	if err := os.RemoveAll(goDiagramsFolder); err != nil {
		log.Fatal(err)
	}

	log.Printf("Cleaned up all temporary files.")
}
