/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package utils

import (
	"strings"
	"unicode"
)

func NormalizeUserId(userId string) string {

	var out strings.Builder

	for _, ch := range []byte(userId) {

		if ch >= '0' && ch <= '9' {
			out.WriteByte(ch)
			continue
		}

		if ch >= 'a' && ch <= 'z' {
			out.WriteByte(ch)
			continue
		}

		if ch >= 'A' && ch <= 'Z' {
			out.WriteByte(ch)
			continue
		}

	}

	return out.String()
}

func NormalizeLogin(login string) string {
	var out strings.Builder
	for _, ch := range login {
		low := unicode.ToLower(ch)
		if low == ':' || low == ' ' {
			continue
		}
		out.WriteRune(low)
	}
	return out.String()
}

func NormalizeEmail(email string) string {
	return NormalizeLogin(email)
}

func NormalizeUsername(username string) string {
	return NormalizeLowerUnreservedCharacters(username)
}

func NormalizeField(field string) string {
	s := strings.TrimSpace(field)
	return strings.ReplaceAll(s, ":", "")
}

func NormalizeCode(code string) string {
	s := strings.TrimSpace(code)
	return strings.ReplaceAll(s, ":", "")
}

func NormalizePageId(pageId string) string {
	return NormalizeLowerUnreservedCharacters(pageId)
}

// RFC 3986 section 2.3 Unreserved Characters (January 2005)
func NormalizeUnreservedCharacters(identity string) string {

	var out strings.Builder

	for _, low := range identity {

		if (low >= 'a' && low <= 'z') ||
			(low >= 'A' && low <= 'Z') ||
			(low >= '0' && low <= '9') ||
			(low == '-' || low == '_' || low == '.' || low == '~' ) {
			out.WriteRune(low)
		}

	}

	return out.String()
}

// RFC 3986 section 2.3 Unreserved Characters (January 2005)
func NormalizeLowerUnreservedCharacters(identity string) string {

	var out strings.Builder

	for _, ch := range identity {

		low := unicode.ToLower(ch)

		if (low >= 'a' && low <= 'z') ||
			(low >= '0' && low <= '9') ||
			(low == '-' || low == '_' || low == '.' || low == '~' ) {
			out.WriteRune(low)
		}

	}

	return out.String()
}