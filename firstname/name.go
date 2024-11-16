package firstname

type Gender uint8

const (
	GenderUnspecified Gender = iota
	GenderMale        Gender = 1
	GenderFemale      Gender = 2
	GenderNeutral     Gender = 3
)

type Name struct {
	Name   string `db:"name"`
	Gender Gender `db:"gender"`
}

func genderFromInt(gender int) Gender {
	const unspecified, male, female, neutral int = 0, 1, 2, 3

	switch gender {
	case unspecified:
		return GenderUnspecified
	case male:
		return GenderMale
	case female:
		return GenderFemale
	case neutral:
		return GenderNeutral
	default:
		return GenderUnspecified
	}
}
