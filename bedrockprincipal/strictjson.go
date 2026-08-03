package bedrockprincipal

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

func strictObject(raw []byte, target any) error {
	if !utf8.Valid(raw) || len(raw) > MaxPayloadBytes {
		return Malformed
	}
	if err := validateJSONStringEscapes(raw); err != nil {
		return err
	}
	stringsBytes := 0
	shape := json.NewDecoder(strings.NewReader(string(raw)))
	shape.UseNumber()
	if err := parseStrictObject(shape, 1, &stringsBytes); err != nil {
		return err
	}
	if stringsBytes > 24576 || len(raw)+stringsBytes > 65536 {
		return Malformed
	}
	if token, err := shape.Token(); err != io.EOF || token != nil {
		return Malformed
	}

	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	return ensureJSONEOF(decoder)
}

func validateJSONStringEscapes(raw []byte) error {
	inString := false
	for i := 0; i < len(raw); i++ {
		switch raw[i] {
		case '"':
			inString = !inString
		case '\\':
			if !inString || i+1 >= len(raw) {
				return Malformed
			}
			if raw[i+1] != 'u' {
				i++
				continue
			}
			value, ok := parseHex16(raw, i+2)
			if !ok {
				return Malformed
			}
			if value >= 0xd800 && value <= 0xdbff {
				if i+12 > len(raw) || raw[i+6] != '\\' || raw[i+7] != 'u' {
					return Malformed
				}
				low, ok := parseHex16(raw, i+8)
				if !ok || low < 0xdc00 || low > 0xdfff {
					return Malformed
				}
				i += 11
				continue
			}
			if value >= 0xdc00 && value <= 0xdfff {
				return Malformed
			}
			i += 5
		}
	}
	if inString {
		return Malformed
	}
	return nil
}

func parseHex16(raw []byte, start int) (uint16, bool) {
	if start+4 > len(raw) {
		return 0, false
	}
	var result uint16
	for _, value := range raw[start : start+4] {
		result <<= 4
		switch {
		case value >= '0' && value <= '9':
			result |= uint16(value - '0')
		case value >= 'a' && value <= 'f':
			result |= uint16(value-'a') + 10
		case value >= 'A' && value <= 'F':
			result |= uint16(value-'A') + 10
		default:
			return 0, false
		}
	}
	return result, true
}

func requireObjectMembers(raw []byte, required []string, optional ...string) error {
	var members map[string]json.RawMessage
	if err := json.Unmarshal(raw, &members); err != nil {
		return Malformed
	}
	allowed := make(map[string]bool, len(required)+len(optional))
	for _, name := range required {
		allowed[name] = true
		if _, ok := members[name]; !ok {
			return Malformed
		}
	}
	for _, name := range optional {
		allowed[name] = true
	}
	for name := range members {
		if !allowed[name] {
			return Malformed
		}
	}
	return nil
}

func rejectDuplicateMembers(raw []byte, maxDepth int) error {
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.UseNumber()
	first, err := decoder.Token()
	if err != nil {
		return err
	}
	if err := validateJSONToken(decoder, first, 1, maxDepth); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return Malformed
	}
	return nil
}

func validateJSONToken(decoder *json.Decoder, token json.Token, depth, maxDepth int) error {
	if depth > maxDepth {
		return Malformed
	}
	delim, compound := token.(json.Delim)
	if !compound {
		if value, ok := token.(string); ok && strings.ContainsRune(value, '\x00') {
			return Malformed
		}
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok || seen[key] {
				return Malformed
			}
			seen[key] = true
			value, err := decoder.Token()
			if err != nil {
				return err
			}
			if err := validateJSONToken(decoder, value, depth+1, maxDepth); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim('}') {
			return Malformed
		}
	case '[':
		for decoder.More() {
			value, err := decoder.Token()
			if err != nil {
				return err
			}
			if err := validateJSONToken(decoder, value, depth+1, maxDepth); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim(']') {
			return Malformed
		}
	default:
		return Malformed
	}
	return nil
}

func parseStrictObject(decoder *json.Decoder, depth int, stringsBytes *int) error {
	if depth > 4 {
		return Malformed
	}
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return Malformed
	}
	seen := map[string]struct{}{}
	members := 0
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			return Malformed
		}
		key, ok := keyToken.(string)
		if !ok || strings.ContainsRune(key, '\x00') {
			return Malformed
		}
		if _, duplicate := seen[key]; duplicate {
			return Malformed
		}
		seen[key] = struct{}{}
		members++
		if members > 32 {
			return Malformed
		}
		*stringsBytes += len(key)
		if err := parseStrictValue(decoder, depth, stringsBytes); err != nil {
			return err
		}
	}
	closing, err := decoder.Token()
	if err != nil || closing != json.Delim('}') {
		return Malformed
	}
	return nil
}

func parseStrictValue(decoder *json.Decoder, depth int, stringsBytes *int) error {
	token, err := decoder.Token()
	if err != nil {
		return Malformed
	}
	switch value := token.(type) {
	case json.Delim:
		if value != '{' {
			return Malformed
		}
		// Token already consumed; parse the object body inline.
		if depth+1 > 4 {
			return Malformed
		}
		seen := map[string]struct{}{}
		members := 0
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return Malformed
			}
			key, ok := keyToken.(string)
			if !ok || strings.ContainsRune(key, '\x00') {
				return Malformed
			}
			if _, duplicate := seen[key]; duplicate {
				return Malformed
			}
			seen[key] = struct{}{}
			members++
			if members > 32 {
				return Malformed
			}
			*stringsBytes += len(key)
			if err := parseStrictValue(decoder, depth+1, stringsBytes); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim('}') {
			return Malformed
		}
	case string:
		if strings.ContainsRune(value, '\x00') {
			return Malformed
		}
		*stringsBytes += len(value)
	case json.Number:
		if strings.ContainsAny(string(value), ".eE") {
			return Malformed
		}
	case bool, nil:
	default:
		return fmt.Errorf("unexpected JSON token")
	}
	return nil
}
