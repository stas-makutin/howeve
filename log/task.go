package log

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stas-makutin/howeve/api"
	"github.com/stas-makutin/howeve/config"
	"github.com/stas-makutin/howeve/tasks"
	"github.com/stas-makutin/howeve/utils"
)

// Task struct
type Task struct {
	cfg                *api.LogConfig
	defaultLogFileName string
	fileName           string
	fileMode           os.FileMode
	maxSizeBytes       int64
	maxAgeDuration     time.Duration
	archive            bool
	rotateLock         uint32
}

// NewTask func
func NewTask(defaultLogFileName string) *Task {
	t := &Task{defaultLogFileName: defaultLogFileName}
	config.AddReader(t.readConfig)
	config.AddWriter(t.writeConfig)
	return t
}

func (t *Task) readConfig(cfg *api.Config, cfgError config.Error) {
	t.cfg = cfg.Log
	t.fileName = ""
	if t.cfg == nil {
		return
	}

	if t.cfg.Dir == "" {
		cfgError("log.dir is required")
	}
	t.fileName = t.cfg.File
	if t.fileName == "" {
		t.fileName = t.defaultLogFileName + ".log"
	}
	dirMode := t.cfg.DirMode.WithDirDefault()
	t.fileMode = t.cfg.FileMode.WithFileDefault()
	err := os.MkdirAll(t.cfg.Dir, dirMode)
	if err != nil {
		cfgError("log.dir is not valid")
	}

	size := t.cfg.MaxSize.Value()
	if size < 0 {
		err = fmt.Errorf("negative value not allowed")
	}
	if err != nil {
		cfgError(fmt.Sprintf("log.maxSize is not valid: %v", err))
	}
	t.maxSizeBytes = size

	duration := t.cfg.MaxAge.Value()
	if duration < 0 {
		err = fmt.Errorf("negative value not allowed")
	}
	if err != nil {
		cfgError(fmt.Sprintf("log.maxAge is not valid: %v", err))
	}
	t.maxAgeDuration = duration

	archive := strings.ToLower(t.cfg.Archive)
	if !(archive == "" || archive == "zip") {
		cfgError("log.archive could be either empty or has \"zip\" value")
	}
	t.archive = archive == "zip"
}

func (t *Task) writeConfig(cfg *api.Config) {
	cfg.Log = t.cfg
}

// Open func
func (t *Task) Open(ctx *tasks.ServiceTaskContext) error {
	logWrite = func(fields ...string) {
		if t.cfg == nil {
			return
		}
		ctx.Wg.Add(1)
		defer ctx.Wg.Done()

		fields = append([]string{time.Now().Local().Format("2006-01-02T15:04:05.999")}, fields...)

		var record bytes.Buffer
		csvw := csv.NewWriter(&record)
		csvw.Write(fields)
		csvw.Flush()

		logFile := filepath.Join(t.cfg.Dir, t.fileName)

		t.rotate(logFile, &ctx.Wg, ctx.Log)

		var f *os.File
		var err error
		for i := 0; i < 6; i++ {
			f, err = os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, t.fileMode)
			if err == nil {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if err == nil {
			defer f.Close()
			_, err = f.Write(record.Bytes())
		}
		if err != nil {
			ctx.Log.Printf("unable to log the record:%v%v%vreason: %v", utils.NewLine, record.String(), utils.NewLine, err)
		}
	}
	return nil
}

// Close func
func (t *Task) Close(ctx *tasks.ServiceTaskContext) error {
	return nil
}

// Stop func
func (t *Task) Stop(ctx *tasks.ServiceTaskContext) {
}
