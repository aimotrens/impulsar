# The shell

The `shell` field allows you to specify the shell to run the job. The default shell is platform dependent. On Windows, the default shell is `powershell` and on Linux, the default shell is `bash`.

## Supported shells

- `bash`
- `powershell`
- `pwsh`
- `docker`
- `ssh`

## Usage

You can use different shells on different platform. For example, you can also use `bash` on Windows (via WSL) or `pwsh` on Linux, when you have installed the necessary shell on your system.

## Examples

### Bash

You can run a job script with the `bash` shell on linux and windows if you have installed the WSL.

```yaml
job:
  shell:
    type: bash
  script:
    - echo "Hello, world!"
```

### Powershell / pwsh

You can run this on Windows. On Linux you can only run the job2 with `pwsh` installed.

```yaml
job1:
  shell:
    type: powershell
  script:
    - echo "Hello, world!"

job2:
  shell:
    type: pwsh
  script:
    - echo "Hello, world!"
```

### Docker

You can run a job script with the `docker` shell and an specified image on linux and windows if you have installed docker.

```yaml
job:
  shell:
    type: docker
    image: debian
  script:
    - echo "Hello from docker container!"
```

### SSH

You can run a job script with the `ssh` shell on linux and windows if you have installed the ssh client.
This runs the script on a remote server.

!!! note
    impulsar only supports passwordless ssh connections via ssh-agent.


```yaml
job:
  shell:
    type: ssh
    server: user@example.net:22
  script:
    - echo "Hello from remote server!"
```
