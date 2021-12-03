package cmd

type EnvironmentData struct {
	Environment               string
	AWSAccountID              string
}

type WebSocketConfig struct {
	Environments              []string
	ProjectName               string
	Language                  string
	InfrastructureFilePath    string
	WebsocketFilePath         string
	AWSRegion                 string
	AuthorizationKey          string
}

type TerraformConfig struct {
	ENVIRONMENT               string
	AWS_REGION                string
	AWS_ACCOUNT_ID            string
	PROJECT_NAME       string
}

type ServerlessConfig struct {
	AWSAccountID              string
	AuthorizationKey          string
	ProjectName               string
	AWSRegion                 string
}
