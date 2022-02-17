package model

type File struct {
	Name      string
	XxHash    string
	Ext       string
	Dir       string
	TempDir   string
	Path      string
	RePath    string
	MediaInfo struct {
		Video           FFProbeStream
		Format          FFProbeFormat
		DurationSeconds int64
		DurationFormat  string
	}
}
