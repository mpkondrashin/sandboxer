/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

task.go

Inspection task
*/
package task

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sandboxer/pkg/globals"
	"sandboxer/pkg/logging"
	"sandboxer/pkg/sandbox"
)

type ID int64

type Task struct {
	Number        ID
	Type          TaskType
	SubmitTime    time.Time
	Path          string
	Channel       Channel
	RiskLevel     sandbox.RiskLevel
	Active        bool
	Message       string
	SandboxID     string
	MD5           string
	SHA1          string
	SHA256        string
	Report        string
	Investigation string
}

func NewTask(id ID, taskType TaskType, path string) *Task {
	return &Task{
		Number:     id,
		Type:       taskType,
		SubmitTime: time.Now(),
		Path:       path,
		Channel:    ChPrefilter,
		RiskLevel:  sandbox.RiskLevelUnknown,
		Active:     false,
		Message:    "",
		SandboxID:  "",
	}
}

func LoadTask(filePath string) (*Task, error) {
	t := &Task{}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Task) SetChannel(newChannel Channel) {
	logging.Debugf("SetChannel(%v)", newChannel)
	t.Channel = newChannel
	logging.LogError(t.SaveIfNeeded())
}

func (t *Task) GetChannel() string {
	if t.Channel == ChDone { //} && t.RiskLevel  != sandbox.RiskLevelUnknown {
		return t.RiskLevel.String()
	}
	return t.Channel.String()
}

func (t *Task) VOneID() string {
	return t.SandboxID
}

func (t *Task) SetSandboxID(sandboxID string) {
	t.SandboxID = sandboxID
}

func (t *Task) String() string {
	return fmt.Sprintf("Task %d; submitted on: %v; channel: %v; id: %s; message: %s, path: %s", t.Number, t.SubmitTime, t.Channel, t.SandboxID, t.Message, t.Path)
}
func (t *Task) SetRiskLevel(riskLevel sandbox.RiskLevel) {
	//t.State = StateDone
	t.RiskLevel = riskLevel
}

func (t *Task) Title() string {
	if t.Type == URLTask {
		return t.Path
	} else {
		return filepath.Base(t.Path)
	}
}

func (t *Task) SetError(err error) {
	t.Channel = ChDone
	t.RiskLevel = sandbox.RiskLevelError
	t.Message = err.Error()
}

func (t *Task) SetMessage(message string) {
	t.Message = message
}

func (t *Task) SetReport(report string) {
	t.Report = report
}

func (t *Task) SetInvestigation(investigation string) {
	t.Investigation = investigation
}

func (t *Task) Activate() {
	t.Active = true
}

func (t *Task) Deactivate() {
	t.Active = false
}

/*
	func (t *Task) SetDigest(MD5, SHA1, SHA256 string) {
		if MD5 != "" {
			t.MD5 = MD5
		}
		if SHA1 != "" {
			t.SHA1 = SHA1
		}
		if SHA256 != "" {
			t.SHA256 = SHA256
		}
	}
*/
func (t *Task) SaveIfNeeded() error {
	return t.Save()
}

const taskFileName = "task.json"

func (t *Task) Save() error {
	tasksFolder, err := t.Folder()
	if err != nil {
		return err
	}
	taskFilePath := filepath.Join(tasksFolder, taskFileName)
	return t.SaveToFile(taskFilePath)
}

func (t *Task) SaveToFile(filePath string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

var ErrMissingHash = errors.New("missing hash")

func (t *Task) Folder() (string, error) {
	if t.SHA256 == "" {
		return "", ErrMissingHash
	}
	tasksFolder, err := globals.TasksFolder()
	if err != nil {
		return "", err
	}
	folder := filepath.Join(tasksFolder, t.SHA256)
	if err := os.MkdirAll(folder, 0755); err != nil {
		return "", err
	}
	return folder, nil
}

func (t *Task) ReportPath() (string, error) {
	folder, err := t.Folder()
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("report_%s.pdf", t.SHA256)
	return filepath.Join(folder, fileName), nil
}

func (t *Task) InvestigationPath() (string, error) {
	folder, err := t.Folder()
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("investigation_%s.zip", t.SHA256)
	return filepath.Join(folder, fileName), nil
}

func (t *Task) CalculateHash() error {
	var source io.Reader
	if t.Type == FileTask {
		src, err := os.Open(t.Path)
		if err != nil {
			return err
		}
		defer src.Close()
		source = src
	} else {
		source = strings.NewReader(t.Path)
	}
	SHA256 := sha256.New()
	srcWithSHA256 := io.TeeReader(source, SHA256)
	SHA1 := sha1.New()
	srcWithSHA256andSHA1 := io.TeeReader(srcWithSHA256, SHA1)
	MD5 := md5.New()
	if _, err := io.Copy(MD5, srcWithSHA256andSHA1); err != nil {
		return err
	}
	t.SHA256 = hex.EncodeToString(SHA256.Sum(nil))
	t.SHA1 = hex.EncodeToString(SHA1.Sum(nil))
	t.MD5 = hex.EncodeToString(MD5.Sum(nil))
	return nil
}

func (t *Task) Delete() error {
	folder, err := t.Folder()
	if err != nil {
		return err

	}
	return os.RemoveAll(folder)
}
