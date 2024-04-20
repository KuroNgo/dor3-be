package exam_repository

import "sync"

type ExamCountMap struct {
	sync.RWMutex
	CountMap map[string]int
}

func NewExamCountMap() *ExamCountMap {
	return &ExamCountMap{
		CountMap: make(map[string]int),
	}
}

func (m *ExamCountMap) UpdateCount(examID string, count int) {
	m.Lock()
	defer m.Unlock()
	m.CountMap[examID] = count
}

func (m *ExamCountMap) GetCount(examID string) (int, bool) {
	m.RLock()
	defer m.RUnlock()
	count, ok := m.CountMap[examID]
	return count, ok
}
