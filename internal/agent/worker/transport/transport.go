package agent

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/klimenkokayot/calc-net-go/internal/shared/models"
)

func GetTask(client *http.Client, url string) (*models.Task, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, ErrNotFound
	case http.StatusInternalServerError:
		return nil, ErrInternalServer
	case http.StatusOK:
		break
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	task := &models.Task{}
	err = json.Unmarshal(body, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func PostTask(client *http.Client, url string, result *models.TaskResult) error {
	data, _ := json.Marshal(result)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	return err
}
