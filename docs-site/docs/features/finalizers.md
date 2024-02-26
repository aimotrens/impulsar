# Finalizers

Finalizers are a way to run a script or job after the main job has finished, regardless of the result. This is useful for cleaning up temporary files, sending notifications, or other tasks that should always be executed.

```yaml
test:
  script:
    - echo "Running tests"

build:
    jobs:pre:
      - test
    script:
      - echo "Building the app"
    jobs:finalize:
      - cleanup

cleanup:
    script:
      - echo "Cleaning temporary build files"
```

This conecpt also exists for the script block. You can define a `script:finalize` block to run a script after the main script has finished, regardless of the result.

```yaml
job:
  script:
    - echo "This is a job"
  script:finalize:
    - echo "This is a finalizer"
```
