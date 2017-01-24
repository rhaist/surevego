package surevego

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strings"
)

// LoadEveJSONFile loads a suricata eve.json from a given path and returns a
// result channel containing marshaled events and a error channel.
func LoadEveJSONFile(path string) (<-chan EveEvent, <-chan error) {
	eventChan := make(chan EveEvent)
	errorChan := make(chan error)

	file, err := os.Open(path)
	if err != nil {
		errorChan <- err
		close(eventChan)
		close(errorChan)
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

// LoadEveString loads a eve formated string and returns a struct of EveEvents
func LoadEveString(eveString string) (EveEvents, error) {
	var ev EveEvent
	var events EveEvents

	scanner := bufio.NewScanner(strings.NewReader(eveString))
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &ev)
		if err != nil {
			return events, errors.New("Error unmarshaling eve.json line: " +
				err.Error() + " " + scanner.Text())
		}
		events.Events = append(events.Events, ev)
	}
	return events, nil
}
