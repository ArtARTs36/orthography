//nolint:gomnd // not need
package firstname

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
				"ем": GenderMale,
				"ём": GenderMale,
				"ей": GenderMale,
				"на": GenderFemale,
				"ил": GenderMale,
				"ия": GenderFemale,
				"ья": GenderFemale,
			},
		},
	},
}
