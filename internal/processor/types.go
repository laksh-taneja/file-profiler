package processor

type FileDimension struct { // input
	Filename   string
	CountLines bool
	CountWords bool
	DoHash     bool
}

type FileResult struct {
	WordCount int
	LineCount int
	CheckSum  string
}
