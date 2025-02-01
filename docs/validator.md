# Validators

## Overview

- **DatatypeValidator**:
Bound to a datatype, contains basic constraints such as minimum string length or allowed value ranges.

- **SectionValidator**:
Bound to sections, controls the behavior of subsections and fields. enables or disables fields and subsections based 
on specific rules, for example `@equals: {fields:[{fieldID:dataIndex},...], onSuccess:<behaviour>}`

- **FieldAggregateValidator**:
  Bound to a field, modifies its behavior based on values in other fields.
  Example:
`@onValue: { fieldID: value, onSuccess:<behaviour>}`

