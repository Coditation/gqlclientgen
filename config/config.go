package config

const (
	ConfigName              = "config"
	ConfigType              = "yaml"
	PackageNameKey          = "packageName"
	OutputDirectoryKey      = "outputDirectory"
	SourceTypeKey           = "sourceType"
	GraphqlServerBaseUrlKey = "graphqlServerBaseUrl"
	SourceFilePathKey       = "sourceFilePath"
	FileSourceType          = "file"
	UrlSourceType           = "remote"
	RemoteURL               = "url"
)

type GqlClientGenConfig struct {
	OutputDirectory      string
	PackageName          string
	SourceType           string
	GraphQLServerBaseUrl string
	SourceFilePath       string
}
