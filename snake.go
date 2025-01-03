package kstrct

import (
	"sync"
)

const underscoreByte = '_'

var (
	snakeCache sync.Map
	// Pool for byte slices used in snake case conversion
	bytePool = sync.Pool{
		New: func() any {
			// Pre-allocate a reasonably sized buffer
			b := make([]byte, 0, 64)
			return &b
		},
	}
	// String intern pool to reduce string allocations
	stringPool = sync.Pool{
		New: func() any {
			return new(string)
		},
	}
)

// internString returns a pooled string to reduce allocations
func internString(b []byte) string {
	strPtr := stringPool.Get().(*string)
	*strPtr = string(b)
	s := *strPtr
	stringPool.Put(strPtr)
	return s
}

// ToSnakeCase converts a string to snake case with caching
func ToSnakeCase(s string) string {
	// Check cache first
	if cached, ok := snakeCache.Load(s); ok {
		return cached.(string)
	}

	// Estimate final size: length + potential underscores
	// Count uppercase letters to estimate underscore count
	upperCount := 0
	for i := 0; i < len(s); i++ {
		if isUpper(s[i]) {
			upperCount++
		}
	}

	// Get byte slice from pool
	bPtr := bytePool.Get().(*[]byte)
	b := *bPtr
	b = b[:0]
	// Pre-allocate with exact size estimate
	if cap(b) < len(s)+upperCount {
		*bPtr = make([]byte, 0, len(s)+upperCount)
		b = *bPtr
	}
	defer bytePool.Put(bPtr)

	idx := 0
	hasLower := false
	hasUnderscore := false
	lowercaseSinceUnderscore := false

	// First pass: check if we need any changes
	for ; idx < len(s); idx++ {
		if isLower(s[idx]) {
			hasLower = true
			if hasUnderscore {
				lowercaseSinceUnderscore = true
			}
			continue
		} else if isDigit(s[idx]) {
			continue
		} else if s[idx] == underscoreByte && idx > 0 && idx < len(s)-1 && (isLower(s[idx+1]) || isDigit(s[idx+1])) {
			hasUnderscore = true
			lowercaseSinceUnderscore = false
			continue
		}
		break
	}

	if idx == len(s) {
		snakeCache.Store(s, s) // Cache the result
		return s               // no changes needed, can just borrow the string
	}

	b = append(b, s[:idx]...)

	if isUpper(s[idx]) && (!hasLower || hasUnderscore && !lowercaseSinceUnderscore) {
		for idx < len(s) && (isUpper(s[idx]) || isDigit(s[idx])) {
			b = append(b, asciiLowercaseArray[s[idx]])
			idx++
		}

		for idx < len(s) && (isLower(s[idx]) || isDigit(s[idx])) {
			b = append(b, s[idx])
			idx++
		}
	}

	for idx < len(s) {
		if !isAlphanumeric(s[idx]) {
			idx++
			continue
		}

		if len(b) > 0 {
			b = append(b, underscoreByte)
		}

		for idx < len(s) && (isUpper(s[idx]) || isDigit(s[idx])) {
			b = append(b, asciiLowercaseArray[s[idx]])
			idx++
		}

		for idx < len(s) && (isLower(s[idx]) || isDigit(s[idx])) {
			b = append(b, s[idx])
			idx++
		}
	}

	result := internString(b)
	snakeCache.Store(s, result) // Cache the result
	return result
}

func SnakeCaseToTitle(s string) string {
	b := make([]byte, 0, 64)
	l := len(s)
	i := 0
	for i < l {

		// skip leading bytes that aren't letters or digits
		for i < l && !isWord(s[i]) {
			i++
		}

		// set the first byte to uppercase if it needs to
		if i < l {
			c := s[i]

			// simply append contiguous digits
			if isDigit(c) {
				for i < l {
					if c = s[i]; !isDigit(c) {
						break
					}
					b = append(b, c)
					i++
				}
				continue
			}

			// the sequence starts with and uppercase letter, we append
			// all following uppercase letters as equivalent lowercases
			if isUpper(c) {
				b = append(b, c)
				i++

				for i < l {
					if c = s[i]; !isUpper(c) {
						break
					}
					b = append(b, toLower(c))
					i++
				}

			} else {
				b = append(b, toUpper(c))
				i++
			}

			// append all trailing lowercase letters
			for i < l {
				if c = s[i]; !isLower(c) {
					break
				}
				b = append(b, c)
				i++
			}
		}
	}

	return string(b)
}

func isAlphanumeric(c byte) bool {
	return isLower(c) || isUpper(c) || isDigit(c)
}

func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func isLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isWord(c byte) bool {
	return isLetter(c) || isDigit(c)
}

func isLetter(c byte) bool {
	return isLower(c) || isUpper(c)
}

func toLower(c byte) byte {
	if isUpper(c) {
		return c + ('a' - 'A')
	}
	return c
}

func toUpper(c byte) byte {
	if isLower(c) {
		return c - ('a' - 'A')
	}
	return c
}

var asciiLowercaseArray = [256]byte{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
	0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
	0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	' ', '!', '"', '#', '$', '%', '&', '\'',
	'(', ')', '*', '+', ',', '-', '.', '/',
	'0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', ':', ';', '<', '=', '>', '?',
	'@',

	'a', 'b', 'c', 'd', 'e', 'f', 'g',
	'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w',
	'x', 'y', 'z',

	'[', '\\', ']', '^', '_',
	'`', 'a', 'b', 'c', 'd', 'e', 'f', 'g',
	'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w',
	'x', 'y', 'z', '{', '|', '}', '~', 0x7f,
	0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87,
	0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
	0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97,
	0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
	0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7,
	0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf,
	0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7,
	0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf,
	0xc0, 0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7,
	0xc8, 0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf,
	0xd0, 0xd1, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7,
	0xd8, 0xd9, 0xda, 0xdb, 0xdc, 0xdd, 0xde, 0xdf,
	0xe0, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7,
	0xe8, 0xe9, 0xea, 0xeb, 0xec, 0xed, 0xee, 0xef,
	0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7,
	0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff,
}
