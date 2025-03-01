package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
)

// S3WorkerAPI представляет клиент для работы с S3 API
type S3WorkerAPI struct {
	BucketName  string
	S3WorkerURL string
	Client      *http.Client
}

// NewS3WorkerAPI создает новый экземпляр S3WorkerAPI
func NewS3WorkerAPI(bucketName string, s3WorkerURL string) *S3WorkerAPI {
	return &S3WorkerAPI{
		BucketName:  bucketName,
		S3WorkerURL: s3WorkerURL,
		Client:      &http.Client{},
	}
}

// NewBucket создает новый бакет
func (s *S3WorkerAPI) NewBucket() error {
	reqURL, err := url.Parse(fmt.Sprintf("%s/new_bucket", s.S3WorkerURL))
	if err != nil {
		return fmt.Errorf("ошибка при формировании URL: %w", err)
	}

	q := reqURL.Query()
	q.Add("name", s.BucketName)
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodPost, reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка при создании бакета: %s, код: %d", string(body), resp.StatusCode)
	}

	return nil
}

// RemoveBucket удаляет бакет
func (s *S3WorkerAPI) RemoveBucket() error {
	reqURL, err := url.Parse(fmt.Sprintf("%s/remove_bucket", s.S3WorkerURL))
	if err != nil {
		return fmt.Errorf("ошибка при формировании URL: %w", err)
	}

	q := reqURL.Query()
	q.Add("name", s.BucketName)
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodDelete, reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка при удалении бакета: %s, код: %d", string(body), resp.StatusCode)
	}

	return nil
}

// ListFiles возвращает список файлов в бакете
func (s *S3WorkerAPI) ListFiles() ([]string, error) {
	reqURL := fmt.Sprintf("%s/files/%s", s.S3WorkerURL, s.BucketName)
	
	resp, err := s.Client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ошибка при получении списка файлов: %s, код: %d", string(body), resp.StatusCode)
	}

	var files []string
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	return files, nil
}

// GetFilesURLs возвращает URL файлов в бакете
func (s *S3WorkerAPI) GetFilesURLs(filenames []string) (map[string]string, error) {
	reqURL := fmt.Sprintf("%s/files/%s", s.S3WorkerURL, s.BucketName)
	
	body, err := json.Marshal(filenames)
	if err != nil {
		return nil, fmt.Errorf("ошибка при маршалинге данных: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ошибка при получении URL файлов: %s, код: %d", string(body), resp.StatusCode)
	}

	var urls map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&urls); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	return urls, nil
}

// UploadFile загружает файл в бакет
func (s *S3WorkerAPI) UploadFile(fileData []byte, filename string, mimetype string) (string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	
	// Создаем часть для файла
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании формы для файла: %w", err)
	}
	
	// Записываем данные файла
	if _, err = fw.Write(fileData); err != nil {
		return "", fmt.Errorf("ошибка при записи данных файла: %w", err)
	}
	
	// Закрываем writer
	w.Close()

	// Формируем URL с параметрами
	reqURL, err := url.Parse(fmt.Sprintf("%s/upload_file", s.S3WorkerURL))
	if err != nil {
		return "", fmt.Errorf("ошибка при формировании URL: %w", err)
	}

	q := reqURL.Query()
	q.Add("bucket", s.BucketName)
	if filename != "" {
		q.Add("filename", filename)
	}
	if mimetype != "" {
		q.Add("mimetype", mimetype)
	}
	reqURL.RawQuery = q.Encode()

	// Создаем запрос
	req, err := http.NewRequest(http.MethodPost, reqURL.String(), &b)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Выполняем запрос
	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка при загрузке файла: %s, код: %d", string(body), resp.StatusCode)
	}

	// Читаем ответ (имя файла)
	var resultFilename string
	if err := json.NewDecoder(resp.Body).Decode(&resultFilename); err != nil {
		return "", fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	return resultFilename, nil
}

// UploadFileFromMultipart загружает файл из multipart.FileHeader в бакет
func (s *S3WorkerAPI) UploadFileFromMultipart(file *multipart.FileHeader) (string, error) {
	// Открываем файл
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("ошибка при открытии файла: %w", err)
	}
	defer src.Close()

	// Читаем данные файла
	fileData, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	// Определяем MIME-тип на основе расширения файла
	ext := filepath.Ext(file.Filename)
	var mimetype string
	switch ext {
	case ".jpg", ".jpeg":
		mimetype = "image/jpeg"
	case ".png":
		mimetype = "image/png"
	case ".gif":
		mimetype = "image/gif"
	case ".webp":
		mimetype = "image/webp"
	default:
		mimetype = "application/octet-stream"
	}

	// Загружаем файл
	return s.UploadFile(fileData, file.Filename, mimetype)
}

// RemoveFile удаляет файл из бакета
func (s *S3WorkerAPI) RemoveFile(filename string) error {
	reqURL, err := url.Parse(fmt.Sprintf("%s/remove_file", s.S3WorkerURL))
	if err != nil {
		return fmt.Errorf("ошибка при формировании URL: %w", err)
	}

	q := reqURL.Query()
	q.Add("bucket", s.BucketName)
	q.Add("filename", filename)
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodDelete, reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка при удалении файла: %s, код: %d", string(body), resp.StatusCode)
	}

	return nil
}

// GetFileURL возвращает URL файла
func (s *S3WorkerAPI) GetFileURL(filename string) (string, error) {
	reqURL := fmt.Sprintf("%s/file/%s/%s", s.S3WorkerURL, s.BucketName, filename)
	
	resp, err := s.Client.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка при получении URL файла: %s, код: %d", string(body), resp.StatusCode)
	}

	var fileURL string
	if err := json.NewDecoder(resp.Body).Decode(&fileURL); err != nil {
		return "", fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	return fileURL, nil
}
