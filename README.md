# negima

This Action runs your [Jest](https://github.com/facebook/jest) tests and notify Slack of failed tests.

## Simple Usage

```
- name: Test
  uses: skuwa229/negima@master
  env:
    INCOMING_WEBHOOK_URL: "https://hooks.slack.com/services/XXXXXXXXX/YYYYYYYYYYYYYYYYYYYYYY"
    JEST_FILE_PATH: "result.json"
```
