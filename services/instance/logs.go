package instance

const (
	LogKindOut = "out"
	LogKindErr = "err"
)

type LogLine struct {
	Id      int    `json:"id"`
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

// TODO: Use a better structure than array

type Logs struct {
	Lines []LogLine `json:"lines"`
}

func (l *Logs) Add(line LogLine) LogLine {
	line.Id = len(l.Lines) + 1
	l.Lines = append(l.Lines, line)
	return line
}
