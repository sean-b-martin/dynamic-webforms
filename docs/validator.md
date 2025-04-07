# Validators

## Overview

- **FieldValidator**:
  Bound to a field, controls its behaviour. Example: "value must be less than x"
- **FieldAggregateValidator**:
  Bound to a field, modifies its behavior based on values in other fields.
  Example:
`@onValue: { fieldID: value, true:<behaviour>}`

- **ConditionalValidator**: Returns a boolean based on  its configuration and its provided data. used inside other
  conditional validators. These validators must implement the Conditional interface

## List of Validators

The schema definitions contain different types of attributes. Attributes with a `?` suffix, such as `opt?` are optional
attributes. Attributes or value encapsulated by `<>` are defined as variables.

Additional attributes that are not required for functionality by the validator, such as custom error messages,
are excluded in this documentation.


### General structure

Any validator has the following basic json schema definition:

```json
{
  "name": "<validatorName>",
  "attributes": {
    "<attribute1>": "<value1>"
  }
}
```

#### Conditionals

```go
type Conditional interface {
    Evaluate(data []byte) bool
}
```

### not

```json
{
  "name": "@not",
  "attributes": {
    "name": "<conditionalValidator>",
    "attributes": {}
  }
}
```

returns !Evaluate() of the inner conditional validator

### equals

```json
{
  "name": "@equals",
  "attributes": {
    "static": {
      "value": "<value>", 
      "datatype": "<datatype>"
    },
    "form": {
      "id": 0,
      "index": 0
    },
    "type": "<type>"
  }
}
```

Compares the value of the form element `attributes.form.id[attributes.form.index]` with the value of 
`attributes.static`

`type` defines the type of function used to determine the equality. Examples of types are: 
 - strict: where the datatype provided by the data field must match the value of `static`,
 - string: where both the static value and the form value are converted to strings and compared
 - regex: where the value of `static` must be a regex, the dynamic type must be converted to string and evaluated using
  the regex.

### all

```json
{
  "name": "@all",
  "attributes": {
    "conditions": [{}],
    "true": {},
    "false": {}
  }
}
```

Uses `true` when all conditional validators in `conditions` return true.


### any

```json
{
  "name": "@any",
  "attributes": {
    "conditions": [{}],
    "true": {},
    "false": {}
  }
}
```

Uses `true` when at least one conditional validator in `conditions` returns true.



### if

The `if` validator allows for conditional validation and has the following structure:

```json
{
  "name": "@if",
  "attributes": {
    "conditions": [{}],
    "true": {},
    "false": {}
  }
}
```

`conditions`: a list of conditional validators such as `any`, `not`, `equals` and `all`

