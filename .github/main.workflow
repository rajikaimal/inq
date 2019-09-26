workflow "Build" {
  on = "push"
  resolves = ["Setup Go for use with actions"]
}

action "Setup Go for use with actions" {
  uses = "actions/setup-go@75259a5ae02e59409ee6c4fa1e37ed46ea4e5b8d"
  runs = "go build"
}
