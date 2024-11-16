//nolint:gomnd // not need
package firstname

import (
	"strings"
)

type NaiveGenderDetector struct {
	rules NaiveGenderRules
}

type NaiveGenderRules struct {
	Suffixes []NaiveGenderRulesSuffixes
}

type NaiveGenderRulesSuffixes struct {
	Length  int
	Genders map[string]Gender
}

var RussianNaiveGenderRules = NaiveGenderRules{
	Suffixes: []NaiveGenderRulesSuffixes{
		{
			Length: 5,
			Genders: map[string]Gender{
				"слава": GenderFemale,
			},
		},
		{
			Length: 4,
			Genders: map[string]Gender{
				"слав": GenderMale,
			},
		},
		{
			Length: 3,
			Genders: map[string]Gender{
				"ика": GenderFemale,
				"ида": GenderFemale,
				"ана": GenderFemale,
			},
		},
		{
			Length: 2,
			Genders: map[string]Gender{
				"ан": GenderMale,
				"он": GenderMale,
				"ий": GenderMale,
				"ав": GenderMale,
				"ей": GenderMale,
				"ия": GenderFemale,
				"ья": GenderFemale,
			},
		},
	},
}

func NewNaiveGenderDetector(rules NaiveGenderRules) NaiveGenderDetector {
	return NaiveGenderDetector{rules: rules}
}

func (d *NaiveGenderDetector) Detect(name string) Gender {
	if name == "" {
		return GenderUnspecified
	}

	name = strings.ToLower(name)

	for _, rule := range d.rules.Suffixes {
		if rule.Length > len(name) {
			continue
		}

		suffix := d.getSuffix(name, rule.Length)

		gender, ok := rule.Genders[suffix]
		if ok {
			return gender
		}
	}

	return GenderUnspecified
}

func (d *NaiveGenderDetector) getSuffix(name string, suffixLength int) string {
	nameRunes := []rune(name)

	return string(nameRunes[len(nameRunes)-suffixLength:])
}
