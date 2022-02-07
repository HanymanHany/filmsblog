package documents

type PeopleDocument struct {
	Idimg      string `bson:"_id,omitempty"`
	Profession string
	Fio        string
	Image      string
}
