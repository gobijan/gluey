package codegen

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// titleCaser is a shared instance for converting strings to title case
var titleCaser = cases.Title(language.English)

// ToTitle converts a string to title case using proper Unicode handling
func ToTitle(s string) string {
	return titleCaser.String(s)
}

// ToSingular converts a plural resource name to singular form
func ToSingular(plural string) string {
	if strings.HasSuffix(plural, "ies") {
		return plural[:len(plural)-3] + "y"
	}
	if strings.HasSuffix(plural, "es") {
		return plural[:len(plural)-2]
	}
	if strings.HasSuffix(plural, "s") {
		return plural[:len(plural)-1]
	}
	return plural
}

// ToPlural converts a singular resource name to plural form
func ToPlural(singular string) string {
	if strings.HasSuffix(singular, "y") {
		return singular[:len(singular)-1] + "ies"
	}
	if strings.HasSuffix(singular, "s") || strings.HasSuffix(singular, "x") ||
		strings.HasSuffix(singular, "ch") || strings.HasSuffix(singular, "sh") {
		return singular + "es"
	}
	return singular + "s"
}

// ToCamelCase converts snake_case to CamelCase
func ToCamelCase(snakeCase string) string {
	parts := strings.Split(snakeCase, "_")
	for i, part := range parts {
		if part != "" {
			parts[i] = ToTitle(part)
		}
	}
	return strings.Join(parts, "")
}
