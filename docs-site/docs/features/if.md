# The simple If

The `if` feature allows you to run a job conditionally. You can define multiple conditions to run the job. The conditions are evaluated in the order they are defined. If a condition is met, the job is executed.

```yaml
job:
  if: 
   - "env.os === 'windows'"
  script:
    - echo "This job runs only on windows"
```

!!! note
    The conditions are javascript expressions. You can use the `env` object to access environment variables.
