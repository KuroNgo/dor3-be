package file_internal

type Course struct {
	Name        string
	Description string
}

type Lesson struct {
	CourseID string
	Name     string
	Content  string
	Level    int
}

type Unit struct {
	LessonID string
	Name     string
	Content  string
}

type Vocabulary struct {
	Word          string
	PartOfSpeech  string
	Pronunciation string
	Example       string
	FieldOfIT     string
	LinkURL       string
	UnitID        string
}

type Mean struct {
	VocabularyID string
	LessonID     string
	Description  string
	Example      string
	Synonym      string
	Antonym      string
}
