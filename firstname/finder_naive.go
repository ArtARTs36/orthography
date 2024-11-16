package firstname

import (
	"context"
	"strings"
)

type NaiveGenderFinder struct {
	rules NaiveGenderRules
}

type NaiveGenderRules struct {
	Suffixes []NaiveGenderRulesSuffixes
}

type NaiveGenderRulesSuffixes struct {
	Length  int
	Genders map[string]Gender
}

func NewNaiveGenderFinder(rules NaiveGenderRules) NaiveGenderFinder {
	return NaiveGenderFinder{rules: rules}
}

func (f *NaiveGenderFinder) Find(_ context.Context, names []string) (*FindResult, error) {
	res := &FindResult{
		Found:    map[string]Gender{},
		NotFound: make([]string, 0),
	}

	for _, name := range names {
		g := f.FindOne(name)

		if g == GenderUnspecified {
			res.NotFound = append(res.NotFound, name)
		} else {
			res.Found[name] = g
		}
	}

	return res, nil
}

func (f *NaiveGenderFinder) FindOne(name string) Gender {
	if name == "" {
		return GenderUnspecified
	}

	name = strings.ToLower(name)

	nameLength := len([]rune(name))

	for _, rule := range f.rules.Suffixes {
		if rule.Length > nameLength {
			continue
		}

		suffix := f.getSuffix(name, rule.Length)

		gender, ok := rule.Genders[suffix]
		if ok {
			return gender
		}
	}

	return GenderUnspecified
}

func (f *NaiveGenderFinder) getSuffix(name string, suffixLength int) string {
	nameRunes := []rune(name)

	return string(nameRunes[len(nameRunes)-suffixLength:])
}
