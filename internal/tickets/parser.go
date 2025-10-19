package tickets

import (
	"strconv"
	"strings"
)

// ParsePriority parses a priority string and returns an integer value
// Supports Spanish enum values: "Alta", "Media", "Baja"
// Also supports English values: "High", "Medium", "Low", "Urgent", "Critical", "Normal"
func ParsePriority(priorityStr string) int {
	// Default priority
	priority := 1

	if priorityStr == "" {
		return priority
	}

	// Convert to lowercase for case-insensitive comparison
	priorityLower := strings.ToLower(strings.TrimSpace(priorityStr))

	switch priorityLower {
	// Spanish enum values
	case "alta":
		priority = 3
	case "media":
		priority = 2
	case "baja":
		priority = 1
	}

	return priority
}

// ParseCooldown parses a cooldown string and returns the number of seconds
// Supports various formats:
// - "2d" = 2 days
// - "1d" = 1 day
// - "30 min" = 30 minutes
// - "2 hour" = 2 hours
// - "60" = 60 minutes (assumes minutes if no unit)
func ParseCooldown(cooldownStr string) int {
	// Default cooldown: 1 hour in seconds
	cooldown := 3600

	if cooldownStr == "" {
		return cooldown
	}

	// Clean and normalize the string
	cooldownLower := strings.ToLower(strings.TrimSpace(cooldownStr))

	// Handle "2d" format (days)
	if strings.HasSuffix(cooldownLower, "d") {
		if num, err := strconv.Atoi(strings.TrimSuffix(cooldownLower, "d")); err == nil {
			cooldown = num * 24 * 3600 // Convert days to seconds
		}
		return cooldown
	}

	// Handle "30 min" format (minutes)
	if strings.Contains(cooldownLower, "min") {
		// Extract number from string like "30 min" or "30 minutes"
		fields := strings.Fields(cooldownLower)
		if len(fields) > 0 {
			if num, err := strconv.Atoi(fields[0]); err == nil {
				cooldown = num * 60 // Convert minutes to seconds
			}
		}
		return cooldown
	}

	// Handle "2 hour" format (hours)
	if strings.Contains(cooldownLower, "hour") {
		// Extract number from string like "2 hour" or "2 hours"
		fields := strings.Fields(cooldownLower)
		if len(fields) > 0 {
			if num, err := strconv.Atoi(fields[0]); err == nil {
				cooldown = num * 3600 // Convert hours to seconds
			}
		}
		return cooldown
	}

	// Handle "2 sec" format (seconds)
	if strings.Contains(cooldownLower, "sec") {
		// Extract number from string like "30 sec" or "30 seconds"
		fields := strings.Fields(cooldownLower)
		if len(fields) > 0 {
			if num, err := strconv.Atoi(fields[0]); err == nil {
				cooldown = num // Already in seconds
			}
		}
		return cooldown
	}

	// Handle plain numbers (assume minutes)
	if num, err := strconv.Atoi(cooldownLower); err == nil {
		cooldown = num * 60 // Assume minutes
	}

	return cooldown
}

// ParseCooldownToDays parses a cooldown string and returns the number of days
// This is a convenience function for when you need days instead of seconds
func ParseCooldownToDays(cooldownStr string) float64 {
	seconds := ParseCooldown(cooldownStr)
	return float64(seconds) / (24 * 3600) // Convert seconds to days
}

// ParseCooldownToHours parses a cooldown string and returns the number of hours
// This is a convenience function for when you need hours instead of seconds
func ParseCooldownToHours(cooldownStr string) float64 {
	seconds := ParseCooldown(cooldownStr)
	return float64(seconds) / 3600 // Convert seconds to hours
}

// ParseCooldownToMinutes parses a cooldown string and returns the number of minutes
// This is a convenience function for when you need minutes instead of seconds
func ParseCooldownToMinutes(cooldownStr string) float64 {
	seconds := ParseCooldown(cooldownStr)
	return float64(seconds) / 60 // Convert seconds to minutes
}
