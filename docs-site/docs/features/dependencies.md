# Dependencies between jobs

Sometimes you want to define dependencies between jobs. For example, you want to build your app, but before this, you want to run some tests. You can define this with the `pre` and `post` job extensions.

```yaml
test:
  script:
    - echo "Running tests"

build:
    jobs:pre:
      - test
    script:
      - echo "Building the app"
    jobs:post:
      - cleanup

cleanup:
    script:
      - echo "Cleaning temporary build files"
```

Execute this with `impulsar build` and the `test` job will run before the `build` job. After the `build` job, the `cleanup` job will run.

