# Datatypes

## Simple types
| Datatype | Description                        | Status      | 
|----------|------------------------------------|-------------|
| Integer  | integer 64 bit                     | in progress | 
| Float    | float 64 bit                       | in progress |
| BigInt   | large integers > 64 bit            | planned     |
| BigFloat | large floats > 64 bit              | planned     |
| Date     | date object without time component | planned     |
| Time     | object for time                    | planned     |
| Datetime | date+time                          | planned     |


## File based
| Datatype   | Description                        | Status  |
|------------|------------------------------------|---------|
| file small | small files ([]byte/base64 string) | planned |
| file large | large file, using io.reader        | planned |


## Combined datatypes
| Datatype | Description                               | Status  |
|----------|-------------------------------------------|---------|
| table    | contains different datatypes as subfields | planned |