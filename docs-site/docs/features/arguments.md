# Necessary environment variables: Arguments

You can define necessary arguments for a job using the `arguments` field. This field is a list of key-value pairs where the key is the name of the argument and the value is the description of the argument.

If the argument is provided as environment variable, the value of the environment variable is used. If the argument is not provided, impulsar will ask you to provide the value before the jobs are started.

```yaml
job:
  arguments:
    ARG1: Description of ARG1
    ARG2: Description of ARG2
  script:
    - echo "This is a job with arguments"
    - echo "The value of ARG1 is $ARG1"
    - echo "The value of ARG2 is $ARG2"
```

## Default values for arguments

You can also define default values for arguments. To do this, extend the `arguments` field with a `description` and `default` field. 

```yaml
job:
  arguments:
    ARG1:
      description: Description of ARG1
      default: value1
    ARG2: 
      description: Description of ARG2
      default: value2
  script:
    - echo "This is a job with arguments"
    - echo "The value of ARG1 is $ARG1"
    - echo "The value of ARG2 is $ARG2"
```

## Mixing short and long arguments definitions

You can mix short and long arguments definitions. 

```yaml
job:
  arguments:
    ARG1:
      description: Description of ARG1
      default: value1
    ARG2: Description of ARG2
  script:
    - echo "This is a job with arguments"
    - echo "The value of ARG1 is $ARG1"
    - echo "The value of ARG2 is $ARG2"
```
