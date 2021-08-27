package configurations

// PathConfiguration - Struct for Configurations
type PathConfiguration struct {
	RawCorpusRoot     string
	CleanedCorpusRoot string
	ArchivePath       string
	ArticlePath       string
}

// GetPathConfigurations - Retrieves the email data
func GetPathConfigurations() PathConfiguration {
	pathconfiguration := PathConfiguration{}

	// Change below 2 lines as per your directory structure and
	// to store the corpus data
	pathconfiguration.RawCorpusRoot = "/media/user/corpus/raw_corpus"
	pathconfiguration.CleanedCorpusRoot = "/media/user/corpus/cleaned_corpus"
	pathconfiguration.ArchivePath = "archive"
	pathconfiguration.ArticlePath = "articles"

	return pathconfiguration
}
