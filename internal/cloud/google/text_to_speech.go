package google

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"io"
	"os"
	"path/filepath"
)

func GenerateMP3File() error {
	// Khởi tạo phiên AWS SDK
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Tạo phiên Polly
	svc := polly.New(sess)

	// Tạo yêu cầu chuyển đổi văn bản thành giọng nói
	input := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		Text:         aws.String("Hello, this is a test."),
		VoiceId:      aws.String("Joanna"), // Chọn giọng đọc
	}

	// Gửi yêu cầu đến Polly
	output, err := svc.SynthesizeSpeech(input)
	if err != nil {
		return err
	}

	// Đường dẫn đến thư mục nội bộ
	internalFolderPath := "internal"

	// Tạo đường dẫn đầy đủ của file MP3 trong thư mục nội bộ
	internalFilePath := filepath.Join(internalFolderPath, "output.mp3")

	// Tạo file MP3 trong thư mục nội bộ
	file, err := os.Create(internalFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy dữ liệu âm thanh từ phản hồi của Polly vào file MP3
	files, err := io.Copy(file, output.AudioStream)
	if err != nil {
		return err
	}

	fmt.Println("File MP3 đã được tạo và gắn vào thư mục nội bộ thành công.", files)
	return nil
}
