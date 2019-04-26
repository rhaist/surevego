package surevego

import (
	"encoding/json"
	"log"
	"sync"
	"testing"
)

func ExampleLoadEveJSONFile() {
	ee, ec := LoadEveJSONFile("pathto/eve.json")

	// Fork handling of parsing errors to a gofunc
	go func() {
		for err := range ec {
			log.Fatal("[ERR]", err)
		}
	}()

	// Range over the events and print dns answers to stdout
	for event := range ee {
		if event.DNS != nil && event.DNS.Type == "answer" {
			log.Println(event.DNS)
		}
	}
}

func TestLoadEveJSONFile(t *testing.T) {
	var countTotal int
	var countDNS int
	var countFlow int

	ee, ec := LoadEveJSONFile("testdata/eve.json")

	go func() {
		for err := range ec {
			t.Error(err)
		}
	}()

	for event := range ee {
		if event.DNS != nil {
			countDNS++
		}
		if event.Flow != nil {
			countFlow++
		}
		if event.EventType == "" || event.Timestamp.IsZero() {
			t.Error("Mandatory field missing")
		}
		countTotal++
	}

	if countDNS != 48 || countFlow != 13 || countTotal != 266 {
		t.Error("Event count mismatch")
	}
}

func TestLoadBrokenEveJSONFile(t *testing.T) {
	var countErrors int
	var wg sync.WaitGroup

	ee, ec := LoadEveJSONFile("testdata/eve_broken.json")

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range ec {
			t.Log(err)
			countErrors++
		}
	}()

	for event := range ee {
		if event.EventType == "" || event.Timestamp.IsZero() {
			t.Error("Mandatory field missing")
		}
	}

	wg.Wait()
	if countErrors < 1 {
		t.Error("Error count mismatch")
	}
}

func TestMissingJSONFile(t *testing.T) {
	var countErrors int
	var wg sync.WaitGroup

	_, ec := LoadEveJSONFile("nonexistant")

	wg.Add(1)
	go func(myWg *sync.WaitGroup, myEc <-chan error) {
		defer myWg.Done()
		for err := range myEc {
			t.Log(err)
			countErrors++
			break
		}
	}(&wg, ec)

	wg.Wait()
	if countErrors < 1 {
		t.Error("Error count mismatch")
	}
}

func TestMarshalWithTimestamp(t *testing.T) {
	ee, ec := LoadEveJSONFile("testdata/eve.json")

	go func() {
		for err := range ec {
			t.Error(err)
		}
	}()

	e := <-ee

	out, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}

	var inEVE EveEvent
	err = json.Unmarshal(out, &inEVE)
	if err != nil {
		t.Error(err)
	}

	if inEVE.Timestamp.Time != e.Timestamp.Time {
		t.Fatalf("timestamp round-trip failed: %v <-> %v", inEVE.Timestamp, e.Timestamp)
	}
}
