/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

sandbox_test.go

Test basic sandbox functions
*/
package sandbox

import (
	"context"
	"errors"
	"net/url"
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

func GetAnalyzer(t *testing.T) ddan.ClientInterface {
	configFile := "analyzer.yaml"
	f, err := os.Open(configFile)
	if err != nil {
		t.Fatalf("error opening %s config: %v", configFile, err)
	}
	conf := new(config.DDAn)
	if err := yaml.NewDecoder(f).Decode(conf); err != nil {
		t.Fatalf("error parsing %s config: %v", configFile, err)
	}
	t.Logf("Config: %v", conf)
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatal(err)
	}
	url, err := url.Parse(conf.URL)
	if err != nil {
		t.Fatal(err)
	}
	d := ddan.NewClient(conf.ProductName, hostname).
		SetAnalyzer(url, conf.APIKey, conf.IgnoreTLSErrors).
		SetUUID(conf.ClientUUID).
		SetSource(conf.SourceID, conf.SourceName)
	if conf.ProtocolVersion != "" {
		d.SetProtocolVersion(conf.ProtocolVersion)
	}
	t.Logf("ddan: %v", d)
	return d
}
func TestAnalyzer(t *testing.T) {
	analyzer := GetAnalyzer(t)
	analyzerSandbox := NewDDAnSandbox(analyzer)
	t.Run("file", func(t *testing.T) {
		SandboxRunTests(analyzerSandbox, t)
	})
	t.Run("urls", func(t *testing.T) {
		SandboxRunTestsURLs(analyzerSandbox, t)
	})
	err := analyzer.Unregister(context.TODO())
	if err != nil {
		t.Error(err)
	}
}

func TestVOne(t *testing.T) {
	vOne := GetVOne(t)
	vOneSandbox := NewVOneSandbox(vOne)
	t.Run("file", func(t *testing.T) {
		SandboxRunTests(vOneSandbox, t)
	})
	t.Run("urls", func(t *testing.T) {
		SandboxRunTestsURLs(vOneSandbox, t)
	})

}

func SandboxRunTestsURLs(sandbox Sandbox, t *testing.T) {
	testCases := []struct {
		url          string
		expectedRisk RiskLevel
	}{
		{"http://www.ru", RiskLevelNoRisk},
		{"http://wrs21.winshipway.com", RiskLevelHigh},
	}
	for _, tCase := range testCases {
		t.Run(tCase.url, func(t *testing.T) {
			SubmitTest(false, tCase.url, tCase.expectedRisk, sandbox, t)
		})
	}
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
		//{"novirus", RiskLevelNoRisk},
		{"spyware", RiskLevelHigh},
		//{"downloader", RiskLevelHigh},
		{"dropper", RiskLevelLow},
		//{"encryptor", RiskLevelHigh},
		{"rmpy", RiskLevelNoRisk},
		//{"timeout", RiskLevelNoRisk},
	}
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			fileName := folder + "/" + tCase.name + ".exe"
			Compile("./gmw/"+tCase.name, fileName, t)
			SubmitTest(true, fileName, tCase.expectedRisk, sandbox, t)
		})
	}
}

func SubmitTest(file bool, filePath string, expectedRisk RiskLevel, sandbox Sandbox, t *testing.T) {
	//t.Parallel()
	var id string
	var err error
	if file {
		id, err = sandbox.SubmitFile(filePath)
	} else {
		id, err = sandbox.SubmitURL(filePath)
	}
	//t.Logf("%p.Submit(%s): %s, %v", sandbox, filePath, id, err)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 240; i++ {
		time.Sleep(5 * time.Second)
		riskeLevel, virusName, err := sandbox.GetResult(id)
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
