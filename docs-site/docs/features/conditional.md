# Conditional

The `conditional` keyword allows you to run a modify the job conditionally. This is useful if you want to run different scripts on different platforms.

You can overwrite all of the job properties, except the `name` and `conditional` fields.

```yaml
job:
  conditional:
    - if: ["env.os === 'windows'"]
      overwrite:
        script:
          - echo "Do windows stuff"
    - if: ["env.os === 'linux'"]
      overwrite:
        script:
        - echo "Do linux stuff"
  script:
    - echo "Error message, if no condition is met"
```