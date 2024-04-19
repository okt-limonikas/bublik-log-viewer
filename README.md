## To run this example

1. Serve application via http server - `python3 -m http.server`
2. Visit app

## How it works

This app expects to find revelant assets for logs in those locations:

- `json` - folder should contain all log files

  - `node_1_0.json` - should contain run log
  - `node_id<NODE_ID>.json` - should contain log for an id
  - `node_id<NODE_ID>_p2.json` - if log is paginated

- `tree.json` - should contain the logs tree json
