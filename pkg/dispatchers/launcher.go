package dispatchers

import (
	"sandboxer/pkg/config"
	"sandboxer/pkg/fifo"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/task"
	"sync"
)

const (
	PrefilterDispatchers = 1
	UploadDispatchers    = 5
	WaitDispatchers      = 5
	ResultDispatchers    = 5
)

type Launcher struct {
	conf     *config.Configuration
	channels *Channels
	list     *task.TaskList
}

func NewLauncher(conf *config.Configuration, channels *Channels, list *task.TaskList) *Launcher {
	return &Launcher{
		conf:     conf,
		channels: channels,
		list:     list,
	}
}

func (l *Launcher) Run() {
	base := NewBaseDispatcher(l.conf, l.channels, l.list)
	dispatchers := []struct {
		count      int
		dispatcher Dispatcher
	}{
		{ResultDispatchers, NewResultDispatch(base)},
		{WaitDispatchers, NewWaitDispatch(base)},
		{UploadDispatchers, NewUploadDispatch(base)},
		{PrefilterDispatchers, NewPrefilterDispatch(base)},
	}
	var wg sync.WaitGroup
	for _, d := range dispatchers {
		for i := 0; i < d.count; i++ {
			wg.Add(1)
			go l.RunDispatcher(d.dispatcher, &wg)
		}
	}
	submit := NewSubmitDispatch(base)
	wg.Add(1)
	go submit.Run(&wg)
	//wg.Wait()
}

func (l *Launcher) RunDispatcher(disp Dispatcher, wg *sync.WaitGroup) {
	//ctx, cancel := context.WithCancel(context.TODO())
	logging.Debugf("Start %T", disp)
	ch := disp.InboundChannel()
	for id := range l.channels.TaskChannel[ch] {
		_ = l.list.Task(id, func(tsk *task.Task) error {
			//tsk := l.list.Get(id)
			err := disp.ProcessTask(tsk)
			if err != nil {
				tsk.SetError(err)
				logging.Errorf("Task #%d: %v (%T)", id, err, disp)
			}
			return nil
		})
	}
	wg.Done()
}

func (l *Launcher) Stop() error {
	l.channels.Close() // Should we move it to the end?
	fifoWriter, err := fifo.NewWriter()
	if err != nil {
		return err
	}
	defer func() {
		logging.LogError(fifoWriter.Close())
	}()
	if err = fifoWriter.Write(StopPath); err != nil {
		return err
	}
	return nil
}

/*
type Scan struct {
	list   *task.TaskList
	config *config.Configuration
	//check  *goperic.Periculosum
}

func NewScan(config *config.Configuration , list *task.TaskList) *Scan {

	return &Scan{
		list:   list,
		config: config,
		//check:  check,
	}
}

func (s *Scan) vOne() *vone.VOne {
	return vone.NewVOne(s.config.VisionOne.Domain, s.config.VisionOne.Token)
}

func (s *Scan) InspecfFolder(folderPath string) {
	logging.Debugf("InspectFolder(%s)", folderPath)
	err := filepath.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				//go
				s.InspectFile(path)
			}
			return nil
		})
	logging.LogError(err)
}

func (s *Scan) ShouldIgnore(filePath string) bool {
	fileName := filepath.Base(filePath)
	for _, mask := range s.config.Ignore {
		result, err := filepath.Match(strings.ToLower(mask), strings.ToLower(fileName))
		logging.LogError(err)
		if result {
			logging.Debugf("%s: ignore by mask \"%s\"", filePath, mask)
			return true
		}
	}
	return false
}

func (s *Scan) InspectFile(filePath string) {
	logging.Debugf("InspectFile(%s)", filePath)
	info, err := os.Lstat(filePath)
	if err != nil {
		logging.Errorf("%s", err)
		return
	}
	if info.IsDir() {
		s.InspecfFolder(filePath)
		return
	}
	if !info.Mode().IsRegular() {
		logging.Errorf("%s: not regular file", filePath)
	}
	if s.ShouldIgnore(filePath) {
		return
	}
	t := task.New(filePath)
	if err := s.Submit(t); err != nil {
		logging.LogError(err)
		task.SetError(t, err)
		return
	}
}

func (s *Scan) Submit(t task.ID) error {
	if s.config.VisionOne.Token == "" {
		return errors.New("token is not set")
	}
	if s.config.VisionOne.Domain == "" {
		return errors.New("domain is not set")
	}
	task.SetState(t, state.StateUpload)
	//task.SetState(t, state.State(rand.Int()%int(state.StateCount)))
	//return nil
	f, err := s.vOne().SandboxSubmitFile().SetFilePath(task.Path(t))
	if err != nil {
		return err
	}
	response, headers, err := f.Do(context.TODO())
	_ = headers
	if err != nil {
		return err
	}
	task.SetSandboxID(t, response.ID)
	logging.Infof("Accepted: %v", t)
	if err := s.WaitForResult(t); err != nil {
		return err
	}
	//c.LogQuota(id, headers)
	if err := s.GetResult(t); err != nil {
		return err
	}
	return nil
}

func (s *Scan) WaitForResult(t task.ID) error {
	task.SetState(t, state.StateInspected)
	for {
		// Should we set temporary state "checking"?
		status, err := s.vOne().SandboxSubmissionStatus(task.GetSandboxID(t)).Do(context.TODO())
		if err != nil {
			return fmt.Errorf("check status: %w", err)
		}
		logging.Debugf("%s Status: %v", task.GetSandboxID(t), status.Status)
		switch status.Status {
		case vone.StatusSucceeded:
			return nil
		case vone.StatusRunning:
			//if time.Now().After(endTime) {
			//	return ErrTimeout
			//}
			time.Sleep(s.config.VisionOne.Sleep)
		case vone.StatusFailed:
			return fmt.Errorf("%s: %s", status.Error.Code, status.Error.Message)
		default:
			return fmt.Errorf("unknown status: %s", status.Status)
		}
	}
}

func (s *Scan) GetResult(t task.ID) error {
	results, err := s.vOne().SandboxAnalysisResults(task.GetSandboxID(t)).Do(context.TODO())
	if err != nil {
		return err
	}
	detectionName := strings.Join(results.DetectionNames, ", ")
	threatType := strings.Join(results.ThreatTypes, ", ")
	logging.Debugf("Type: %s, TrueFileType: %s, RiskLevel: %s, DetectionNames: %s, threatTypes: %s; for task %v",
		results.Type, results.TrueFileType, results.RiskLevel, detectionName, threatType, t)
	switch results.RiskLevel {
	case vone.RiskLevelHigh:
		task.SetState(t, state.StateHighRisk)
	case vone.RiskLevelMedium:
		task.SetState(t, state.StateMediumRisk)
	case vone.RiskLevelLow:
		task.SetState(t, state.StateLowRisk)
	case vone.RiskLevelNoRisk:
		task.SetState(t, state.StateNoRisk)
	default:
		err := fmt.Errorf("unknown risk level: %d", results.RiskLevel)
		logging.LogError(err)
		task.SetError(t, err)
	}
	return nil
}

func RunService(conf *config.Configuration, list *task.TaskList) (func(), error) {
	inbox := make(StringChannel, 10000)
	//stop := make(chan struct{})
	go SubmitDispatchFunc(inbox)
	//	pericPath, err := config.PericulosumPath()
	//	if err != nil {
	//	}
	//	goperic.NewPericulosum()
	scan := NewScan(conf, list)
	go func() {
		for s := range inbox {
			logging.Debugf("Got %s", s)
			go scan.InspectFile(s)
		}
	}()
	return func() { close(inbox) }, nil
}
*/
