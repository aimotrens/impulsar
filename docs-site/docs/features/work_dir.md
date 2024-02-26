# Change your working directory for a job

You can change the working directory for a job with the `workDir` field.

```yaml
job-backend:
  workDir: backend
  script:
    - echo "Run this in the backend directory"

job-frontend:
  workDir: frontend
  script:
    - echo "Run this in the frontend directory"
```
!!! note
    This feature is available in all shell types except `ssh`.
