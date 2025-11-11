package userutil

import (
	"strings"
	"unicode"
)

const (
	ethiopiaCountryCode = "251"
	localNumberLength   = 9
)

// NormalizePhone converts various local/E.164 representations of an Ethiopian phone
// number into a canonical digits-only string prefixed with the country code (e.g. 2519...).
// An empty string is returned when no digits are found.
func NormalizePhone(phone string) string {
	if phone == "" {
		return ""
	}

	digits := extractDigits(phone)
	if digits == "" {
		return ""
	}

	// Remove leading zeros if present
	for len(digits) > localNumberLength && digits[0] == '0' {
		digits = digits[1:]
	}

	switch {
	case len(digits) >= localNumberLength && strings.HasPrefix(digits, ethiopiaCountryCode):
		// Keep only the last 9 digits for the local part to avoid duplicated prefixes
		local := digits[len(digits)-localNumberLength:]
		return ethiopiaCountryCode + local
	case len(digits) == localNumberLength && digits[0] == '9':
		return ethiopiaCountryCode + digits
	case len(digits) == localNumberLength+1 && digits[0] == '0':
		return ethiopiaCountryCode + digits[1:]
	case len(digits) == localNumberLength:
		return ethiopiaCountryCode + digits
	default:
		return digits
	}
}

// FormatPhoneE164 returns the phone in +251... format if it can be normalised,
// otherwise returns an empty string.
func FormatPhoneE164(phone string) string {
	normalized := NormalizePhone(phone)
	if normalized == "" {
		return ""
	}
	return "+" + normalized
}

// PhoneVariants returns a set of possible phone representations covering local and
// international patterns. The canonical representation is provided first.
func PhoneVariants(phone string) []string {
	normalized := NormalizePhone(phone)
	if normalized == "" {
		return nil
	}

	raw := strings.TrimSpace(phone)
	digits := extractDigits(phone)
	local := ""
	if len(normalized) > len(ethiopiaCountryCode) {
		local = "0" + normalized[len(ethiopiaCountryCode):]
	}

	add := func(slice *[]string, seen map[string]struct{}, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		*slice = append(*slice, value)
		seen[value] = struct{}{}
	}

	variants := make([]string, 0, 6)
	seen := make(map[string]struct{})

	add(&variants, seen, normalized)
	add(&variants, seen, "+"+normalized)
	add(&variants, seen, digits)
	add(&variants, seen, raw)
	add(&variants, seen, local)

	return variants
}

func extractDigits(phone string) string {
	var builder strings.Builder
	for _, r := range phone {
		if unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
