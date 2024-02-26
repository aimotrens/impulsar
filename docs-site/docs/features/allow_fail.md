# Allow to fail

Sometimes you want to allow a job to fail without stopping the entire execution.


```yaml
# impulsar.yml

job1:
  allowFail: true
  script:
    - echo "A job that can fail"

job2:
  script:
    - echo "A second job"
```

If you run this two jobs with `impulsar job1 job2`, the second job will run even if the first one fails.
