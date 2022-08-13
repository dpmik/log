package log

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"testing"
)

const ts = `^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{6} `

var lp = [...]string{"INFO> ", "WARN> ", "ERROR> "}

var tt = []struct {
	name     string
	f        func()
	minLevel Severity
	prefix   string
	want     string
}{
	{"Info normal", func() { Info("Ciao") }, LevelInfo, lp[0], "Ciao"},
	{"Info double string", func() { Info("Ciao", "ciao") }, LevelInfo, lp[0], "Ciaociao"},
	{"Info string number", func() { Info("Ciao", 7) }, LevelInfo, lp[0], "Ciao7"},
	{"Info number string", func() { Info(7, "Ciao") }, LevelInfo, lp[0], "7Ciao"},
	{"Info double number", func() { Info(3, 7) }, LevelInfo, lp[0], "3 7"},
	{"Info level warning", func() { Info("Ciao") }, LevelWarning, "", ""},
	{"Info level error", func() { Info("Ciao") }, LevelError, "", ""},
	{"Infof normal", func() { Infof("Ciao") }, LevelInfo, lp[0], "Ciao"},
	{"Infof format", func() { Infof("fmt: %s %v", "ciao", 7) }, LevelInfo, lp[0], "fmt: ciao 7"},
	{"Infof level warning", func() { Infof("Ciao") }, LevelWarning, "", ""},
	{"Infof level error", func() { Infof("Ciao") }, LevelError, "", ""},
	{"Infoln normal", func() { Infoln("Ciao") }, LevelInfo, lp[0], "Ciao"},
	{"Infoln double string", func() { Infoln("Ciao", "ciao") }, LevelInfo, lp[0], "Ciao ciao"},
	{"Infoln string number", func() { Infoln("Ciao", 7) }, LevelInfo, lp[0], "Ciao 7"},
	{"Infoln number string", func() { Infoln(7, "Ciao") }, LevelInfo, lp[0], "7 Ciao"},
	{"Infoln double number", func() { Infoln(3, 7) }, LevelInfo, lp[0], "3 7"},
	{"Infoln level warning", func() { Infoln("Ciao") }, LevelWarning, "", ""},
	{"Infoln level error", func() { Infoln("Ciao") }, LevelError, "", ""},
	{"Warning normal", func() { Warning("Ciao") }, LevelInfo, lp[1], "Ciao"},
	{"Warning double string", func() { Warning("Ciao", "ciao") }, LevelInfo, lp[1], "Ciaociao"},
	{"Warning string number", func() { Warning("Ciao", 7) }, LevelInfo, lp[1], "Ciao7"},
	{"Warning number string", func() { Warning(7, "Ciao") }, LevelInfo, lp[1], "7Ciao"},
	{"Warning double number", func() { Warning(3, 7) }, LevelInfo, lp[1], "3 7"},
	{"Warning level warning", func() { Warning("Ciao") }, LevelWarning, "", ""},
	{"Warning level error", func() { Warning("Ciao") }, LevelError, "", ""},
	{"Warningf normal", func() { Warningf("Ciao") }, LevelInfo, lp[1], "Ciao"},
	{"Warningf format", func() { Warningf("fmt: %s %v", "ciao", 7) }, LevelInfo, lp[1], "fmt: ciao 7"},
	{"Warningf level warning", func() { Warningf("Ciao") }, LevelWarning, lp[1], "Ciao"},
	{"Warningf level error", func() { Warningf("Ciao") }, LevelError, "", ""},
	{"Warningln normal", func() { Warningln("Ciao") }, LevelInfo, lp[1], "Ciao"},
	{"Warningln double string", func() { Warningln("Ciao", "ciao") }, LevelInfo, lp[1], "Ciao ciao"},
	{"Warningln string number", func() { Warningln("Ciao", 7) }, LevelInfo, lp[1], "Ciao 7"},
	{"Warningln number string", func() { Warningln(7, "Ciao") }, LevelInfo, lp[1], "7 Ciao"},
	{"Warningln double number", func() { Warningln(3, 7) }, LevelInfo, lp[1], "3 7"},
	{"Warningln level warning", func() { Warningln("Ciao") }, LevelWarning, lp[1], "Ciao"},
	{"Warningln level error", func() { Warningln("Ciao") }, LevelError, "", ""},
	{"Error normal", func() { Error("Ciao") }, LevelInfo, lp[2], "Ciao"},
	{"Error double string", func() { Error("Ciao", "ciao") }, LevelInfo, lp[2], "Ciaociao"},
	{"Error string number", func() { Error("Ciao", 7) }, LevelInfo, lp[2], "Ciao7"},
	{"Error number string", func() { Error(7, "Ciao") }, LevelInfo, lp[2], "7Ciao"},
	{"Error double number", func() { Error(3, 7) }, LevelInfo, lp[2], "3 7"},
	{"Error level warning", func() { Error("Ciao") }, LevelWarning, lp[2], "Ciao"},
	{"Error level error", func() { Error("Ciao") }, LevelError, lp[2], "Ciao"},
	{"Errorf normal", func() { Errorf("Ciao") }, LevelInfo, lp[2], "Ciao"},
	{"Errorf format", func() { Errorf("fmt: %s %v", "ciao", 7) }, LevelInfo, lp[2], "fmt: ciao 7"},
	{"Errorf level warning", func() { Errorf("Ciao") }, LevelWarning, lp[2], "Ciao"},
	{"Errorf level error", func() { Errorf("Ciao") }, LevelError, lp[2], "Ciao"},
	{"Errorln normal", func() { Errorln("Ciao") }, LevelInfo, lp[2], "Ciao"},
	{"Errorln double string", func() { Errorln("Ciao", "ciao") }, LevelInfo, lp[2], "Ciao ciao"},
	{"Errorln string number", func() { Errorln("Ciao", 7) }, LevelInfo, lp[2], "Ciao 7"},
	{"Errorln number string", func() { Errorln(7, "Ciao") }, LevelInfo, lp[2], "7 Ciao"},
	{"Errorln double number", func() { Errorln(3, 7) }, LevelInfo, lp[2], "3 7"},
	{"Errorln level warning", func() { Errorln("Ciao") }, LevelWarning, lp[2], "Ciao"},
	{"Errorln level error", func() { Errorln("Ciao") }, LevelError, lp[2], "Ciao"},
	{"Verbose", func() { Verbose(true); Info("Ciao") }, LevelInfo, "log_test.go:[0-9]+: " + lp[0], "Ciao"},
	{"Verbose enabled and disabled", func() { Verbose(true); Verbose(false); Info("Ciao") }, LevelInfo, lp[0], "Ciao"},
}

