package api

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type Command struct {
	ID             int    `json:"id"`
	ShortName      string `json:"short_name"`
	Description    string `json:"description"`
	ExecutableType string `json:"executable_type"`
	ExecutablePath string `json:"executable_path"`
	DockerImage    string `json:"docker_image"`
	ContainerArgs  string `json:"arguments"`
	OutputFileName string `json:"output_filename"`
}

type Commands []Command

type Task struct {
	ID      int    `json:"id"`
	Summary string `json:"summary"`
	Status  string `json:"status"`
}

type Vulnerability struct {
	ID      int    `json:"id"`
	Summary string `json:"summary"`
	Status  string `json:"status"`
	Risk    string `json:"risk"`
}
