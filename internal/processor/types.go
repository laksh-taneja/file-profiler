package processor

type FileDimension struct { // input
	Filename   string
	CountLines bool
	CountWords bool
	DoHash     bool
	LongLines  bool // > for line 64kb, change buffer size
}

type FileResult struct {
	Filename  string
	WordCount int
	LineCount int
	CheckSum  string
}
