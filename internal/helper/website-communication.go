package helper

import "errors"

// ConvertPosition converts positions from 'top-left' to 'NW' and vice versa.
func ConvertPosition(position string) (string, error) {
	positionMap := map[string]string{
		"top-left":     "NW",
		"top-right":    "NE",
		"bottom-left":  "SW",
		"bottom-right": "SE",
		"NW":           "top-left",
		"NE":           "top-right",
		"SW":           "bottom-left",
		"SE":           "bottom-right",
	}

	if converted, exists := positionMap[position]; exists {
		return converted, nil
	}
	return "", errors.New("invalid position")
}

// ValidateVersion checks if the version is one of 1, 2, 3, or 'D'.
func ValidateVersion(version string) error {
	validVersions := map[string]bool{
		"1": true,
		"2": true,
		"3": true,
		"D": true,
	}

	if _, exists := validVersions[version]; exists {
		return nil
	}
	return errors.New("invalid version")
}
