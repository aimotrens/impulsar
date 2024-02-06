package sshexecutor

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/aimotrens/impulsar/engine"
	"github.com/aimotrens/impulsar/model"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type SshExecutor struct {
	*engine.Engine
}

func init() {
	engine.RegisterExecutor(model.SHELL_TYPE_SSH, func(e *engine.Engine) engine.Shell {
		return &SshExecutor{Engine: e}
	})
}

func (e *SshExecutor) Execute(j *model.Job, script string) {
	user, server, port, err := splitServerString(j)
	if checkError(j, err) {
		return
	}

	agentConn, err := newAgentConnection()
	if checkError(j, err) {
		return
	}

	agentClient := agent.NewClient(agentConn)
	defer agentConn.Close()

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server, port), sshConfig)
	if checkError(j, err) {
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if checkError(j, err) {
		return
	}
	defer session.Close()

	session.Stdout = &engine.JobOutputUnifier{Job: j, ScriptLine: &script, Writer: os.Stdout}
	session.Stderr = &engine.JobOutputUnifier{Job: j, ScriptLine: &script, Writer: os.Stderr}

	scriptExpanded := os.Expand(script, e.LookupVarFunc(j))

	err = session.Run(scriptExpanded)

	switch err.(type) {
	case *ssh.ExitMissingError:
		return
	default:
		checkError(j, err)
	}
}

func checkError(j *model.Job, err error) bool {
	if err != nil {
		if j.AllowFail {
			return true
		}

		fmt.Println(err)
		os.Exit(1)
	}

	return false
}

func splitServerString(j *model.Job) (user string, server string, port uint16, err error) {
	re := regexp.MustCompile(`^(?P<user>[^@]+)@(?P<server>[^:]+)(?::(?P<port>\d{1,5}))?$`)
	matches := re.FindStringSubmatch(j.Shell.Server)

	if len(matches) == 0 {
		err = fmt.Errorf("[%s] invalid server string: %s", j.Name, j.Shell.Server)
		return
	}

	user = matches[re.SubexpIndex("user")]
	server = matches[re.SubexpIndex("server")]
	port = 22

	if matches[re.SubexpIndex("port")] != "" {
		intPort, _ := strconv.Atoi(matches[re.SubexpIndex("port")])
		port = uint16(intPort)
	}

	return
}
