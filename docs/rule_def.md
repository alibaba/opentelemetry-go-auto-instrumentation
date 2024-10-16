## Fields of a rule definition

## Instument a function
- `ImportPath`: The import path of the package that contains the function to be instrumented.
- `Function`: The name of the function to be instrumented.
- `OnEnter`: The name of the function to be called when the instrumented function is called.
- `OnExit`: The name of the function to be called when the instrumented function returns.
- `Order`: The order of the probe code in the instrumented function.
- `Path`: The path to the directory containing the probe code.


## Add a new file during compiling package
- `ImportPath`: The import path of the package that contains the function to be instrumented.
- `FileName` : The name of the file to be added.
- `Path`: The path to the directory containing the probe code.

## Add a new field to a struct
- `ImportPath`: The import path of the package that contains the struct to be instrumented.
- `StructType`: The name of the struct to be instrumented.
- `FieldName`: The name of the field to be added.
- `FieldType`: The type of the field to be added.