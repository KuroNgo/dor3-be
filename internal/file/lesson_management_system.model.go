package file_internal

type Course struct {
	Name        string
	Description string
}

type Lesson struct {
	CourseID string
	Name     string
	Level    int
}

type Unit struct {
	LessonID string
	Name     string
	Level    int
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
	UnitLevel     int
}

type Final struct {
	CourseName              string
	LessonCourseID          string
	LessonName              string
	LessonLevel             int
	UnitLessonID            string
	UnitName                string
	VocabularyWord          string
	VocabularyPartOfSpeech  string
	VocabularyPronunciation string
	VocabularyExample       string
	VocabularyFieldOfIT     string
	VocabularyUnitLevel     int
	MeanLessonID            string
	MeanVocabularyID        string
	MeanExplainVie          string
	MeanExplainEng          string
	MeanExampleVie          string
	MeanExampleEng          string
	MeanSynonym             string
	MeanAntonym             string
}
