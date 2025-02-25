package utils

import (
    "github.com/mrz1836/go-sanitize"
)

// InputType is a simple, string-based "enum".
type InputType string

// We define three constants for different input types.
const (
    TEXT     InputType = "TEXT"
    EMAIL    InputType = "EMAIL"
    PASSWORD InputType = "PASSWORD"
)

// DangerousCharsRegex is a regular expression that matches
// characters commonly associated with SQL injection attacks.
const DangerousCharsRegex = `[;'"\-#/\*\\]+`

// sanitizeText applies different sanitization rules depending
// on the provided inputType. Each case uses a helper function
// from go-sanitize.
func sanitizeText(input string, inputType InputType) string {
    switch inputType {

    // For TEXT, we only allow alphabetical characters.
    case "TEXT":
        return sanitize.Alpha(input, true)

    // For EMAIL, we clean and normalize the string into a valid email format.
    case "EMAIL":
        return sanitize.Email(input, false)

    // For PASSWORD, we use a custom regex to remove characters
    // deemed dangerous for SQL injection.
    case "PASSWORD":
        return sanitize.Custom(input, DangerousCharsRegex)

    // If the inputType isn't recognized, we simply return the original string.
    default:
        return input
    }
}
