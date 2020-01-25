# inq

Take notes on GitHub with inq CLI

## Install

```
$ go build inq.go
$ cp inq /usr/local/bin
$ inq
```

## Usage

### Configure

- Create a repository to save notes (inq-notes)
- Provide repository name to inq

```
$ inq config [githubRepositoryUrl]
```

### Save note by date

```
$ inq save
```

### Save note by topic

```
$ inq --topic=[topicName] save 
```

### Save note by name

```
$ inq --name=[fileName] save 
```

### Push to GitHub

```
$ inq push
```

MIT
