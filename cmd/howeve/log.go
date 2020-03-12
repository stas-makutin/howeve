package main

type logTask struct {
}

func newLogTask() *logTask {
	return &logTask{}
}

func (t *logTask) open(ctx *serviceTaskContext) error {
	return nil
}

func (t *logTask) close(ctx *serviceTaskContext) error {
	return nil
}

func (t *logTask) stop(ctx *serviceTaskContext) {
}

func logrus(fields ...string) {

}
