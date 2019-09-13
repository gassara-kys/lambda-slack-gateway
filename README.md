# lambda-slack-geteway

gateway function form slack requsts.
```
slack(slash_command) -> API-Gateway -> Lambda(this)
```

## Gateway
- authorization
  - valid time and signature with token.
  - see [slack documet](https://api.slack.com/docs/verifying-requests-from-slack)
- proxy request
  - get parameter
    - get data from DynamoDB table.
  - delete parameter
    - delete table
  - other
    - usage parameters in response.

## Pre-Require
- Create DynamoDB Table.(ex. `sns_alert` table)

- Set the required IAMRole
  - Write to the DynamoDB Table resource

## Run Local
```bash
# make sure that the required IAM is set in advance
$ make
```

## Test DATA data
- generate signature with test-key
- sample: https://play.golang.org/p/QGmUt31yfB-

## Lambda zip file
```bash
$ make zip
```
