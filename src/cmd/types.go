package cmd

type EnvironmentData struct {
	Environment               string
	AWSAccountID              string
}

type WebsocketConfig struct {
	Environments              []string
	ProjectName               string
	Language                  string
	InfrastructureFilePath    string
	WebsocketFilePath         string
	AwsRegion                 string
	AuthorizationKey          string
}
