package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"terraform-plan-parser/internals"
)

func contains(slice []string, entity string) bool {
	for _, sliceItem := range slice {
		if sliceItem == entity {
			return true
		}
	}
	return false
}

func main() {
	planPath := flag.String("planjson", "", "Path to the plan JSON file")
	useColour := flag.Bool("colour", false, "Output result to screen with colour")
	outFile := flag.String("output", "", "Optional file to output results to")
	flag.Parse()

	if *planPath == "" {
		log.Fatalf("planjson is a required argument")
	}

	parser := internals.NewParser()

	result, err := parser.ParseFile(*planPath)
	if err != nil {
		log.Fatalf("can't process file %s: %v", *planPath, err)
	}

	interested := make([]internals.ResourceChange, 0)

	var outFh *os.File
	outputToFile := false
	if *outFile != "" {
		outFh, err = os.Create(*outFile)
		if err != nil {
			log.Fatalf("Failed to open output file: %v", err)
		}
		defer outFh.Close()
		outputToFile = true
	}

	for _, resChg := range result.ResourceChanges {
		if contains(resChg.Changes.Actions, "no-op") {
			continue
		}

		interested = append(interested, resChg)
		changeTextPlain := strings.Join(resChg.Changes.Actions, "/")
		changeText := changeTextPlain
		if *useColour {
			switch changeText {
			case "create":
				changeText = AddColour(changeTextPlain, COLOUR_GREEN)
			case "update":
				changeText = AddColour(changeTextPlain, COLOUR_ORANGE)
			case "delete":
				changeText = AddColour(changeTextPlain, COLOUR_RED)
			case "delete/create":
				changeText = AddColour(changeTextPlain, COLOUR_PURPLE)
			}
		}
		msg := fmt.Sprintf("%s %s", changeText, resChg.Address)
		log.Print(msg)
		if outputToFile {
			outFh.WriteString(fmt.Sprintf("%s %s\n", changeTextPlain, resChg.Address))

		}
	}
}

type TextColour int64

const (
	COLOUR_RED TextColour = iota
	COLOUR_GREEN
	COLOUR_ORANGE
	COLOUR_PURPLE
)

func AddColour(msg string, colour TextColour) string {
	code := ""
	switch colour {
	case COLOUR_GREEN:
		code = "1;32"
	case COLOUR_RED:
		code = "1;31"
	case COLOUR_ORANGE:
		code = "1;33"
	case COLOUR_PURPLE:
		code = "1;35"
	}
	result := fmt.Sprintf("\033[%sm%s\033[0m", code, msg)
	return result
}
