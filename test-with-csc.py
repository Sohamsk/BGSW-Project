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

def test_cs_compilation(output_dir, temp_dir, specific_file=None):
    """Test compilation of C# files"""
    if not check_output_dir(output_dir):
        return
    
    # Make sure the temp directory exists
    os.makedirs(temp_dir, exist_ok=True)
    
    # Check for compilers
    csc_available, csc_cmd = check_csc_installed()
    
    # Get C# files to compile
    if specific_file:
        cs_files = [specific_file]
    else:
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

            # Initialize variables
            roslyn_success = False
            roslyn_message = "Skipped (csc not available)"
            
            # Try with Roslyn if CSC is available
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
    """Main entry point with argument parsing"""
    parser = argparse.ArgumentParser(description="VB6 to C# Conversion and Compilation Tool")
    
    # Mutually exclusive group for conversion targets
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument("-f", "--file", type=str, 
                       help="Convert and compile a single VB6 file")
    group.add_argument("-d", "--directory", action="store_true", 
                       help="Convert and compile all VB6 files in the testfiles directory")
    
    # Optional output and temp directories
    parser.add_argument("-o", "--output", type=str, default="output1", 
                        help="Output directory for converted C# files (default: output1)")
    parser.add_argument("-t", "--temp", type=str, default="compilation_results", 
                        help="Temporary directory for compilation logs (default: compilation_results)")
    
    # Parse arguments
    args = parser.parse_args()
    
    # Perform conversion and compilation based on arguments
    if args.file:
        # Convert and compile single file
        if convert_vb6_to_csharp(args.file):
            # Assume output file is in the output directory with .cs extension
            output_cs_file = os.path.join(args.output, 
                                          os.path.splitext(os.path.basename(args.file))[0] + ".cs")
            test_cs_compilation(args.output, args.temp, specific_file=output_cs_file)
    
    elif args.directory:
        # Convert all VB6 files in testfiles directory
        print("\nConverting VB6 files to C#...")
        process_vb6_directory("testfiles")
        
        # Compile all converted files
        print("\nTesting C# compilation...")
        test_cs_compilation(args.output, args.temp)

if __name__ == "__main__":
    main()