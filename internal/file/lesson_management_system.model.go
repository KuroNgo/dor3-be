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
	ExplainVie    string
	ExplainEng    string
	ExampleVie    string
	ExampleEng    string
	UnitID        string
}

type Mean struct {
	LessonID     string
	VocabularyID string
	ExplainVie   string
	ExplainEng   string
	ExampleVie   string
	ExampleEng   string
	Synonym      string
	Antonym      string
}

type Final struct {
	CourseName              string
	LessonCourseID          string
	LessonName              string
	LessonContent           string
	LessonLevel             int
	UnitLessonID            string
	UnitName                string
	UnitContent             string
	VocabularyWord          string
	VocabularyPartOfSpeech  string
	VocabularyPronunciation string
	VocabularyExample       string
	VocabularyFieldOfIT     string
	VocabularyUnitID        string
	MeanLessonID            string
	MeanVocabularyID        string
	MeanExplainVie          string
	MeanExplainEng          string
	MeanExampleVie          string
	MeanExampleEng          string
	MeanSynonym             string
	MeanAntonym             string
}
