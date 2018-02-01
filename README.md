<a href="https://github.com/fsm"><p align="center"><img src="https://user-images.githubusercontent.com/2105067/35464215-a014d512-02a9-11e8-8913-63a066f6064e.png" alt="FSM" width="350px" align="center;"/></p></a>

# Dynamo Store

This package is a [DynamoDB](https://aws.amazon.com/dynamodb) implementation of a [fsm](https://github.com/fsm/fsm).[Store](https://github.com/fsm/fsm/blob/master/fsm.go#L26-L29).

## Environment Variables

When using this store, you must set four environment variables:

```sh
DYNAMO_REGION=""
DYNAMO_ACCESS_KEY_ID=""
DYNAMO_SECRET_ACCESS_KEY=""
DYNAMO_TABLE=""
```

## Getting Started

> Note: The environment variables above are assumed to be set in this example code:

```go
package main

import "github.com/fsm/dynamo-store"

func main() {
    store := dynamostore.New()
    // ...
}
```

## License

[MIT](LICENSE.md)