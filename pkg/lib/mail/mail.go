package _mailUtils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/aaronchen2k/deeptest"
	"github.com/aaronchen2k/deeptest/internal/pkg/config"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/pkg/lib/comm"
	"github.com/aaronchen2k/deeptest/pkg/lib/i118"
	"github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/pkg/lib/string"
	"github.com/kataras/iris/v12/websocket"
	"gopkg.in/gomail.v2"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func Send(to, subject, tmpl string, mp map[string]string) (err error) {
	d := gomail.NewDialer(
		config.CONFIG.Mail.SmtpAddress,
		config.CONFIG.Mail.SmtpPort,
		config.CONFIG.Mail.Account,
		config.CONFIG.Mail.Password)

	s, err := d.Dial()
	if err != nil {
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", consts.SupportEmail)
	m.SetAddressHeader("To", to, mp["name"])
	m.SetHeader("Subject", subject)

	tmplFile := filepath.Join("res", "tmpl", tmpl+".ftl")
	content, _ := deeptest.ReadResData(tmplFile)

	body := os.Expand(string(content), func(k string) string { return mp[k] })

	m.SetBody("text/html", body)

	err = gomail.Send(s, m)
	if err != nil {
		_logUtils.Infof("Failed to send email to %q: %v", mp["toEmail"], err)
	}
	m.Reset()

	return
}

func ExeShell(cmdStr string) (string, error) {
	return ExeShellInDir(cmdStr, "")
}

func ExeShellInDir(cmdStr string, dir string) (ret string, err error) {
	ret, err, _ = ExeShellInDirWithPid(cmdStr, dir)
	return
}

func ExeShellWithPid(cmdStr string) (string, error, int) {
	return ExeShellInDirWithPid(cmdStr, "")
}

func ExeShellInDirWithPid(cmdStr string, dir string) (ret string, err error, pid int) {
	var cmd *exec.Cmd
	if _commUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}
	if dir != "" {
		cmd.Dir = dir
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		_logUtils.Error(_i118Utils.Sprintf("fail_to_exec_command", cmdStr, cmd.Dir, err))
	}

	pid = cmd.Process.Pid
	ret = _stringUtils.TrimAll(out.String())
	return
}

func ExeShellWithOutput(cmdStr string) ([]string, error) {
	return ExeShellWithOutputInDir(cmdStr, "")
}

func ExeShellWithOutputInDir(cmdStr string, dir string) ([]string, error) {
	return ExeShellWithEnvVarsAndOutputInDir(cmdStr, dir, nil)
}

func ExeShellWithEnvVarsAndOutputInDir(cmdStr, dir string, envVars []string) ([]string, error) {
	var cmd *exec.Cmd
	if _commUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	if dir != "" {
		cmd.Dir = dir
	}
	if envVars != nil && len(envVars) > 0 {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envVars...)
	}

	output := make([]string, 0)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return output, err
	}

	cmd.Start()

	if err != nil {
		return output, err
	}

	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		_logUtils.Info(strings.TrimRight(line, "\n"))
		output = append(output, line)
	}

	cmd.Wait()

	return output, nil
}

func ExeShellCallback(ch chan int, cmdStr, dir string,
	fun func(s string, msg websocket.Message), msg websocket.Message) (err error) {

	var cmd *exec.Cmd
	if _commUtils.IsWin() {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	if dir != "" {
		cmd.Dir = dir
	}

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		fmt.Println(err)
		return
	}

	cmd.Start()

	if err != nil {
		return
	}

	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}

		line = strings.Trim(line, "\n")
		fun(line, msg)

		select {
		case <-ch:
			fmt.Println("exiting...")
			ch <- 1
			return
		default:
			fmt.Println("continue...")
		}
	}

	cmd.Wait()
	return
}

func GetProcess(app string) (string, error) {
	var cmd *exec.Cmd

	tmpl := ""
	cmdStr := ""
	if _commUtils.IsWin() {
		tmpl = `tasklist`
		cmdStr = fmt.Sprintf(tmpl)

		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		tmpl = `ps -ef | grep "%s" | grep -v "grep" | awk '{print $2}'`
		cmdStr = fmt.Sprintf(tmpl, app)

		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	output := ""
	if _commUtils.IsWin() {
		arr := strings.Split(out.String(), "\n")
		for _, line := range arr {
			if strings.Index(line, app+".exe") > -1 {
				arr2 := regexp.MustCompile(`\s+`).Split(line, -1)
				output = arr2[1]
				break
			}
		}
	} else {
		output = out.String()
	}

	return output, err
}

func KillProcess(app string) (string, error) {
	var cmd *exec.Cmd

	tmpl := ""
	cmdStr := ""
	if _commUtils.IsWin() {
		// tasklist | findstr ztf.exe
		tmpl = `taskkill.exe /f /im %s.exe`
		cmdStr = fmt.Sprintf(tmpl, app)

		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		tmpl = `ps -ef | grep '%s' | grep -v "grep" | awk '{print $2}' | xargs kill -9`
		cmdStr = fmt.Sprintf(tmpl, app)

		cmd = exec.Command("/bin/bash", "-c", cmdStr)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	output := out.String()

	return output, err
}

func KillProcessById(pid int) {
	cmdStr := fmt.Sprintf("kill -9 %d", pid)
	ExeShell(cmdStr)
}
