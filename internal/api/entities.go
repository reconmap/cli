package api

type Command struct {
	ShortName      string `json:"short_name"`
	DockerImage    string `json:"docker_image"`
	ContainerArgs  string `json:"arguments"`
	OutputFileName string `json:"output_filename"`
}
