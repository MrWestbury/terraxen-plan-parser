package internals

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Parser struct{}

func NewParser() *Parser {
	p := &Parser{}

	return p
}

func (p *Parser) ParseFile(planPath string) (*Plan, error) {
	fh, err := os.Open(planPath)
	if err != nil {
		log.Printf("failed to open file: %v", err)
		return nil, err
	}
	defer fh.Close()

	byteValue, _ := ioutil.ReadAll(fh)

	var result Plan
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		log.Printf("failed to bind json: %v", err)
		return nil, err
	}

	return &result, nil
}
