package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

func ReplaceLastURLID(url string) string {
	re := regexp.MustCompile(`/\d+$`)
	return re.ReplaceAllString(url, "/{id}")
}

func ResizeImage(inputPath string, outputPath string, width uint, height uint) error {
	log.Println(inputPath)
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	resizedImg := resize.Resize(width, height, img, resize.Lanczos3)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output image: %v", err)
	}
	defer outFile.Close()

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(outFile, resizedImg, nil)
		if err != nil {
			return fmt.Errorf("failed to encode jpeg image: %v", err)
		}
	case "png":
		err = png.Encode(outFile, resizedImg)
		if err != nil {
			return fmt.Errorf("failed to encode png image: %v", err)
		}
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}

	return nil
}

func SeparateStringIntoUint64Arr(str string, sep string) []uint64 {
	strArr := strings.Split(str, sep)
	resArr := make([]uint64, len(strArr))
	for i, s := range strArr {
		resArr[i], _ = strconv.ParseUint(s, 10, 64)
	}
	return resArr
}

func GenerateCode(length uint64) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func ReadFileContent(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GenerateStoryName() string {
	currentDate := time.Now().Format("20060102")

	uniqueID := uuid.New().String()

	baseName := "story"

	fileName := fmt.Sprintf("%s_%s_%s.txt", baseName, currentDate, uniqueID)
	return fileName
}

func GenerateStoryFile(fileName, storyContent string) error {
	contentByte := []byte(storyContent)
	err := os.WriteFile("./temp-stories/"+fileName, contentByte, 0644)
	if err != nil {
		return err
	}
	return nil
}

func WriteLogFile() (*os.File, error) {
	year, month, day := time.Now().Date()
	monthStr, dayStr := fmt.Sprintf("%d", month), fmt.Sprintf("%d", day)
	if month < 10 {
		monthStr = "0" + monthStr
	}
	if day < 10 {
		dayStr = "0" + dayStr
	}
	dateFile := fmt.Sprintf("%d%s%s.log", year, monthStr, dayStr)

	logFile, err := os.OpenFile(dateFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return nil, err
	}

	return logFile, nil
}

func ConvertDateString(dateString string) (*time.Time, error) {
	layout := `2006-01-02`
	t, err := time.Parse(layout, dateString)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func RandomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func Join(arr []string, joiner string) string {
	result := ""
	if len(arr) == 0 {
		return result
	}
	for index, r := range arr {
		result += r
		if index != len(arr)-1 {
			result += joiner
		}
	}
	return result
}

func StrToTime(str string) time.Time {
	layout := "15:04:05"
	parsedTime, err := time.Parse(layout, str)
	if err != nil {
		log.Fatal("Failed to parse time:", err)
	}

	now := time.Now()
	return time.Date(
		now.Year(), now.Month(), now.Day(),
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, now.Location(),
	)
}

func CafeNameToSnakeCase(cafeName string) string {
	cafeName = strings.ReplaceAll(cafeName, " ", "_")
	cafeName = strings.ToLower(cafeName)
	return cafeName
}

func GetCafeOpenCloseStatus(startTime string, endTime string) string {
	currTime := time.Now()

	if startTime == "00:00:00" && endTime == "00:00:00" {
		return "OPEN"
	}

	startTimeTime := StrToTime(startTime)
	endTimeTime := StrToTime(endTime)

	if currTime.After(startTimeTime) && currTime.Before(endTimeTime) {
		return "OPEN"
	}

	return "CLOSED"
}

func ImageStorage(imgUri string) string {
	return imgUri
}

func InQuery(ids []uint64) string {
	query := `(`
	for index, id := range ids {
		query += fmt.Sprintf("%d", id)
		if index != len(ids)-1 {
			query += `,`
		}
	}
	query += `)`
	return query
}

func ReadJSONObject[T any](jsonStr string) *T {
	var result T
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		log.Printf("error unmarshaling JSON: %v", err)
		return nil
	}
	return &result
}

func WriteJSONObject[T any](obj *T) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Println(`err json `, err)
		return ""
	}
	return string(bytes)
}

func FirstLastName(name string) (firstName, lastName string) {
	userName := strings.Split(name, " ")
	return userName[0], userName[1]
}

func Retry(ctx context.Context, maxRetries int, fn func(ctx context.Context) error) error {
	numbOfRetry := 0

	for numbOfRetry < maxRetries {
		err := fn(ctx)
		if err == nil {
			return nil
		}

		numbOfRetry++
		log.Printf("Retry %d/%d failed: %v", numbOfRetry, maxRetries, err)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(30 * time.Second):
		}
	}

	log.Println("⚠️ Max retries reached")
	return fmt.Errorf("max retries reached")
}

func MutexLock(ctx context.Context, fn func(ctx context.Context) error) error {
	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	if err := fn(ctx); err != nil {
		return err
	}

	return nil
}
