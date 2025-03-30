# VB6 To C# Migration Tool Documentation

## Contents

1. [Introduction](#introduction)
   1. [Problem statement](#problem-statement)
   2. [Milestone-1](#milestone-1-prototype)
   3. [Milestone-2](#milestone-2-comprehensive-code-conversion)
2. [Steps to build and run the tool](#steps-to-build-and-run-the-tool)
3. [Workflow](#workflow)
   1. [Workflow Components](#workflow-components)
   2. [Supported VB6 Constructs](#supported-vb6-constructs)
   3. [VB6 to C# Type Mapping](#vb6-to-c-type-mapping)
   4. [Limitations](#limitations)
4. [System Design](#system-design)
5. [Pending Project Components](#pending-project-components)

## Introduction

### Problem Statement

The VB6 to C#.NET Migration Tool offers an effective solution for updating legacy VB6 applications to modern C# .NET platforms. It includes tools for code analysis. The framework ensures efficient migration, modernises applications for better performance and scalability, and provides a foundation for future growth.

### Milestone 1: Prototype

- **Tasks:**
  - **Project Setup:**
    - Set up the Go project structure.
    - Install necessary dependencies for JSON handling and type mappings.
    - Create initial structs to represent VB6 constructs (e.g., Dim, FuncRule).
  - **Basic Functionality:**
    - Implement basic JSON parsing to load VB6 constructs into Go structs.
    - Write a simple handler for translating VB6 variable declarations to C#.
  - **Testing:**
    - Develop unit tests for variable declarations and JSON parsing.
- **Outcome:** A basic prototype capable of parsing VB6 variable declarations and converting them to C#.

### Milestone 2: Comprehensive Code Conversion

#### Sub-Milestone 1.3.1: Core Parsing and Handling

- **Tasks:**
  - Enhance JSON parsing capabilities.
  - Implement handlers for:
    - Variable Declarations (DeclareVariableRule)
    - Function Calls (FuncCallRule)
  - Write unit tests for these handlers.
- **Outcome:** Reliable parsing and translation for VB6 variables and function calls.

#### Sub-Milestone 1.3.2: Statement Processing

- **Tasks:**
  - Implement handlers for:
    - Expressions (ExpressionRuleHandler)
    - Sub-statements (SubStmtHandler)
    - Function Declarations (FunctionHandler)
  - Develop reusable utility functions, such as:
    - ProcessCondition for conditional parsing.
    - handleBody for nested statement handling.
  - Ensure accurate translation of nested structures.
- **Outcome:** Functional support for processing VB6 statements.

#### Sub-Milestone 1.3.3: Control Flow Structures

- **Tasks:**
  - Implement handlers for:
    - If-Then-Else statements (IfThenElseStmtHandler, ElseIfHandler, ElseHandler).
    - Loops (DoLoopStmtHandler, ForNextRule).
    - Select Case (already implemented; perform additional testing if needed).
  - Ensure control flows are accurately translated to C# syntax.
- **Outcome:** Complete support for VB6 control flow structures.

# -----------------------------------------------------------------

## Steps to Build and Run the Tool

If the project is already cloned locally or zip file is present:

To build and run the VB6 to C#.NET Migration Tool using a Go binary on Windows, where VB6 files are provided as a path through the command line, you can follow these steps. The instructions assume you have already built the Go binary and that the Go tool is available.

### 1. Build the Go Binary (if not already built)

Since we're using make running (cmake needs to be present on system if want to build):

```
make build
```

It runs following command:

```
go build -o ./dist/vb2sharp.exe main.go
```

### 2. Command to Run the Tool with VB6 Path Input

Once the Go binary (vb2sharp.exe) is built, you can use the following command format in your Windows Command Prompt to run the migration tool, providing the path to the VB6 project as input:

```
vb2sharp.exe "C:\path\to\your\vb6\project\cls or bas or frm file"
```

**Command Breakdown:**

- `vb2sharp.exe`: The name of your built Go binary.
- `"C:\path\to\your\vb6\project\cls or bas or frm file"`: Replace it with the actual path of your VB6 source file.

**Example:**
If your VB6 project is located at:  
C:\Users\Sample\Documents\VB6_Projects\MyLegacyApp.bas,

the command would look like this:

```
vb2sharp.exe "C:\Users\Sample\Documents\VB6_Projects\MyLegacyApp.bas"
```

The output will be created inside the output directory from where the executable is called.

## Workflow

This project is a prototype of a VB6 to C# source-to-source converter. It utilizes ANTLR for parsing VB6 constructs and transforming them into C# code. The workflow is divided into multiple stages involving event listening, JSON-based intermediate representations, and structured conversion to C#.

### Workflow Components

#### 3.1.1 Listener Package

- Files such as ForNextStmt.go, VariableStmt.go, and IfThenElse.go correspond to specific VB6 constructs.
- This package listens to the syntax tree events emitted by ANTLR.
- Converts VB6 constructs to JSON objects for further processing.

#### 3.1.2 Converter Package

- **converter.go:**
  - Reads the generated JSON from the listener package.
  - Converts JSON into mapped Go structs defined in models.go.
  - Invokes specific conversion handlers from funcMap.go to generate equivalent C# code.
  - Outputs the C# code to a file.
- **funcMap.go:**
  - Maintains mappings of VB6 construct names to handler functions. Each handler is responsible for converting its respective construct.
- **models.go:**
  - Defines Go data structures for unmarshaling the JSON representation of VB6 constructs.

### Supported VB6 Constructs

The current implementation supports the following constructs:

#### 1. Control Flow Statements

- If-Then-Else Statement → IfThenElseStmtRule, ElseIfRule, ElseRule
- For-Next Loop → ForNext
- Do-Loop Statement → DoloopStmt
- For Each Loop → ForEachStmt
- With Statement → WithStmt
- Return Statement → ReturnStmt
- Break/Exit Statement → BreakStmt

#### 2. Functions & Procedures

- Function Declaration → FuncDecl
- Subroutine Declaration → SubStmt
- Function Call → FuncRule
- Function Arguments → FuncArg, DeclArg

#### 3. Variables & Data Types

- Literal Values → Literal
- Variable Declaration → Not explicitly present, but DeclArg could serve a similar role
- Set Statement (Assigning Objects) → SetStmt
- Type Declaration (User-Defined Types - UDTs) → TypeStmt
- Enum Statement → EnumStmt

#### 4. Object-Oriented Constructs

- Property Statement (Similar to Get/Set in VB6 Classes) → PropertyStatement

#### 5. Expressions & Rules

- Expression Rule → ExpressionRule
- Rule System → Rule (used as a base struct for many constructs)

#### 6. Miscellaneous

- File Context (General VB6 File Handling) → FileContext
- Comments in VB6 Code → Comment, MultiLineComment
- Print Statement → PrintStmt

### VB6 to C# Type Mapping

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

### Limitations

- This is an incomplete implementation, and the conversion is restricted to the constructs listed above.
- More complex VB6 features such as third-party ActiveX components or custom controls are not yet supported.

## System Design

- The listener package is tightly coupled with the ANTLR parser to handle VB6 syntax and semantics, generating intermediate JSON outputs for constructs.
- The converter package ensures proper translation to C#, leveraging function maps for modular handling of constructs.

## Pending Project Components

1. File handling constructs for mapping to C# library.
2. Third Party COM component replacement choice to User.
3. UI code context as commented code in C# file.
4. Conversion of Entire Projects utilizing file to file conversion inside a loop, resolving references.
5. Terminal UI to run the application with user-friendly flags.

By:

Sameer Godse Soham Kudalkar Piyush Shah Mohit Barade

Under the Guidance of:  
Prof. Mrs. Jayaprabha M. Kanase,  
Associate Professor,  
Progressive Education Society's Modern College Of Engineering, Pune
