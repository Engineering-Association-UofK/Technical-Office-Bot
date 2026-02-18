package models

type ActuatorHealthResponse struct {
	Status     string             `json:"status"`
	Components ActuatorComponents `json:"components"`
	Groups     []string           `json:"groups"`
}

type ActuatorComponents struct {
	DB             DbComponent        `json:"db"`
	DiskSpace      DiskSpaceComponent `json:"diskSpace"`
	LivenessState  SimpleComponent    `json:"livenessState"`
	Mail           MailComponent      `json:"mail"`
	Ping           SimpleComponent    `json:"ping"`
	ReadinessState SimpleComponent    `json:"readinessState"`
	SSL            SSLComponent       `json:"ssl"`
}

type SimpleComponent struct {
	Status string `json:"status"`
}

type DbComponent struct {
	Status  string `json:"status"`
	Details struct {
		Database        string `json:"database"`
		ValidationQuery string `json:"validationQuery"`
	} `json:"details"`
}

type DiskSpaceComponent struct {
	Status  string `json:"status"`
	Details struct {
		Total     uint64 `json:"total"`
		Free      uint64 `json:"free"`
		Threshold uint64 `json:"threshold"`
		Exists    bool   `json:"exists"`
	} `json:"details"`
}

type MailComponent struct {
	Status  string `json:"status"`
	Details struct {
		Location string `json:"location"`
	} `json:"details"`
}

type SSLComponent struct {
	Status  string `json:"status"`
	Details struct {
		ValidChains []string `json:"validChains"`
	} `json:"details"`
}
