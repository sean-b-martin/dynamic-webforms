# Dynamic-WebForms

## About

The aim of this library is to provide a simple API for creating, modifying and validating complex user-generated web forms,
which can be defined using a structured data format (JSON). This allows creating and modifying web forms without
needing to create new HTTP endpoints or additional source code for each form.

A form consists of the following components:

- Sections
- Fields
- Subfields

A form can contain 0 to n sections. A section itself can contain 0 to n subsections and 0 to n fields.
A field can be either a simple input field such as a string, a checkbox or a complex input such as file uploads or
dynamic tables. Subfields within fields are used to implement complex types, like columns in a table.

The data that defines a form is called form schema. Each type of field allows different validation rules to be set.
When a user fills out a form, the entered data can be validated against the rules set in the schema. The validation
rules can be either simple, like a number field where the entered number must be between "min $A$ and max $B$" or complex, for
example, "when field $X$ has the value $Y$, then set field $Z$ to disabled".

## Motivation

For the student research project during my dual study program at DHBW Mosbach, I was tasked with
creating a library to handle complex dynamic web forms. The initial implementation was functional, but limited,
due to time constraints. Additionally, it revealed some design issues that this library aims to address.
The original implementation was also written in a different programming language.