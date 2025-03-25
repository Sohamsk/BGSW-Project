import os
import subprocess
import glob
import shutil
import sys
import platform
import argparse

def convert_vb6_to_csharp(vb6_file):
    """Convert a single VB6 file to C# using the Go converter"""
    print(f"Converting {os.path.basename(vb6_file)} to C#...")
    
    try:
        # Call your Go program with the VB6 file as input
        result = subprocess.run(
            ["go", "run", ".", vb6_file],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            check=True
        )
        print(f"Conversion output: {result.stdout}")
        return True
    except subprocess.CalledProcessError as e:
        print(f"Conversion failed: {e.stderr}")
        return False
    except Exception as e:
        print(f"Error during conversion: {str(e)}")
        return False

def process_vb6_directory(input_dir):
    """Process all VB6 files in the input directory"""
    # Find all VB6 files (adjust extensions as needed)
    vb6_files = []
    for ext in ['.vb', '.bas', '.cls', '.frm']:
        vb6_files.extend(glob.glob(os.path.join(input_dir, f"*{ext}")))
    
    print(f"Found {len(vb6_files)} VB6 files to process")
    success_count = 0
    
    for vb6_file in vb6_files:
        if convert_vb6_to_csharp(vb6_file):
            success_count += 1
    
    print(f"Converted {success_count}/{len(vb6_files)} VB6 files to C#")
    return success_count > 0

def check_output_dir(output_dir):
    """Check if the output directory exists and has any .cs files"""
    if not os.path.exists(output_dir):
        print(f"Output directory '{output_dir}' does not exist.")
        return False
    
    cs_files = glob.glob(os.path.join(output_dir, "*.cs"))
    if not cs_files:
        print(f"No C# files found in the output directory '{output_dir}'.")
        return False
    
    return True
