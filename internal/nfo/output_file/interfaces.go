package output_file

type Element interface {
	GetLines() []string
	GetComment() string
}

type FileWriter interface {
	WriteToFile(file *File)
}
