package models

type People struct {
	Idimg      string
	Profession string
	Fio        string
	Image      string
}

func NewPeople(idimg, profession, fio, image string) *People {
	return &People{idimg, profession, fio, image}

}
