## Steps to Build and Run the Tool

If the project is already cloned locally or zip file is present:

To build and run the VB6 to C#.NET Migration Tool using a Go binary on Windows, where VB6 files are provided as a path through the command line, you can follow these steps. The instructions assume you have already built the Go binary and that the Go tool is available.

---
### 1. Build the Go Binary (if not already built)

Since we're using make running (cmake needs to be present on system if want to build):

```
make build
```

It runs following command:

```
go build -o ./dist/vb2sharp.exe main.go
```

---
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


---
