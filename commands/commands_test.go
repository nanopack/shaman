package commands_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"

	"github.com/nanopack/shaman/api"
	"github.com/nanopack/shaman/commands"
	"github.com/nanopack/shaman/config"
)

func init() {
	shamanTool.AddCommand(commands.AddDomain)
	shamanTool.AddCommand(commands.DelDomain)
	shamanTool.AddCommand(commands.ListDomains)
	shamanTool.AddCommand(commands.GetDomain)
	shamanTool.AddCommand(commands.UpdateDomain)
	shamanTool.AddCommand(commands.ResetDomains)

	config.AddFlags(shamanTool)
}

type (
	execable func() error // cobra.Command.Execute() 'alias'
)

var shamanTool = &cobra.Command{
	Use:   "shaman",
	Short: "shaman - api driven dns server",
	Long:  ``,

	Run: startShaman,
}

func startShaman(ccmd *cobra.Command, args []string) {
	return
}

func TestMain(m *testing.M) {
	// manually configure
	initialize()

	// start api
	go api.Start()
	<-time.After(time.Second)
	rtn := m.Run()

	os.Exit(rtn)
}

func TestAddRecord(t *testing.T) {
	commands.ResetVars()

	args := strings.Split("add -d nanobox.io -A 127.0.0.1", " ")
	shamanTool.SetArgs(args)

	out, err := capture(shamanTool.Execute)
	if err != nil {
		t.Errorf("Failed to execute - %v", err.Error())
	}

	if string(out) != "{\"domain\":\"nanobox.io.\",\"records\":[{\"ttl\":60,\"class\":\"IN\",\"type\":\"A\",\"address\":\"127.0.0.1\"}]}\n" {
		t.Errorf("Unexpected output: %+q", string(out))
	}
}

func TestListRecords(t *testing.T) {
	commands.ResetVars()

	args := strings.Split("list", " ")
	shamanTool.SetArgs(args)

	out, err := capture(shamanTool.Execute)
	if err != nil {
		t.Errorf("Failed to execute - %v", err.Error())
	}

	if string(out) != "[\"nanobox.io\"]\n" {
		t.Errorf("Unexpected output: %+q", string(out))
	}

	args = strings.Split("list -f", " ")
	shamanTool.SetArgs(args)

	out, err = capture(shamanTool.Execute)
	if err != nil {
		t.Errorf("Failed to execute - %v", err.Error())
	}

	if string(out) != "[{\"domain\":\"nanobox.io.\",\"records\":[{\"ttl\":60,\"class\":\"IN\",\"type\":\"A\",\"address\":\"127.0.0.1\"}]}]\n" {
		t.Errorf("Unexpected output: %+q", string(out))
	}
}

func TestResetRecords(t *testing.T) {
	commands.ResetVars()

	args := strings.Split("reset -j [{\"domain\":\"nanopack.io\"}]", " ")
	shamanTool.SetArgs(args)

	out, err := capture(shamanTool.Execute)
	if err != nil {
		t.Errorf("Failed to execute - %v", err.Error())
	}

	if string(out) != "[{\"domain\":\"nanopack.io.\",\"records\":null}]\n" {
		t.Errorf("Unexpected output: %+q", string(out))
	}

	args = strings.Split("list", " ")
	shamanTool.SetArgs(args)

	out, err = capture(shamanTool.Execute)
	if err != nil {
		t.Errorf("Failed to execute - %v", err.Error())
	}

	if string(out) != "[\"nanopack.io\"]\n" {
		t.Errorf("Unexpected output: %+q", string(out))
	}
}

func TestUpdateRecord(t *testing.T) {
	commands.ResetVars()

	args := strings.Split("update -d nanopack.io -A 127.0.0.5", " ")
	shamanTool.SetArgs(args)

	out, err := capture(shamanTool.Execute)
	if err != nil {
		t.Errorf("Failed to execute - %v", err.Error())
	}

	if string(out) != "{\"domain\":\"nanopack.io.\",\"records\":[{\"ttl\":60,\"class\":\"IN\",\"type\":\"A\",\"address\":\"127.0.0.5\"}]}\n" {
		t.Errorf("Unexpected output: %+q", string(out))
	}

	args = strings.Split("list", " ")
	shamanTool.SetArgs(args)

	out, err = capture(shamanTool.Execute)
	if err != nil {
		t.Errorf("Failed to execute - %v", err.Error())
	}

	if string(out) != "[\"nanopack.io\"]\n" {
		t.Errorf("Unexpected output: %+q", string(out))
	}
}

func TestGetRecord(t *testing.T) {
	commands.ResetVars()

	args := strings.Split("get -d nanopack.io", " ")
	shamanTool.SetArgs(args)

	out, err := capture(shamanTool.Execute)
	if err != nil {
		t.Errorf("Failed to execute - %v", err.Error())
	}

	if string(out) != "{\"domain\":\"nanopack.io.\",\"records\":[{\"ttl\":60,\"class\":\"IN\",\"type\":\"A\",\"address\":\"127.0.0.5\"}]}\n" {
		t.Errorf("Unexpected output: %+q", string(out))
	}
}

func TestDeleteRecord(t *testing.T) {
	commands.ResetVars()

	args := strings.Split("delete -d nanopack.io", " ")
	shamanTool.SetArgs(args)

	out, err := capture(shamanTool.Execute)
	if err != nil {
		t.Errorf("Failed to execute - %v", err.Error())
	}

	if string(out) != "{\"msg\":\"success\"}\n" {
		t.Errorf("Unexpected output: %+q", string(out))
	}
}

///////////////////////////////////////////////////
// PRIVS
///////////////////////////////////////////////////

// function to capture output of cli
func capture(fn execable) ([]byte, error) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := fn()
	os.Stdout = oldStdout
	w.Close() // do not defer after os.Pipe()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r)
}

// manually configure and start internals
func initialize() {
	config.Insecure = true
	config.L2Connect = "none://"
	config.ApiListen = "127.0.0.1:1634"
	config.Log = lumber.NewConsoleLogger(lumber.LvlInt("FATAL"))
	config.LogLevel = "FATAL"
}
