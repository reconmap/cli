package api

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type Command struct {
	ID             int    `json:"id"`
	ShortName      string `json:"short_name"`
	DockerImage    string `json:"docker_image"`
	ContainerArgs  string `json:"arguments"`
	OutputFileName string `json:"output_filename"`
}

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
