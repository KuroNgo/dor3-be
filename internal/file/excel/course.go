package excel

import (
	"clean-architecture/internal/file"
	"fmt"
	"github.com/xuri/excelize/v2"
)

func ReadFileForCourse(filename string) (file_internal.Course, error) {
	f, err := excelize.OpenFile(filename)

	if err != nil {
		return file_internal.Course{}, err
	}
	defer func() {
		if err := f.Close(); err != nil { // Kiểm tra lỗi khi đóng tệp
			fmt.Printf("Failed to close file: %v\n", err)
		}
	}()

	course := file_internal.Course{
		Name:        "English for IT",
		Description: "",
	}

	return course, nil
}
