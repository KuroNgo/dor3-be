package google

import (
	"fmt"
)

func GenerateStudyReminderNotifications(hour, minute int) string {
	return fmt.Sprintf("%d %d * * *", minute, hour)
}

func GenerateStudyReminder3DaysAbsence(hour, minute int) string {
	return fmt.Sprintf("%d %d * * 3", minute, hour)
}

func GenerateWelcomeToFEIT() string {
	return fmt.Sprintf("20 * * * *")
}
