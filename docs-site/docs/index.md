# Simplify your dev jobs!

impulsar is a tool, to run named tasks in your shell. It is a easy to use tool that can help you to automate your daily tasks.


## Features

- **Simple**: impulsar has a simple but powerful YAML based configuration file.
- **Easy to use**: Start one or more tasks (called jobs) with a single command.
- **Job dependencies**: You can define dependencies between jobs.
- **Foreach**: Run jobs with different sets of EnvVars in a foreach like loop.


## Quickstart

### Download
Download the latest release, that meets your platform, from the [releases page](https://github.com/aimotrens/impulsar/releases/latest) and extract it to a directory in your PATH.


### Create a configuration file

Create a file named `impulsar.yml` with the following content:
```yaml
hello:
  script:
    - echo "Hello from impulsar!"
```

### Run impulsar

Run impulsar with the following command:
```bash
impulsar hello
```
`hello` is the name of the job defined in the configuration file.