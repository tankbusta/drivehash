# Drivehash

A quick tool to hash a folder or an entire file system to generate a list of MD5, SHA1, and SHA256. The results are split into 3 files for each hash type.

This was primarily developed for generating [Redline](https://www.fireeye.com/services/freeware/redline.html) whitelists on gold images.

## Usage

```
  -backlog-processing int
        The maximum number of filepaths in the processing backlog at any given time (default 100)
  -backlog-writer int
        The maximum number of hashes in the writer backlog at any given time (default 
250)
  -ignore-admin
        Ignore the check if admin flag
  -result-dir string
        the directory to save the results (default "./results")
  -start-dir string
        the directory or drive path to walk
```

## Requirements

* Go 1.15+
* Windows/Mac/Linux

## Building

    go build -o drivehash.bin cmd\drivehash\drivehash.go

## License

MIT