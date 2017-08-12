# STOPPROCESSOR
Stop after emitting a blank event on start
Allow you to put first event and then stop processors as soon as they finish their job.

Permit to launch bitfan with a pipeline and quit when work is done.

## Synopsys


|   SETTING   | TYPE | REQUIRED | DEFAULT VALUE |
|-------------|------|----------|---------------|
| exit_bitfan | bool | false    | true          |


## Details

### exit_bitfan
* Value type is bool
* Default value is `true`

Stop bitfan after stopping the pipeline ?



## Configuration blueprint

```
stopprocessor{
	exit_bitfan => true
}
```
