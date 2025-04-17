
This project is a prototype of a VB6 to C# source-to-source converter. It utilizes ANTLR for parsing VB6 constructs and transforming them into C# code. The workflow is divided into multiple stages involving event listening, JSON-based intermediate representations, and structured conversion to C#.

---
### 🔹Workflow Components
---
#### 3.1.1 Listener Package

- Files such as ForNextStmt.go, VariableStmt.go, and IfThenElse.go correspond to specific VB6 constructs.
- This package listens to the syntax tree events emitted by ANTLR.
- Converts VB6 constructs to JSON objects for further processing.

#### 3.1.2 Converter Package

####1.converter.go:
  - Reads the generated JSON from the listener package.
  - Converts JSON into mapped Go structs defined in models.go.
  - Invokes specific conversion handlers from funcMap.go to generate equivalent C# code.
  - Outputs the C# code to a file.
####2.funcMap.go:
  - Maintains mappings of VB6 construct names to handler functions. Each handler is responsible for converting its respective construct.
####3.models.go:
  - Defines Go data structures for unmarshaling the JSON representation of VB6 constructs.

---
### 🔹Supported VB6 Constructs

The current implementation supports the following constructs:

#### 1. Control Flow Statements

- If-Then-Else Statement → IfThenElseStmtRule, ElseIfRule, ElseRule
- For-Next Loop → ForNext
- Do-Loop Statement → DoloopStmt
- For Each Loop → ForEachStmt
- With Statement → WithStmt
- Return Statement → ReturnStmt
- Break/Exit Statement → BreakStmt
---
#### 2. Functions & Procedures

- Function Declaration → FuncDecl
- Subroutine Declaration → SubStmt
- Function Call → FuncRule
- Function Arguments → FuncArg, DeclArg
---
#### 3. Variables & Data Types

- Literal Values → Literal
- Variable Declaration → Not explicitly present, but DeclArg could serve a similar role
- Set Statement (Assigning Objects) → SetStmt
- Type Declaration (User-Defined Types - UDTs) → TypeStmt
- Enum Statement → EnumStmt
---
#### 4. Object-Oriented Constructs

 Property Statement (Similar to Get/Set in VB6 Classes) → PropertyStatement

---
#### 5. Expressions & Rules

- Expression Rule → ExpressionRule
- Rule System → Rule (used as a base struct for many constructs)
---
#### 6. Miscellaneous

- File Context (General VB6 File Handling) → FileContext
- Comments in VB6 Code → Comment, MultiLineComment
- Print Statement → PrintStmt
---
### 🔹VB6 to C# Type Mapping

The vb_cs_types map defines the conversion of VB6 data types to C# equivalents.

| VB6 Type  | C# Type  |
| --------- | -------- |
| boolean   | bool     |
| byte      | byte     |
| currency  | decimal  |
| date      | DateTime |
| double    | double   |
| integer   | int      |
| long      | long     |
| object    | object   |
| single    | float    |
| string    | string   |
| variant   | object   |
| byte()    | byte[]   |
| integer() | short[]  |
| long()    | int[]    |
---


### 🔹Limitations

- This is an incomplete implementation, and the conversion is restricted to the constructs listed above.
- More complex VB6 features such as third-party ActiveX components or custom controls are not yet supported.

---

