package models

type ProfileCard struct {
	Label        *ProfileLabel `json:"label"`
	Job          string        `json:"job"`
	InterestTags []string      `json:"interestTags"`
	SkillTags    []string      `json:"skillTags"`
	IsSubTarget  bool          `json:"isSubTarget"`
}