'''
def check_dotnet_installed():
    """Check if dotnet CLI is installed and working"""
    try:
        result = subprocess.run(
            ["dotnet", "--version"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        if result.returncode == 0:
            print(f"Using .NET version: {result.stdout.strip()}")
            return True
        return False
    except:
        return False
'''
def check_csc_installed():
    """Check if csc compiler is installed and working"""
    # Different command based on OS
    if platform.system() == "Windows":
        csc_cmd = "csc"
    else:
        csc_cmd = "csc" 
    
    try:
        result = subprocess.run(
            [csc_cmd, "/help"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        if result.returncode == 0:
            print("C# compiler (csc) found")
            return True, csc_cmd
        return False, csc_cmd
    except FileNotFoundError:
        try:
            # Try with mono on Unix systems
            if platform.system() != "Windows":
                result = subprocess.run(
                    ["mcs", "--version"],
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE,
                    text=True
                )
                if result.returncode == 0:
                    print("Mono C# compiler (mcs) found")
                    return True, "mcs"
        except:
            pass
        return False, csc_cmd
'''
def direct_syntax_check(cs_file, temp_dir):
    """Perform basic syntax checking using a simple compilation approach"""
    file_name = os.path.basename(cs_file)
    print(f"Basic syntax checking {file_name}...")
    
    # Read the file content
    try:
        with open(cs_file, 'r', encoding='utf-8', errors='replace') as f:
            content = f.read()
            
        # Check for basic syntax issues
        has_class = "class" in content
        has_namespace = "namespace" in content
        has_main = "Main" in content
        has_braces = "{" in content and "}" in content
        
        # Very basic validation
        if has_braces and (has_class or has_namespace):
            # Create a simple report
            report = f"""
Basic syntax check for {file_name}:
- Contains class declaration: {has_class}
- Contains namespace: {has_namespace}
- Contains Main method: {has_main}
- Contains proper braces: {has_braces}

Note: This is a very basic check and doesn't guarantee the code will compile.
            """
            return True, report
        else:
            return False, "Basic syntax check failed: Missing essential C# structure"
            
    except Exception as e:
        return False, f"Error during basic syntax check: {str(e)}"


def compile_cs_file_with_dotnet(cs_file, temp_dir):
    """Compile a single C# file using dotnet CLI"""
    file_name = os.path.basename(cs_file)
    base_name = os.path.splitext(file_name)[0]
    temp_project_dir = os.path.join(temp_dir, base_name)
    
    print(f"Compiling {file_name}...")
    
    # Clean up any previous compilation attempt for this file
    if os.path.exists(temp_project_dir):
        shutil.rmtree(temp_project_dir)
    
    # Create a new temp project directory
    os.makedirs(temp_project_dir)
    
    try:
        # Create a new project with specific template
        result = subprocess.run(
            ["dotnet", "new", "classlib", "--force", "-o", temp_project_dir],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            timeout=60  # Add timeout
        )
        
        if result.returncode != 0:
            return False, f"Failed to create project: {result.stderr}"
        
        # Read the C# file content
        with open(cs_file, 'r', encoding='utf-8', errors='replace') as source_file:
            source_content = source_file.read()
        
        # Find the default class file created by dotnet new
        class_files = glob.glob(os.path.join(temp_project_dir, "*.cs"))
        if not class_files:
            return False, "Could not find default class file in new project"
        
        class_path = class_files[0]
        
        # Replace the default class file with our content
        with open(class_path, 'w', encoding='utf-8') as target_file:
            target_file.write(source_content)
        
        # Build the project
        build_result = subprocess.run(
            ["dotnet", "build", temp_project_dir, "-c", "Release"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            timeout=60  # Add timeout
        )
        
        if build_result.returncode == 0:
            return True, build_result.stdout
        else:
            return False, f"Compilation errors:\n{build_result.stderr}\n{build_result.stdout}"
            
    except subprocess.TimeoutExpired:
        return False, "Command timed out"
    except subprocess.CalledProcessError as e:
        return False, f"Process error: {e}"
    except Exception as e:
        return False, f"Error: {str(e)}"
'''
def compile_cs_with_roslyn(cs_file, csc_cmd):
    """Try to compile using Roslyn compiler directly"""
    file_name = os.path.basename(cs_file)
    print(f"Syntax checking {file_name} with {csc_cmd}...")
    
    try:
        # Use the appropriate command for the C# compiler
        if csc_cmd == "mcs":
            # Mono C# compiler syntax
            cmd = [csc_cmd, "-out:/dev/null", cs_file]
        else:
            # Microsoft C# compiler syntax
            cmd = [csc_cmd, "/nologo", "/t:library", "/out:NUL", cs_file]
        
        result = subprocess.run(
            cmd,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            timeout=30  # Add timeout
        )
        
        if result.returncode == 0:
            return True, "Syntax check passed"
        else:
            return False, f"Syntax errors:\n{result.stderr}\n{result.stdout}"
    
    except subprocess.TimeoutExpired:
        return False, "Command timed out"
    except FileNotFoundError:
        return False, f"{csc_cmd} compiler not found"
    except Exception as e:
        return False, f"Error: {str(e)}"

def test_cs_compilation(output_dir, temp_dir):
    """Test compilation of all C# files in the output directory"""
    if not check_output_dir(output_dir):
        return
    
    # Make sure the temp directory exists
    os.makedirs(temp_dir, exist_ok=True)
    
    # Check for compilers
    csc_available, csc_cmd = check_csc_installed()
    
    # Get all C# files
    cs_files = glob.glob(os.path.join(output_dir, "*.cs"))
    total_files = len(cs_files)
    success_count = 0
    
    # Prepare the log file
    log_path = os.path.join(temp_dir, "compilation_results.log")
    with open(log_path, 'w', encoding='utf-8') as log_file:
        log_file.write("C# Compilation Results\n")
        log_file.write("=====================\n\n")
        log_file.write(f"C# compiler available: {csc_available}\n\n")
        
        for i, cs_file in enumerate(cs_files):
            file_name = os.path.basename(cs_file)
            print(f"[{i+1}/{total_files}] Processing {file_name}")

            # Initialize variables for other methods
            roslyn_success = False
            roslyn_message = "Skipped (csc not available)"
            
            # Try with Roslyn if dotnet failed and CSC is available
            if csc_available:
                try:
                    roslyn_success, roslyn_message = compile_cs_with_roslyn(cs_file, csc_cmd)
                except Exception as e:
                    roslyn_message = f"Unexpected error: {str(e)}"
          
            # Determine status
            if roslyn_success:
                status = "SYNTAX CHECK AND COMPILATION SUCCESSFUL"
                success_count += 1
                print(f"✅ {file_name} syntax check passed")
            else:
                status = "COMPILATION FAILED"
                print(f"❌ {file_name} failed to compile")
            
            # Log the results
            log_file.write(f"File: {file_name}\n")
            log_file.write(f"Status: {status}\n")
            
            if roslyn_success:
                log_file.write("Method: csc syntax check\n")
                log_file.write(f"Output:\n{roslyn_message}\n")
            else:
                log_file.write(f"csc errors:\n{roslyn_message}\n")
            
            log_file.write("-----------------------------------\n\n")
    
    print(f"\nCompilation Summary: {success_count}/{total_files} files passed validation.")
    print(f"Detailed compilation results written to {log_path}")

def main():
    """Main entry point"""
    INPUT_DIR = "testfiles"  # Directory containing VB6 files
    OUTPUT_DIR = "output1"
    TEMP_DIR = "compilation_results"
    
    # Step 1: Convert VB6 to C#
    print("\nConverting VB6 files to C#...")
    process_vb6_directory(INPUT_DIR)
    
    # Step 2: Test C# compilation - Always run this regardless of conversion success
    print("\nTesting C# compilation...")
    test_cs_compilation(OUTPUT_DIR, TEMP_DIR)

if __name__ == "__main__":
    main()