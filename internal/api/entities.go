package api

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type Command struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OutputParser   string `json:"output_parser"`
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
