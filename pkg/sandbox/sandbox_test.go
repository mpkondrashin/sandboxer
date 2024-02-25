package sandbox

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	"sandboxer/pkg/config"

	"github.com/mpkondrashin/ddan"
	"github.com/mpkondrashin/vone"
	"gopkg.in/yaml.v2"
)

func GetVOne(t *testing.T) *vone.VOne {
	configFile := "visionone.yaml"
	f, err := os.Open(configFile)
	if err != nil {
		t.Fatalf("error opening %s config: %v", configFile, err)
	}
	conf := new(config.VisionOne)
	if err := yaml.NewDecoder(f).Decode(conf); err != nil {
		t.Fatalf("error parsing %s config: %v", configFile, err)
	}
	return vone.NewVOne(conf.Domain, conf.Token)
}

func GetAnalyzer(t *testing.T) *ddan.Client {
	return nil
}

func TestAnalyzer(t *testing.T) {
	analyzer := GetAnalyzer(t)
	analyzerSandbox := NewDDAnSandbox(analyzer)
	SandboxRunTests(analyzerSandbox, t)
	err := analyzer.Unregister(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

func TestVOne(t *testing.T) {
	vOne := GetVOne(t)
	vOneSandbox := NewVOneSandbox(vOne)
	SandboxRunTests(vOneSandbox, t)
}

func SandboxRunTests(sandbox Sandbox, t *testing.T) {
	folder := "testing"
	if err := os.MkdirAll(folder, 0755); err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name         string
		expectedRisk RiskLevel
	}{
		{"novirus", RiskLevelNoRisk},
		{"spyware", RiskLevelHigh},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			Compile("./gmw/"+tCase.name, folder+"/"+tCase.name+".exe", t)
			FileTest(folder+"/nov.exe", RiskLevelNoRisk, sandbox, t)
		})
	}
}

func FileTest(filePath string, expectedRisk RiskLevel, sandbox Sandbox, t *testing.T) {
	id, err := sandbox.SubmitFile(filePath)
	t.Logf("%p.Submit(%s): %s, %v", sandbox, filePath, id, err)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 30; i++ {
		time.Sleep(5 * time.Second)
		riskeLevel, virusName, err := sandbox.GetResult(id)
		t.Logf("GetResult[%d]: %v, %s, %v", i, riskeLevel, virusName, err)
		if err != nil {
			t.Fatal(err)
		}
		if riskeLevel == RiskLevelNotReady {
			continue

		}
		if riskeLevel == expectedRisk {
			return
		}
		t.Errorf("Wrong sandbox response: %v (%s)", riskeLevel, virusName)
		return
	}
	t.Errorf("%s: timeout", filePath)
}

func Compile(sourceFolder, targetFolder string, t *testing.T) {
	_, err := os.Stat(targetFolder)
	if !errors.Is(err, os.ErrNotExist) {
		t.Logf("Skip recompiling %s", sourceFolder)
		return
	}
	c := []string{"build", "-o", targetFolder, sourceFolder}
	cmd := exec.Command("go", c...)
	cmd.Env = append(os.Environ(), "GOOS=windows", "GOARCH=amd64")
	if err := cmd.Run(); err != nil {
		t.Fatalf("error running go executable: %v", err)
	}
}
