package cmd

type EnvironmentData struct {
	AwsAccountId              string `json:"awsAccountId"`
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
