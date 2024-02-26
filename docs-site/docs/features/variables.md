# Job environment variables

You can set environment variables for a job using the `variables` field. This field is a list of key-value pairs where the key is the name of the environment variable and the value is the value of the environment variable.

```yaml
job:
  variables:
    - KEY1: value1
    - KEY2: value2
  script:
    - echo "This is a job with environment variables"
    - echo "The value of KEY1 is $KEY1"
    - echo "The value of KEY2 is $KEY2"
```

!!! note
    **Variable expansion**

    The variables are expanded before they are passed to the shell. This means that you can use the same "variable-style" in each shell (bash, powershell, ...).

    Normally, an environment variable in powershell would be accessed via $env:KEY1.
    You can use the `$KEY` variant in impulsar with powershell as well.
