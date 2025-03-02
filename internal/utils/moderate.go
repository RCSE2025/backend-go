package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RCSE2025/backend-go/internal/config"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type ModeratorAPI struct {
	Client *http.Client
}

func NewModeratorAPI() *ModeratorAPI {
	return &ModeratorAPI{
		Client: &http.Client{},
	}
}

type ModerateRequest struct {
	Files []byte `json:"files"`
	Text  string `json:"text"`
}

func (m *ModeratorAPI) IsModerateContent(content string, files *[]*multipart.FileHeader, staticFile bool) (bool, error) {
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Добавляем поле text
	err := writer.WriteField("text", content)
	if err != nil {
		return false, err
	}

	// Добавляем файлы из списка
	if files != nil {
		for _, f := range *files {
			src, err := f.Open()
			if err != nil {
				return false, err
			}
			defer src.Close()

			// Создаем новый параметр для файла в multipart
			part, err := writer.CreateFormFile("files", f.Filename)
			if err != nil {
				return false, err
			}

			// Копируем данные файла в тело запроса
			_, err = io.Copy(part, src)
			if err != nil {
				return false, err
			}
		}
	}
	// Загрузка статического файла
	if staticFile == true {
		staticFile, err := os.Open("internal/static/pixel.jpg")
		if err != nil {
			return false, fmt.Errorf("unable to open static file: %v", err)
		}
		defer staticFile.Close()

		// Извлекаем имя файла и тип контента
		staticFileName := filepath.Base("internal/static/pixel.jpg")
		// Получаем расширение файла, чтобы корректно указать MIME-тип (например, image/png)
		//ext := filepath.Ext(staticFileName)
		//var mimeType string
		//switch ext {
		//case ".png":
		//	mimeType = "image/png"
		//case ".jpg", ".jpeg":
		//	mimeType = "image/jpeg"
		//case ".gif":
		//	mimeType = "image/gif"
		//default:
		//	mimeType = "application/octet-stream"
		//}

		// Создаем новый параметр для статического файла в multipart
		part, err := writer.CreateFormFile("files", staticFileName)
		if err != nil {
			return false, fmt.Errorf("unable to create form file part for static file: %v", err)
		}

		// Устанавливаем MIME-тип для файла
		//part.Header.Set("Content-Type", mimeType)

		// Копируем данные статического файла в тело запроса
		_, err = io.Copy(part, staticFile)
		if err != nil {
			return false, fmt.Errorf("unable to copy static file data: %v", err)
		}
	}

	// Закрываем writer, чтобы завершить формирование multipart данных
	err = writer.Close()
	if err != nil {
		return false, err
	}

	// Создаем новый HTTP-запрос
	req, err := http.NewRequest(http.MethodPost, config.Get().ModerateModelURL+"/moderate", &b)
	if err != nil {
		return false, err
	}

	// Устанавливаем заголовки для multipart запроса
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Отправляем запрос через HTTP клиент
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Проверка на успешный ответ
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	var r map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&r)

	t, ok := r["result"].(float64)
	if ok == true && t == 0 {
		return true, nil
	}
	return false, nil
}
