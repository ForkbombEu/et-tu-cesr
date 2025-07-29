package cesr

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
)

var hdrRE = regexp.MustCompile(`\{"v":"((?:KERI|ACDC)[0-9]{2}[A-Z]{4}[0-9A-F]{6}_)`)

type Event struct {
	KED         map[string]interface{}
	AttachBytes int
}

// ParseCESR parses a concatenated CESR stream into events.
func ParseCESR(stream string) ([]Event, error) {
	var out []Event
	offset := 0

	for offset < len(stream) {
		loc := hdrRE.FindStringSubmatchIndex(stream[offset:])
		if loc == nil {
			break
		}
		hdrStart := offset + loc[0]
		header := stream[offset+loc[2] : offset+loc[3]]

		bodyLen, err := bodyLength(header)
		if err != nil {
			return nil, err
		}
		bodyStart := hdrStart
		bodyEnd := bodyStart + bodyLen
		if bodyEnd > len(stream) {
			return nil, errors.New("truncated body")
		}
		bodyRaw := stream[bodyStart:bodyEnd]

		// Search for next header in the slice AFTER bodyEnd
		next := hdrRE.FindStringIndex(stream[bodyEnd:])
		attEnd := len(stream)
		if next != nil {
			attEnd = bodyEnd + next[0]
		}
		atcRaw := stream[bodyEnd:attEnd]

		var ked map[string]interface{}
		if err := json.Unmarshal([]byte(bodyRaw), &ked); err != nil {
			return nil, err
		}
		out = append(out, Event{KED: ked, AttachBytes: len(atcRaw)})

		offset = attEnd
	}
	return out, nil
}

// bodyLength extracts the hex body‑length from a 17‑char CESR header string.
func bodyLength(h string) (int, error) {
	if len(h) != 17 {
		return 0, errors.New("bad header length")
	}
	n64, err := strconv.ParseInt(h[10:16], 16, 0)
	return int(n64), err
}
