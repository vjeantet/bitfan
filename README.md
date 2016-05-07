# logstack

Logstash in GO.

## Usage

```./logstack -f debug.conf```

### TODO

- [x] parse logstash config file
- [x] generic input support
- [x] generic filter support
- [x] generic output support
- [x] configuration condition (if else) support
- [x] dynamic %{field.key} support in config file
- [x] gracefully stop
- [x] gracefully start
- [ ] codec support
- [ ] log to file
- [ ] name all contributors and imported packages


# supported inputs, filters and outputs 
can be found here : https://github.com/veino

## input
* beats
* exec
* file
* stdin
* twitter

## filter
* date
* drop
* grok
* json
* mutate
* split
* uuid

## output
* elasticsearch
* mongodb
* null
* stdout

## Used package
...TODO


