package utils

import (
    "testing"
)

func TestSanitizeText(t *testing.T) {

    // --- TEXT Subtests ---
    t.Run("TEXT input", func(t *testing.T) {
        tests := []struct {
            name     string
            input    string
            expected string
        }{
            {
                name:     "Alphabetic only",
                input:    "HelloWorld",
                // Alpha(..., true) keeps alpha chars,
                // so result should be the same.
                expected: "HelloWorld",
            },
            {
                name:     "Mixed letters, digits, punctuation",
                input:    "Hi123!!!",
                // Digits and punctuation are removed, leaving letters.
                expected: "Hi",
            },
            {
                name:     "Spaces kept, non-alpha removed",
                input:    "Hello 123 World!!!",
                // The digits and punctuation vanish, spaces remain.
                expected: "Hello  World",
            },
            {
                name:     "Empty string",
                input:    "",
                // No characters to sanitize, remains empty.
                expected: "",
            },
            {
                name:     "Symbols only",
                input:    "!@#$%^&*()_+",
                // All non-alpha are removed, leaving nothing.
                expected: "",
            },
        }

        for _, tc := range tests {
            t.Run(tc.name, func(t *testing.T) {
                got := sanitizeText(tc.input, TEXT)
                if got != tc.expected {
                    t.Errorf("TEXT test failed: got %q, want %q", got, tc.expected)
                }
            })
        }
    })

    // --- EMAIL Subtests ---
    t.Run("EMAIL input", func(t *testing.T) {
        tests := []struct {
            name     string
            input    string
            expected string
        }{
            {
                name:     "Valid email",
                input:    "user@example.com",
                // Should remain unchanged if recognized as valid.
                expected: "user@example.com",
            },
            {
                name:     "Common plus address format",
                input:    "My+ALias@Example.Com",
                expected: "my+alias@example.com",
            },
            {
                name:     "Invalid email",
                input:    "not-an-email",
                // The library might return the original or empty; test actual behavior.
                expected: "not-an-email",
            },
            {
                name:     "Empty string",
                input:    "",
                // Typically remains empty.
                expected: "",
            },
        }

        for _, tc := range tests {
            t.Run(tc.name, func(t *testing.T) {
                got := sanitizeText(tc.input, EMAIL)
                if got != tc.expected {
                    t.Errorf("EMAIL test failed: got %q, want %q", got, tc.expected)
                }
            })
        }
    })

    // --- PASSWORD Subtests ---
    t.Run("PASSWORD input", func(t *testing.T) {
        tests := []struct {
            name     string
            input    string
            expected string
        }{
            {
                name:     "SQL injection attempt",
                input:    `'; DROP TABLE users; --`,
                // Strips quotes, semicolons, dashes, etc.
                // Likely leaves: " DROP TABLE users "
                expected: " DROP TABLE users ",
            },
            {
                name:     "No dangerous chars",
                input:    "SafePassword123!",
                // None of the characters in the regex are removed.
                expected: "SafePassword123!",
            },
            {
                name:     "All dangerous chars only",
                input:    `'"-#/*\`,
                // Everything is matched by `[;'"\-#/\*\\]+`,
                // so result is empty.
                expected: "",
            },
            {
                name:     "Mixed dangerous and safe",
                input:    `myPa$$';DROP--table#`,
                // `$` isn't in the regex, so it stays. 
                // Single quote, semicolon, dash, and hash are removed.
                // We might end up with: "myPa$$DROPtable"
                expected: "myPa$$DROPtable",
            },
            {
                name:     "Empty string",
                input:    "",
                // Remains empty.
                expected: "",
            },
        }

        for _, tc := range tests {
            t.Run(tc.name, func(t *testing.T) {
                got := sanitizeText(tc.input, PASSWORD)
                if got != tc.expected {
                    t.Errorf("PASSWORD test failed: got %q, want %q", got, tc.expected)
                }
            })
        }
    })

    // --- Default case (unrecognized InputType) ---
    t.Run("Unrecognized InputType", func(t *testing.T) {
        tests := []struct {
            name     string
            input    string
            inputType InputType
            expected string
        }{
            {
                name:     "Unknown type returns original",
                input:    "SomeInput!!",
                inputType: "UNKNOWN",
                expected: "SomeInput!!",
            },
        }

        for _, tc := range tests {
            t.Run(tc.name, func(t *testing.T) {
                got := sanitizeText(tc.input, tc.inputType)
                if got != tc.expected {
                    t.Errorf("Default test failed: got %q, want %q", got, tc.expected)
                }
            })
        }
    })
}
