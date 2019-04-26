package surevego

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

// LoadEveJSONFile reads a suricata eve.json from a given path and returns two
// channels. One for parsed EveEvents and one for parsing errors. Those two
// channels need to be handled separately.
func LoadEveJSONFile(path string) (<-chan EveEvent, <-chan error) {
	eventChan := make(chan EveEvent)
	errorChan := make(chan error)

	file, err := os.Open(path)
	if err != nil {
		go func() {
			errorChan <- err
			close(eventChan)
			close(errorChan)
		}()
		return nil, errorChan
	}

	scanner := bufio.NewScanner(file)

	go func() {
		for scanner.Scan() {
			ev := EveEvent{}
			err = json.Unmarshal(scanner.Bytes(), &ev)
			if err != nil {
				errorChan <- errors.New("Error unmarshaling eve.json line: " +
					err.Error() + " " + scanner.Text())
			} else {
				eventChan <- ev
			}
		}
		close(eventChan)
		close(errorChan)
		defer file.Close()
	}()

	return eventChan, errorChan
}
