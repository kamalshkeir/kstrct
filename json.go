package kstrct

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Forward declarations
var parseObjectToMapRecursive func(str string, currentIdx int) (m map[string]any, nextIdx int, err error)
var parseArrayToListRecursive func(str string, currentIdx int) (list []any, nextIdx int, err error)

// --- New Byte-based parsing functions ---

var (
	// Forward declarations for BYTE based recursive parsers
	parseObjectToMapRecursiveBytes func(data []byte, currentIdx int) (m map[string]any, nextIdx int, err error)
	parseArrayToListRecursiveBytes func(data []byte, currentIdx int) (list []any, nextIdx int, err error)
)

// parseArbitraryValueRecursiveBytes parses any JSON value from a byte slice.
func parseArbitraryValueRecursiveBytes(data []byte, currentIdx int) (value any, nextIdx int, err error) {
	localSkipWhitespaceBytes := func(d []byte, i int) int {
		for i < len(d) && (d[i] == ' ' || d[i] == '\n' || d[i] == '\r' || d[i] == '\t') {
			i++
		}
		return i
	}

	idx := localSkipWhitespaceBytes(data, currentIdx)

	if idx >= len(data) {
		return nil, idx, fmt.Errorf("bytes: unexpected end of data when expecting a value at or after index %d", currentIdx)
	}

	switch data[idx] {
	case '"': // String value parsing with escape sequence handling
		idx++ // Consume opening '"'
		stringStartOriginalIndex := idx
		var sb strings.Builder
		builderUsed := false
		segmentStart := idx

		for idx < len(data) {
			charByte := data[idx]

			if charByte == '\\' {
				if !builderUsed {
					builderUsed = true
				}
				sb.Write(data[segmentStart:idx]) // Append segment before '\\'

				idx++ // Consume '\\'
				if idx >= len(data) {
					return nil, stringStartOriginalIndex - 1, fmt.Errorf("bytes: unterminated escape sequence at end of data starting at %d", stringStartOriginalIndex-1)
				}

				escapeChar := data[idx]
				switch escapeChar {
				case '"', '\\', '/':
					sb.WriteByte(escapeChar)
					idx++
				case 'b':
					sb.WriteByte('\b')
					idx++
				case 'f':
					sb.WriteByte('\f')
					idx++
				case 'n':
					sb.WriteByte('\n')
					idx++
				case 'r':
					sb.WriteByte('\r')
					idx++
				case 't':
					sb.WriteByte('\t')
					idx++
				case 'u':
					if idx+4 >= len(data) {
						return nil, stringStartOriginalIndex - 1, fmt.Errorf("bytes: unterminated \\\\uXXXX: not enough chars for hex sequence after \\\\u at index %d", idx)
					}
					hexSeq := data[idx+1 : idx+1+4]
					unicodeVal, errHex := strconv.ParseUint(string(hexSeq), 16, 16)
					if errHex != nil {
						return nil, stringStartOriginalIndex - 1, fmt.Errorf("bytes: invalid \\\\uXXXX escape sequence '%s' in data starting at %d: %w", string(hexSeq), stringStartOriginalIndex-1, errHex)
					}
					sb.WriteRune(rune(unicodeVal))
					idx += (1 + 4)
				default:
					return nil, stringStartOriginalIndex - 1, fmt.Errorf("bytes: invalid escape character '%c' (0x%x) after \\\\ in data starting at %d", escapeChar, escapeChar, stringStartOriginalIndex-1)
				}
				segmentStart = idx
				continue
			} else if charByte == '"' {
				if !builderUsed {
					value = string(data[segmentStart:idx]) // Conversion to string
				} else {
					sb.Write(data[segmentStart:idx])
					value = sb.String()
				}
				nextIdx = idx + 1
				return value, nextIdx, nil
			}
			idx++
		}
		return nil, stringStartOriginalIndex - 1, fmt.Errorf("bytes: unterminated string literal starting at %d", stringStartOriginalIndex-1)

	case 't':
		if idx+3 < len(data) && bytes.Equal(data[idx:idx+4], []byte("true")) {
			return true, idx + 4, nil
		}
		return nil, idx, fmt.Errorf("bytes: expected 'true' at index %d but found '%s'", idx, string(data[idx:minBytes(idx+4, len(data))]))

	case 'f':
		if idx+4 < len(data) && bytes.Equal(data[idx:idx+5], []byte("false")) {
			return false, idx + 5, nil
		}
		return nil, idx, fmt.Errorf("bytes: expected 'false' at index %d but found '%s'", idx, string(data[idx:minBytes(idx+5, len(data))]))

	case 'n':
		if idx+3 < len(data) && bytes.Equal(data[idx:idx+4], []byte("null")) {
			return nil, idx + 4, nil
		}
		return nil, idx, fmt.Errorf("bytes: expected 'null' at index %d but found '%s'", idx, string(data[idx:minBytes(idx+4, len(data))]))

	case '{':
		return parseObjectToMapRecursiveBytes(data, idx)

	case '[':
		return parseArrayToListRecursiveBytes(data, idx)

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		numParseStartIdx := idx
		scanEndIdx := idx
		hasDecimal := false

		// Scanner simplement les caractères potentiels d'un nombre
		// 1. Signe optionnel
		if scanEndIdx < len(data) && data[scanEndIdx] == '-' {
			scanEndIdx++
		}
		// 2. Partie entière
		for scanEndIdx < len(data) && (data[scanEndIdx] >= '0' && data[scanEndIdx] <= '9') {
			scanEndIdx++
		}
		// 3. Partie décimale optionnelle
		if scanEndIdx < len(data) && data[scanEndIdx] == '.' {
			hasDecimal = true
			scanEndIdx++
			startFractional := scanEndIdx
			for scanEndIdx < len(data) && (data[scanEndIdx] >= '0' && data[scanEndIdx] <= '9') {
				scanEndIdx++
			}
			// Vérifier si au moins un chiffre a suivi le point décimal
			if scanEndIdx == startFractional {
				return nil, numParseStartIdx, fmt.Errorf("bytes: malformed number, decimal point not followed by digits at index %d", numParseStartIdx)
			}
		}

		// Extraire la slice scannée
		numStrSlice := data[numParseStartIdx:scanEndIdx]

		// Vérification finale de validité (ex: "", "-", ".")
		if len(numStrSlice) == 0 || bytes.Equal(numStrSlice, []byte("-")) || bytes.Equal(numStrSlice, []byte(".")) {
			return nil, numParseStartIdx, fmt.Errorf("bytes: invalid number format scanned: '%s'", string(numStrSlice))
		}
		// Vérifier aussi que si on a commencé par '-', il y a bien des chiffres après
		if data[numParseStartIdx] == '-' && len(numStrSlice) == 1 {
			return nil, numParseStartIdx, fmt.Errorf("bytes: invalid number format scanned: '%s'", string(numStrSlice))
		}
		// Vérifier aussi qu'on n'a pas juste des chiffres suivis d'un point (ex: "123.")
		if hasDecimal && data[scanEndIdx-1] == '.' { // Vérifie si le dernier caractère scanné était le point
			return nil, numParseStartIdx, fmt.Errorf("bytes: invalid number format scanned, ends with decimal point: '%s'", string(numStrSlice))
		}

		// Convertir avec strconv
		if hasDecimal {
			fVal, errFloat := strconv.ParseFloat(string(numStrSlice), 64)
			if errFloat != nil {
				return nil, numParseStartIdx, fmt.Errorf("bytes: could not parse float '%s' from index %d: %w", string(numStrSlice), numParseStartIdx, errFloat)
			}
			return fVal, scanEndIdx, nil
		} else {
			iVal, errInt := strconv.ParseInt(string(numStrSlice), 10, 64)
			if errInt != nil {
				// Try float for cases like 1e5 which ParseInt fails on
				fVal, errFloat := strconv.ParseFloat(string(numStrSlice), 64)
				if errFloat != nil {
					return nil, numParseStartIdx, fmt.Errorf("bytes: could not parse number (tried int then float) '%s' from index %d: int_err=%w, float_err=%v", string(numStrSlice), numParseStartIdx, errInt, errFloat)
				}
				return fVal, scanEndIdx, nil // Return float if it parsed
			}
			return iVal, scanEndIdx, nil
		}

	default:
		return nil, idx, fmt.Errorf("bytes: unsupported value type or invalid character '%c' starting at index %d", safeCharBytes(data, idx), idx)
	}
}

