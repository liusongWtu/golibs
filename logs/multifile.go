package logs

import (
	"encoding/json"
	"time"
)

// A filesLogWriter manages several fileLogWriter
// filesLogWriter will write logs to the file in json configuration  and write the same level log to correspond file
// means if the file name in configuration is project.log filesLogWriter will create project.error.log/project.debug.log
// and write the error-level logs to project.error.log and write the debug-level logs to project.debug.log
// the rotate attribute also  acts like fileLogWriter
type multiFileLogWriter struct {
	writers       [LevelDebug + 1 + 1]*fileLogWriter // the last one for fullLogWriter
	fullLogWriter *fileLogWriter
	Separate      []string `json:"separate"`
}

var levelNames = [...]string{"emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"}

// Init file logger with json config.
// jsonConfig like:
//	{
//	"filename":"logs/test.log",
//	"maxLines":0,
//	"maxsize":0,
//	"daily":true,
//	"maxDays":15,
//	"rotate":true,
//  	"perm":0600,
//	"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],
//	}

func (f *multiFileLogWriter) Init(config string) error {
	writer := newFileWriter().(*fileLogWriter)
	err := writer.Init(config)
	if err != nil {
		return err
	}
	f.fullLogWriter = writer
	f.writers[LevelDebug+1] = writer

	//unmarshal "separate" field to f.Separate
	json.Unmarshal([]byte(config), f)

	jsonMap := map[string]interface{}{}
	json.Unmarshal([]byte(config), &jsonMap)

	for i := LevelEmergency; i < LevelDebug+1; i++ {
		for _, v := range f.Separate {
			if v == levelNames[i] {
				jsonMap["filename"] = f.fullLogWriter.fileNameOnly + "." + levelNames[i] + f.fullLogWriter.suffix
				jsonMap["level"] = i
				bs, _ := json.Marshal(jsonMap)
				writer = newFileWriter().(*fileLogWriter)
				writer.Init(string(bs))
				f.writers[i] = writer
			}
		}
	}

	return nil
}

func (f *multiFileLogWriter) Destroy() {
	for i := 0; i < len(f.writers); i++ {
		if f.writers[i] != nil {
			f.writers[i].Destroy()
		}
	}
}

func (f *multiFileLogWriter) WriteMsg(when time.Time, msg string, level int) error {
	if f.fullLogWriter != nil {
		f.fullLogWriter.WriteMsg(when, msg, level)
	}
	for i := 0; i < len(f.writers)-1; i++ {
		if f.writers[i] != nil {
			if level == f.writers[i].Level {
				f.writers[i].WriteMsg(when, msg, level)
			}
		}
	}
	return nil
}

func (f *multiFileLogWriter) Flush() {
	for i := 0; i < len(f.writers); i++ {
		if f.writers[i] != nil {
			f.writers[i].Flush()
		}
	}
}

// newFilesWriter create a FileLogWriter returning as LoggerInterface.
func newFilesWriter() Logger {
	return &multiFileLogWriter{}
}

func init() {
	Register(AdapterMultiFile, newFilesWriter)
}