func TestOutput(t *testing.T) {
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := new(bytes.Buffer)
			SetWriter(w)
			SetLevel(tc.minLevel)
			tc.f()

			line := w.String()
			if len(line) > 0 {
				line = line[0 : len(line)-1]
			}
			var pattern string
			if tc.want != "" {
				pattern = ts + tc.prefix + tc.want + "$"
			}
			matched, err := regexp.MatchString(pattern, line)
			if err != nil {
				t.Fatalf("unable to compile regex %q: %v", pattern, err)
			}
			if !matched {
				t.Fatalf("mismatch! Pattern %q, got %q", pattern, line)
			}
		})
	}
}

func TestFatals(t *testing.T) {
	tt := []struct {
		name     string
		f        func()
		minLevel Severity
		want     string
	}{
		{"Fatal normal", func() { Fatal("Ciao") }, LevelInfo, "Ciao"},
		{"Fatal double string", func() { Fatal("Ciao", "ciao") }, LevelInfo, "Ciaociao"},
		{"Fatal string number", func() { Fatal("Ciao", 7) }, LevelInfo, "Ciao7"},
		{"Fatal number string", func() { Fatal(7, "Ciao") }, LevelInfo, "7Ciao"},
		{"Fatal double number", func() { Fatal(3, 7) }, LevelInfo, "3 7"},
		{"Fatal level warning", func() { Fatal("Ciao") }, LevelWarning, "Ciao"},
		{"Fatal level error", func() { Fatal("Ciao") }, LevelError, "Ciao"},
		{"Fatalf normal", func() { Fatalf("Ciao") }, LevelInfo, "Ciao"},
		{"Fatalf format", func() { Fatalf("fmt: %s %v", "ciao", 7) }, LevelInfo, "fmt: ciao 7"},
		{"Fatalf level warning", func() { Fatalf("Ciao") }, LevelWarning, "Ciao"},
		{"Fatalf level error", func() { Fatalf("Ciao") }, LevelError, "Ciao"},
		{"Fatalln normal", func() { Fatalln("Ciao") }, LevelInfo, "Ciao"},
		{"Fatalln double string", func() { Fatalln("Ciao", "ciao") }, LevelInfo, "Ciao ciao"},
		{"Fatalln string number", func() { Fatalln("Ciao", 7) }, LevelInfo, "Ciao 7"},
		{"Fatalln number string", func() { Fatalln(7, "Ciao") }, LevelInfo, "7 Ciao"},
		{"Fatalln double number", func() { Fatalln(3, 7) }, LevelInfo, "3 7"},
		{"Fatalln level warning", func() { Fatalln("Ciao") }, LevelWarning, "Ciao"},
		{"Fatalln level error", func() { Fatalln("Ciao") }, LevelError, "Ciao"},
	}

	idx, err := strconv.Atoi(os.Getenv("FATAL_IDX"))
	if err == nil {
		SetLevel(LevelError)
		tt[idx].f()
		return // just in case...
	}

	for i, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			cmd := exec.Command(os.Args[0], "-test.run=TestFatals")
			cmd.Env = append(os.Environ(), "FATAL_IDX="+strconv.Itoa(i))
			cmd.Stdout = out
			err = cmd.Run()
			var e *exec.ExitError
			if errors.As(err, &e) {
				if ret := e.ExitCode(); ret != 1 {
					t.Fatalf("wrong exit code: want 1, got %d", ret)
				}

				line := out.String()
				if len(line) > 0 {
					line = line[0 : len(line)-1]
				}
				var pattern string
				if tc.want != "" {
					pattern = ts + lp[2] + tc.want + "$"
				}
				matched, err := regexp.MatchString(pattern, line)
				if err != nil {
					t.Fatalf("unable to compile regex %q: %v", pattern, err)
				}
				if !matched {
					t.Fatalf("mismatch! Pattern %q, got %q", pattern, line)
				}
				return
			}
			t.Fatalf("unexpected err value %v, want exit status 1", err)
		})
	}
}

func TestLevel(t *testing.T) {
	l := Level()
	if l != LevelInfo {
		t.Errorf("default log level should be LogInfo (%v), got %v", LevelInfo, l)
	}
	SetLevel(LevelError)
	l = Level()
	if l != LevelError {
		t.Errorf("log level should be LogError (%v), got %v", LevelError, l)
	}
}

func TestWriter(t *testing.T) {
	want := new(bytes.Buffer)
	SetWriter(want)
	got := Writer()
	if want != got {
		t.Error("mismatch on io.Writer parameter after SET/GET cycle")
	}
}

func BenchmarkInfo(b *testing.B) {
	const msg = "Ciao"
	var buf bytes.Buffer
	l := New(LevelInfo)
	l.SetWriter(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Info(msg)
	}
}