func safeCharBytes(d []byte, i int) byte {
	if i < 0 || i >= len(d) {
		return '?'
	}
	return d[i]
}

func minBytes(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ParseBytes parses a JSON byte slice into a slice of KV.
func ParseBytes(data []byte) []KV {
	res := make([]KV, 0, 11) // Estimate based on typical content

	var currentKeyString string // KV.Key is string
	idx := 0

	skipWhitespaceBytes := func() {
		for idx < len(data) && (data[idx] == ' ' || data[idx] == '\n' || data[idx] == '\r' || data[idx] == '\t') {
			idx++
		}
	}

	skipWhitespaceBytes()
	if idx >= len(data) || data[idx] != '{' {
		// Consider returning error: return nil, fmt.Errorf("...")
		return nil
	}
	idx++ // Consume '{'

	firstTopLevelItem := true
	for idx < len(data) {
		skipWhitespaceBytes()
		if idx >= len(data) {
			// fmt.Printf("Error: Unexpected end of data, expecting key or '}' for main object\\n")
			return nil
		}

		if data[idx] == '}' {
			idx++
			break
		}

		if !firstTopLevelItem {
			if data[idx] != ',' {
				// fmt.Printf("Error: Expected ',' or '}' in main object at index %d, found '%c'\\n", idx, safeCharBytes(data, idx))
				return nil
			}
			idx++ // Consume ','
			skipWhitespaceBytes()
			if idx >= len(data) {
				// fmt.Printf("Error: Unexpected end of data after comma in main object\\n")
				return nil
			}
			if data[idx] == '}' { // Trailing comma
				// fmt.Printf("Error: Trailing comma in main object at index %d\\n", idx-1)
				return nil
			}
		}
		firstTopLevelItem = false

		if data[idx] != '"' {
			// fmt.Printf("Error: Expected '\\\"' for key start in main object at index %d, found '%c'\\n", idx, safeCharBytes(data, idx))
			return nil
		}
		idx++ // Consume '"'

		keyStartIdx := idx
		for idx < len(data) && data[idx] != '"' {
			// Assuming keys do not contain escape sequences
			idx++
		}
		if idx >= len(data) {
			// fmt.Printf("Error: Unterminated key string in main object, starting at %d\\n", keyStartIdx-1)
			return nil
		}
		currentKeyString = string(data[keyStartIdx:idx]) // Key converted to string
		idx++                                            // Consume closing '"' for key

		skipWhitespaceBytes()

		if idx >= len(data) || data[idx] != ':' {
			// fmt.Printf("Error: Expected ':' after key '%s' in main object at index %d, found '%c'\\n", currentKeyString, idx, safeCharBytes(data, idx))
			return nil
		}
		idx++ // Consume ':'

		var currentValue any
		var errValue error
		currentValue, idx, errValue = parseArbitraryValueRecursiveBytes(data, idx)
		if errValue != nil {
			// fmt.Printf("Error parsing value for key '%s' in main object: %v\\n", currentKeyString, errValue)
			return nil
		}

		res = append(res, KV{Key: currentKeyString, Value: currentValue})
		skipWhitespaceBytes()
	}
	return res
}

// parseArbitraryValueRecursive parses any JSON value (string, number, boolean, object, array, null)
func parseArbitraryValueRecursive(str string, currentIdx int) (value any, nextIdx int, err error) {
	localSkipWhitespace := func(s string, i int) int {
		for i < len(s) && (s[i] == ' ' || s[i] == '\n' || s[i] == '\r' || s[i] == '\t') {
			i++
		}
		return i
	}

	idx := localSkipWhitespace(str, currentIdx)

	if idx >= len(str) {
		return nil, idx, fmt.Errorf("unexpected end of string when expecting a value at or after index %d", currentIdx)
	}

	switch str[idx] {
	case '"': // String value parsing with escape sequence handling
		idx++                           // Consume opening '"'
		stringStartOriginalIndex := idx // For error reporting, points to the first char of content
		var sb strings.Builder
		builderUsed := false
		segmentStart := idx

		for idx < len(str) {
			charByte := str[idx]

			if charByte == '\\' { // Escape sequence
				if !builderUsed {
					builderUsed = true
					// sb.Grow() could be used here if we have a good estimate
				}
				sb.WriteString(str[segmentStart:idx]) // Append segment before '\'

				idx++ // Consume '\'
				if idx >= len(str) {
					return nil, stringStartOriginalIndex - 1, fmt.Errorf("unterminated escape sequence at end of string starting at %d", stringStartOriginalIndex-1)
				}

				escapeChar := str[idx]
				switch escapeChar {
				case '"', '\\', '/':
					sb.WriteByte(escapeChar)
					idx++
				case 'b':
					sb.WriteByte('\b')
					idx++
				case 'f':
					sb.WriteByte('\f')
					idx++
				case 'n':
					sb.WriteByte('\n')
					idx++
				case 'r':
					sb.WriteByte('\r')
					idx++
				case 't':
					sb.WriteByte('\t')
					idx++
				case 'u':
					// idx currently points at 'u'
					if idx+4 >= len(str) { // Need 4 more characters for XXXX after 'u'
						return nil, stringStartOriginalIndex - 1, fmt.Errorf("unterminated \\uXXXX: not enough chars for hex sequence after \\u at index %d", idx)
					}
					hexSeq := str[idx+1 : idx+1+4] // Get XXXX (4 chars after 'u')
					unicodeVal, errHex := strconv.ParseUint(hexSeq, 16, 16)
					if errHex != nil {
						return nil, stringStartOriginalIndex - 1, fmt.Errorf("invalid \\uXXXX escape sequence '%s' in string starting at %d: %w", hexSeq, stringStartOriginalIndex-1, errHex)
					}
					sb.WriteRune(rune(unicodeVal))
					idx += (1 + 4) // Consume 'u' and 'XXXX'
				default:
					return nil, stringStartOriginalIndex - 1, fmt.Errorf("invalid escape character '%c' (0x%x) after \\ in string starting at %d", escapeChar, escapeChar, stringStartOriginalIndex-1)
				}
				segmentStart = idx // Next segment for direct copy starts after the processed escape sequence
				continue           // Restart loop from new position
			} else if charByte == '"' { // End of string
				if !builderUsed {
					value = str[segmentStart:idx] // No escapes, direct slice
				} else {
					sb.WriteString(str[segmentStart:idx]) // Append last segment
					value = sb.String()
				}
				nextIdx = idx + 1 // Consume closing '"'
				return value, nextIdx, nil
			}
			// Regular character, just advance index for next iteration. Content will be copied by WriteString in batches.
			idx++
		}
		// If loop finishes without finding a closing quote, it's an unterminated string.
		return nil, stringStartOriginalIndex - 1, fmt.Errorf("unterminated string literal starting at %d", stringStartOriginalIndex-1)

	case 't': // true
		if idx+3 < len(str) && str[idx:idx+4] == "true" {
			return true, idx + 4, nil
		}
		return nil, idx, fmt.Errorf("expected 'true' at index %d but found '%s'", idx, str[idx:min(idx+4, len(str))])

	case 'f': // false
		if idx+4 < len(str) && str[idx:idx+5] == "false" {
			return false, idx + 5, nil
		}
		return nil, idx, fmt.Errorf("expected 'false' at index %d but found '%s'", idx, str[idx:min(idx+5, len(str))])

	case 'n': // null
		if idx+3 < len(str) && str[idx:idx+4] == "null" {
			return nil, idx + 4, nil
		}
		return nil, idx, fmt.Errorf("expected 'null' at index %d but found '%s'", idx, str[idx:min(idx+4, len(str))])

	case '{': // Object
		return parseObjectToMapRecursive(str, idx)

	case '[': // Array
		return parseArrayToListRecursive(str, idx)

	// Numbers (integer and float)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		numParseStartIdx := idx
		scanEndIdx := idx
		hasDecimal := false

		if str[scanEndIdx] == '-' {
			// Vérifier AVANT de consommer le '-' s'il y a un caractère après et si c'est un chiffre
			nextIdxCheck := scanEndIdx + 1
			if nextIdxCheck >= len(str) {
				return nil, numParseStartIdx, fmt.Errorf("bytes: invalid number sequence: '-' at end of input at index %d", numParseStartIdx)
			}
			if !(str[nextIdxCheck] >= '0' && str[nextIdxCheck] <= '9') {
				return nil, numParseStartIdx, fmt.Errorf("bytes: invalid number sequence: '-' not followed by a digit at index %d", numParseStartIdx)
			}
			// Si les vérifications passent, ALORS consommer le '-'
			scanEndIdx++
		}

		firstDigitScanEnd := scanEndIdx
		for scanEndIdx < len(str) && (str[scanEndIdx] >= '0' && str[scanEndIdx] <= '9') {
			scanEndIdx++
		}

		if str[numParseStartIdx] == '-' && scanEndIdx == firstDigitScanEnd {
			return nil, numParseStartIdx, fmt.Errorf("bytes: invalid number sequence: '-' not followed by digits at index %d", numParseStartIdx)
		}

		if scanEndIdx < len(str) && str[scanEndIdx] == '.' {
			hasDecimal = true
			scanEndIdx++
			if scanEndIdx >= len(str) || !(str[scanEndIdx] >= '0' && str[scanEndIdx] <= '9') {
				return nil, numParseStartIdx, fmt.Errorf("bytes: malformed number, missing digits after decimal point at index %d", numParseStartIdx)
			}
			for scanEndIdx < len(str) && (str[scanEndIdx] >= '0' && str[scanEndIdx] <= '9') {
				scanEndIdx++
			}
		}

		if scanEndIdx == numParseStartIdx && str[numParseStartIdx] != '-' && !(str[numParseStartIdx] >= '0' && str[numParseStartIdx] <= '9') {
			return nil, numParseStartIdx, fmt.Errorf("bytes: invalid character '%c' where number was expected at index %d", str[numParseStartIdx], numParseStartIdx)
		}

		numStr := str[numParseStartIdx:scanEndIdx]

		if hasDecimal {
			fVal, errFloat := strconv.ParseFloat(numStr, 64)
			if errFloat != nil {
				return nil, numParseStartIdx, fmt.Errorf("bytes: could not parse float '%s' from index %d: %w", numStr, numParseStartIdx, errFloat)
			}
			return fVal, scanEndIdx, nil
		} else {
			iVal, errInt := strconv.ParseInt(numStr, 10, 64)
			if errInt != nil {
				return nil, numParseStartIdx, fmt.Errorf("bytes: could not parse int '%s' from index %d: %w", numStr, numParseStartIdx, errInt)
			}
			return iVal, scanEndIdx, nil
		}

	default:
		return nil, idx, fmt.Errorf("bytes: unsupported value type or invalid character '%c' starting at index %d", str[idx], idx)
	}
}

func init() {
	parseObjectToMapRecursive = func(str string, currentIdx int) (m map[string]any, nextIdx int, err error) {
		localSkipWhitespace := func(s string, i int) int {
			for i < len(s) && (s[i] == ' ' || s[i] == '\n' || s[i] == '\r' || s[i] == '\t') {
				i++
			}
			return i
		}

		idx := localSkipWhitespace(str, currentIdx)

		if idx >= len(str) || str[idx] != '{' {
			return nil, currentIdx, fmt.Errorf("expected '{' for object start at index %d, but found '%c' (or EOF)", safeChar(str, idx), idx)
		}
		idx++ // Consume '{'

		m = make(map[string]any)
		firstItem := true

		for {
			idx = localSkipWhitespace(str, idx)
			if idx >= len(str) {
				return nil, currentIdx, fmt.Errorf("unterminated JSON object, expected '}' or key, starting from object at index %d", currentIdx)
			}

			if str[idx] == '}' {
				idx++ // Consume '}'
				return m, idx, nil
			}

			if !firstItem {
				if str[idx] != ',' {
					return nil, currentIdx, fmt.Errorf("expected ',' or '}' in object at index %d, found '%c'. Object started at index %d", idx, safeChar(str, idx), currentIdx)
				}
				idx++ // Consume ','
				idx = localSkipWhitespace(str, idx)
				if idx >= len(str) {
					return nil, currentIdx, fmt.Errorf("unterminated JSON object after comma, starting from object at index %d", currentIdx)
				}
			}
			firstItem = false

			// Parse key
			if str[idx] != '"' {
				return nil, currentIdx, fmt.Errorf("expected '\"' for object key start at index %d, found '%c'. Object started at index %d", idx, safeChar(str, idx), currentIdx)
			}
			keyStart := idx + 1
			keyEnd := keyStart
			for keyEnd < len(str) && str[keyEnd] != '"' {
				keyEnd++
			}
			if keyEnd >= len(str) {
				return nil, currentIdx, fmt.Errorf("unterminated object key string starting at index %d. Object started at index %d", keyStart-1, currentIdx)
			}
			key := str[keyStart:keyEnd]
			idx = keyEnd + 1 // Consume closing '"' for key

			idx = localSkipWhitespace(str, idx)

			if idx >= len(str) || str[idx] != ':' {
				return nil, currentIdx, fmt.Errorf("expected ':' after object key '%s' at index %d. Object started at index %d", key, idx, currentIdx)
			}
			idx++ // Consume ':'

			var val any
			var valEndIdx int
			val, valEndIdx, err = parseArbitraryValueRecursive(str, idx)
			if err != nil {
				return nil, currentIdx, fmt.Errorf("error parsing value for key '%s' in object starting at index %d: %w", key, currentIdx, err)
			}
			m[key] = val
			idx = valEndIdx
		}
	}

	parseArrayToListRecursive = func(str string, currentIdx int) (list []any, nextIdx int, err error) {
		localSkipWhitespace := func(s string, i int) int {
			for i < len(s) && (s[i] == ' ' || s[i] == '\n' || s[i] == '\r' || s[i] == '\t') {
				i++
			}
			return i
		}

		idx := localSkipWhitespace(str, currentIdx)

		if idx >= len(str) || str[idx] != '[' {
			return nil, currentIdx, fmt.Errorf("expected '[' for array start at index %d, but found '%c' (or EOF)", idx, safeChar(str, idx))
		}
		idx++ // Consume '['

		list = make([]any, 0)
		firstElement := true

		for {
			idx = localSkipWhitespace(str, idx)
			if idx >= len(str) {
				return nil, currentIdx, fmt.Errorf("unterminated JSON array, expected ']' or value, starting from array at index %d", currentIdx)
			}

			if str[idx] == ']' {
				idx++ // Consume ']'
				return list, idx, nil
			}

			if !firstElement {
				if str[idx] != ',' {
					return nil, currentIdx, fmt.Errorf("expected ',' or ']' in array at index %d, found '%c'. Array started at index %d", idx, safeChar(str, idx), currentIdx)
				}
				idx++ // Consume ','
				idx = localSkipWhitespace(str, idx)
				if idx >= len(str) {
					return nil, currentIdx, fmt.Errorf("unterminated JSON array after comma, starting from array at index %d", currentIdx)
				}
				if str[idx] == ']' {
					return nil, currentIdx, fmt.Errorf("trailing comma in array at index %d. Array started at index %d", idx-1, currentIdx)
				}
			}
			firstElement = false

			var element any
			var elementEndIdx int
			element, elementEndIdx, err = parseArbitraryValueRecursive(str, idx)
			if err != nil {
				return nil, currentIdx, fmt.Errorf("error parsing element in array starting at index %d: %w", currentIdx, err)
			}
			list = append(list, element)
			idx = elementEndIdx
		}
	}

	parseObjectToMapRecursiveBytes = func(data []byte, currentIdx int) (m map[string]any, nextIdx int, err error) {
		localSkipWhitespaceBytes := func(d []byte, i int) int {
			for i < len(d) && (d[i] == ' ' || d[i] == '\n' || d[i] == '\r' || d[i] == '\t') {
				i++
			}
			return i
		}
		idx := localSkipWhitespaceBytes(data, currentIdx)
		if idx >= len(data) || data[idx] != '{' {
			return nil, currentIdx, fmt.Errorf("bytes: expected {")
		}
		idx++
		m = make(map[string]any)
		firstItem := true
		for {
			idx = localSkipWhitespaceBytes(data, idx)
			if idx >= len(data) {
				return nil, currentIdx, fmt.Errorf("bytes: unterminated object")
			}
			if data[idx] == '}' {
				idx++
				return m, idx, nil
			}
			if !firstItem {
				if data[idx] != ',' {
					return nil, currentIdx, fmt.Errorf("bytes: expected , or }")
				}
				idx++
				idx = localSkipWhitespaceBytes(data, idx)
			}
			if idx >= len(data) {
				return nil, currentIdx, fmt.Errorf("bytes: unterminated object")
			}
			firstItem = false
			if data[idx] != '"' {
				return nil, currentIdx, fmt.Errorf("bytes: expected key quote")
			}
			keyStart := idx + 1
			keyEnd := keyStart
			for keyEnd < len(data) && data[keyEnd] != '"' {
				keyEnd++
			}
			if keyEnd >= len(data) {
				return nil, currentIdx, fmt.Errorf("bytes: unterminated key")
			}
			key := string(data[keyStart:keyEnd])
			idx = keyEnd + 1
			idx = localSkipWhitespaceBytes(data, idx)
			if idx >= len(data) || data[idx] != ':' {
				return nil, currentIdx, fmt.Errorf("bytes: expected colon")
			}
			idx++
			var val any
			var valEndIdx int
			val, valEndIdx, err = parseArbitraryValueRecursiveBytes(data, idx) // Calls BYTES version
			if err != nil {
				return nil, currentIdx, fmt.Errorf("bytes: %w", err)
			}
			m[key] = val
			idx = valEndIdx
		}
	}

	parseArrayToListRecursiveBytes = func(data []byte, currentIdx int) (list []any, nextIdx int, err error) {
		localSkipWhitespaceBytes := func(d []byte, i int) int {
			for i < len(d) && (d[i] == ' ' || d[i] == '\n' || d[i] == '\r' || d[i] == '\t') {
				i++
			}
			return i
		}
		idx := localSkipWhitespaceBytes(data, currentIdx)
		if idx >= len(data) || data[idx] != '[' {
			return nil, currentIdx, fmt.Errorf("bytes: expected [")
		}
		idx++
		list = make([]any, 0)
		firstElement := true
		for {
			idx = localSkipWhitespaceBytes(data, idx)
			if idx >= len(data) {
				return nil, currentIdx, fmt.Errorf("bytes: unterminated array")
			}
			if data[idx] == ']' {
				idx++
				return list, idx, nil
			}
			if !firstElement {
				if data[idx] != ',' {
					return nil, currentIdx, fmt.Errorf("bytes: expected , or ]")
				}
				idx++
				idx = localSkipWhitespaceBytes(data, idx)
			}
			if idx >= len(data) {
				return nil, currentIdx, fmt.Errorf("bytes: unterminated array")
			}
			if data[idx] == ']' {
				return nil, currentIdx, fmt.Errorf("bytes: trailing comma")
			}
			firstElement = false
			var element any
			var elementEndIdx int
			element, elementEndIdx, err = parseArbitraryValueRecursiveBytes(data, idx) // Calls BYTES version
			if err != nil {
				return nil, currentIdx, fmt.Errorf("bytes: %w", err)
			}
			list = append(list, element)
			idx = elementEndIdx
		}
	}
}

// ParseString parses a flat JSON object string into a slice of kstrct.KV.
func ParseString(str string) []KV {
	res := make([]KV, 0, 11) // Updated capacity for new example fields

	var currentKeyString string
	idx := 0

	skipWhitespace := func() {
		for idx < len(str) && (str[idx] == ' ' || str[idx] == '\n' || str[idx] == '\r' || str[idx] == '\t') {
			idx++
		}
	}

	// 1. Expect '{'
	skipWhitespace()
	if idx >= len(str) || str[idx] != '{' {
		fmt.Printf("Error: Expected '{' at start of JSON object, found '%c' at index %d or EOF\n", safeChar(str, idx), idx)
		return nil
	}
	idx++ // Consume '{'

	firstTopLevelItem := true
	for idx < len(str) {
		skipWhitespace()
		if idx >= len(str) {
			fmt.Printf("Error: Unexpected end of string, expecting key or '}' for main object\n")
			return nil
		}

		if str[idx] == '}' {
			idx++ // Consume '}'
			break
		}

		if !firstTopLevelItem {
			if str[idx] != ',' {
				fmt.Printf("Error: Expected ',' or '}' in main object at index %d, found '%c'\n", idx, safeChar(str, idx))
				return nil
			}
			idx++ // Consume ','
			skipWhitespace()
			if idx >= len(str) {
				fmt.Printf("Error: Unexpected end of string after comma in main object\n")
				return nil
			}
			if str[idx] == '}' {
				fmt.Printf("Error: Trailing comma in main object at index %d\n", idx-1)
				return nil
			}
		}
		firstTopLevelItem = false

		// 3. Expect key start '"'
		if str[idx] != '"' {
			fmt.Printf("Error: Expected '\"' for key start in main object at index %d, found '%c'\n", idx, safeChar(str, idx))
			return nil
		}
		idx++ // Consume '"'

		// 4. Parse key
		keyStartIdx := idx
		for idx < len(str) && str[idx] != '"' {
			idx++
		}
		if idx >= len(str) {
			fmt.Printf("Error: Unterminated key string in main object, starting at %d\n", keyStartIdx-1)
			return nil
		}
		currentKeyString = str[keyStartIdx:idx]
		idx++ // Consume closing '"' for key

		skipWhitespace()

		// 5. Expect ':'
		if idx >= len(str) || str[idx] != ':' {
			fmt.Printf("Error: Expected ':' after key '%s' in main object at index %d, found '%c'\n", currentKeyString, idx, safeChar(str, idx))
			return nil
		}
		idx++ // Consume ':'

		var currentValue any
		var errValue error
		currentValue, idx, errValue = parseArbitraryValueRecursive(str, idx)
		if errValue != nil {
			fmt.Printf("Error parsing value for key '%s' in main object: %v\n", currentKeyString, errValue)
			return nil
		}

		res = append(res, KV{Key: currentKeyString, Value: currentValue})

		skipWhitespace()
	}

	return res
}

func safeChar(s string, i int) byte {
	if i < 0 || i >= len(s) {
		return '?'
	}
	return s[i]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func UnmarshalJson(by []byte, dest any) error {
	sl := ParseBytes(by)
	if len(sl) > 0 {
		return Fill(dest, sl)
	}
	return fmt.Errorf("not working")
}
