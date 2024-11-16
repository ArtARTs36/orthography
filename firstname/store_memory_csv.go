package firstname

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadCSVStore(filepath string) (Store, error) {
	const needRowLength = 2

	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", filepath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse csv file %q: %w", filepath, err)
	}

	mapGenders := func(records [][]string) (map[string]Gender, error) {
		result := make(map[string]Gender)

		for rowIndex, row := range records {
			if len(row) < needRowLength {
				return nil, fmt.Errorf("row[%d] must contain 2 or greather items, got %d", rowIndex, len(row))
			}

			name := strings.Trim(row[0], " ")
			gender := strings.Trim(row[1], " ")

			genderInt, gErr := strconv.Atoi(gender)
			if gErr != nil {
				return nil, fmt.Errorf("row[%d]: failed to parse gender: %w", rowIndex, gErr)
			}

			result[strings.ToLower(name)] = genderFromInt(genderInt)
		}

		return result, nil
	}

	names, err := mapGenders(records)
	if err != nil {
		return nil, fmt.Errorf("failed to map data from csv file %q: %w", filepath, err)
	}

	return NewMemoryStore(names), nil
}

// LoadStorePerGender expects file paths to file, which contain names per line.
func LoadStorePerGender(maleFilepath string, femaleFilepath string, otherFilepath string) (Store, error) {
	names := make(map[string]Gender)

	load := func(filepath string, gender Gender) error {
		file, err := os.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("failed to open file %q: %w", filepath, err)
		}
		fileContent := string(file)

		for _, row := range strings.Split(fileContent, "\n") {
			name := strings.TrimSpace(row)
			if name == "" {
				continue
			}

			names[strings.ToLower(name)] = gender
		}

		return nil
	}

	if maleFilepath != "" {
		err := load(maleFilepath, GenderMale)
		if err != nil {
			return nil, fmt.Errorf("failed to load file for male: %w", err)
		}
	}
	if femaleFilepath != "" {
		err := load(femaleFilepath, GenderFemale)
		if err != nil {
			return nil, fmt.Errorf("failed to load file for female: %w", err)
		}
	}
	if otherFilepath != "" {
		err := load(otherFilepath, GenderNeutral)
		if err != nil {
			return nil, fmt.Errorf("failed to load file for other: %w", err)
		}
	}

	return NewMemoryStore(names), nil
}
