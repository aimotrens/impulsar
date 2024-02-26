# How to use impulsar

Impulsar is a command line tool that allows you to run jobs in a simple and easy way. It is designed to be easy to use.

## Basic usage

The most basic usage is to run a single job. To do this, you can use the `impulsar` command followed by the name of the job you want to run.

```bash
impulsar job1
```

## Running multiple jobs

You can run multiple jobs at once by specifying the names of the jobs you want to run.

```bash
impulsar job1 job2
```

## Running jobs with additional arguments

You can also run jobs with additional arguments. To do this, you can use the `-e` flag followed by the arguments you want to pass to the job(s). An argument is passed to the job as an environment variable.

```bash
impulsar -e "ARG1=value1" -e "ARG2=value2" job1
```

## Shell completion

Impulsar supports shell completion for bash, zsh and powershell. To get the shell completion script, you can let impulsar generate it for you.

```bash
impulsar --gen-bash-completion
impulsar --gen-zsh-completion
impulsar --gen-powershell-completion
```
