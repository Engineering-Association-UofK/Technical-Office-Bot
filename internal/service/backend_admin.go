package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/config"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/models"
)

type credentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AdminAccount struct {
	credentials credentials
	token       string
	Host        string
	client      *http.Client
}

func NewAdminAccount() *AdminAccount {
	return &AdminAccount{
		credentials: credentials{
			Name:     config.App.UserName,
			Password: config.App.Password,
		},
		Host:   config.App.Host,
		client: &http.Client{},
	}
}

func (a *AdminAccount) GetToken() error {
	jsonBody, err := json.Marshal(a.credentials)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", a.Host+"/admin/login", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	token, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		var e models.Error
		if err = json.Unmarshal(token, &e); err != nil {
			return err
		}
		slog.Error("Error trying to log in.", "Status", e.Status, "Message", e.Message, "Time", e.TimeStamp)
		return fmt.Errorf("Failed to get Auth token: ")
	}
	a.token = string(token)
	return nil
}

func (a *AdminAccount) CheckHealth() (*models.ActuatorHealthResponse, *models.Error, error) {
	for {
		req, err := http.NewRequest("GET", a.Host+"/actuator/health", nil)
		if err != nil {
			return nil, nil, err
		}
		req.Header.Add("Authorization", "Bearer "+a.token)

		resp, err := a.client.Do(req)
		if err != nil {
			return nil, nil, err
		}

		statusCode := resp.StatusCode

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, nil, err
		}

		if resp.StatusCode == 200 {
			var health models.ActuatorHealthResponse
			err = json.Unmarshal(body, &health)
			if err != nil {
				return nil, nil, err
			}
			return &health, nil, nil
		}

		var Err models.Error
		if statusCode != 403 {
			if err = json.Unmarshal(body, &Err); err != nil {
				return nil, nil, err
			}
		}

		if Err.Message == "Invalid token." || statusCode == 403 {
			slog.Debug("Invalid token, Getting new one!")
			if a.GetToken() != nil {
				return nil, nil, err
			}
			continue
		}
		return nil, &Err, nil
	}
}
