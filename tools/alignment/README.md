# Alignment Component
Implementation of the alignment component. It provides functionality to find highly similar sub-sequences of two given sequences.
## Usage
```
$ alignment -h
Start and configure the alignment server component of Picapica project

Usage:
  alignment [command]

Available Commands:
  align       Aligns two given texts
  help        Help about any command
  read        Reads an alignment task from a json file

Flags:
  -h, --help               help for alignment
      --log_level string   level of log verbosity (default "info")

Use "alignment [command] --help" for more information about a command.

```

## Example
To circumvent the maximum character length for CLI input, the `alignment` command is called and then reads a JSON task from STDIN. 
Note that in the example below, the `config` part of the JSON is serialized as string, for internal reasons.
Given an ngram-size of 3 and an overlap of 2, the tool correctly computes the common substring "rather small test", complete with character offsets, the matched text, and surrounding text context.
```shell
$ jobs/alignment read --text=true --context=true
{
  "source": "This is a rather small test text",
  "target": "This is another rather small test sample",
  "config": "{\"seeder\": {\"hash\": {\"n\": 3,\"overlap\": 2}},\"extender\": {\"range\": {\"theta\": 250}}}"
}
[{"source":{"begin":9,"end":27,"text":" rather small test","before":"This is a","after":" text"},"target":{"begin":15,"end":33,"text":" rather small test","before":"This is another","after":" sample"}}]

```