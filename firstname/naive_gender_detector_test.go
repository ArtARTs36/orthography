package firstname_test

import (
	"github.com/artarts36/orthography/firstname"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNaiveGenderDetector(t *testing.T) {
	cases := []struct {
		Title     string
		FirstName string
		Gender    firstname.Gender
	}{
		{
			Title:     "empty string",
			FirstName: "",
			Gender:    firstname.GenderUnspecified,
		},
		{
			Title:     "М - Артемий",
			FirstName: "Артемий",
			Gender:    firstname.GenderMale,
		},
		{
			Title:     "Ж - Лилия",
			FirstName: "Лилия",
			Gender:    firstname.GenderFemale,
		},
		{
			Title:     "Ж - Агафья",
			FirstName: "Агафья",
			Gender:    firstname.GenderFemale,
		},
	}

	detector := firstname.NewNaiveGenderDetector(firstname.RussianNaiveGenderRules)

	for _, c := range cases {
		t.Run(c.Title, func(t *testing.T) {
			assert.Equal(t, c.Gender, detector.Detect(c.FirstName))
		})
	}
}
