package sshexecutor

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/aimotrens/impulsar/internal/engine"
	"github.com/aimotrens/impulsar/internal/model"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type SshExecutor struct {
	*engine.Engine
}

func init() {
	engine.RegisterExecutor(model.SHELL_TYPE_SSH, func(e *engine.Engine) engine.Executor {
		return &SshExecutor{Engine: e}
	})
}

func (e *SshExecutor) Execute(j *model.Job, script string) error {
	user, server, port, err := splitServerString(j)
	if checkError(j, err) != nil {
		return err
	}

	agentConn, err := newAgentConnection()
	if checkError(j, err) != nil {
		return err
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
	if checkError(j, err) != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if checkError(j, err) != nil {
		return err
	}
	defer session.Close()

	session.Stdout, session.Stderr = engine.GetCmdOutputTarget(j)

	scriptExpanded := e.ExpandVarsWithTemplateEngine(script, j)

	err = session.Run(scriptExpanded)

	switch err.(type) {
	case *ssh.ExitMissingError:
		return nil
	default:
		if checkError(j, err) != nil {
			return err
		}
	}

	return nil
}

func checkError(j *model.Job, err error) error {
	if err != nil {
		if j.AllowFail {
			return nil
		}

		fmt.Println(err)
		return err
	}

	return nil
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
		intPort, e := strconv.ParseInt(matches[re.SubexpIndex("port")], 10, 16)
		if e != nil {
			err = fmt.Errorf("[%s] invalid port: %s", j.Name, matches[re.SubexpIndex("port")])
			return
		}
		port = uint16(intPort)
	}

	return
}
