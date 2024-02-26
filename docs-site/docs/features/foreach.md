# Foreach, to short your yaml

You can use the `foreach` field to run a job multiple times with different variable sets. This field is a list of key-value pairs where the key is the name of the environment variable and the value is the value of the environment variable.

For each element in the list, the script block is executed with the given variables.

```yaml
job:
  foreach:
    - VAR1: value1
      VAR2: value2
    - VAR1: second value1
      VAR2: second value2
    script:
    - echo "This is a job with foreach"
    - echo "The value of VAR1 is $VAR1"
    - echo "The value of VAR2 is $VAR2"
```

## Pre and Post Scripts

The main-`script` block is executed for each element in the `foreach` list. To run a script only once before or after the main-`script` block, you can use the `pre` and `post` extensions. These fields are a list of commands to run before or after the main-`script` block.

```yaml
job:
  foreach:
    - VAR1: value1
      VAR2: value2
    - VAR1: second value1
      VAR2: second value2
  script:pre:
    - echo "This is a pre script, it is executed before the main script"
  script:
    - echo "This is a job with foreach"
    - echo "The value of VAR1 is $VAR1"
    - echo "The value of VAR2 is $VAR2"
  script:post:
    - echo "This is a post script, it is executed after the main script"
```
