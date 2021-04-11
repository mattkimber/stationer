package output_file

type Element interface {
	GetLines() []string
	GetComment() string
}

type FileWriter interface {
	WriteToFile(file *File)
}

type SortableFileWriter interface {
	FileWriter
	SetID(int)
	GetID() int
}