import os
import subprocess

def convert_all_testfiles(testfiles_dir, output_dir, main_go_path):
    os.makedirs(output_dir, exist_ok=True)
    
    for testfile in os.listdir(testfiles_dir):
        testfile_path = os.path.join(testfiles_dir, testfile)
        if not os.path.isfile(testfile_path):
            continue
        
        try:
            print(f"Converting {testfile}...")
            result = subprocess.run(
                ["go", "run", main_go_path, testfile_path],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True
            )
            if result.returncode != 0:
                print(f"Error converting {testfile}: {result.stderr}")
                continue

            print(f"Successfully converted {testfile}. Output is in the {output_dir} directory.")
        except Exception as e:
            print(f"Failed to convert {testfile} due to error: {e}")
            
def compare_files(expected_file, output_file):
    with open(expected_file, 'r') as ef, open(output_file, 'r') as of:
        expected_content = ef.read().strip()
        output_content = of.read().strip()
    return expected_content == output_content

def compare_test_cases(expected_dir, output_dir):
    passed_count = 0
    total_count = 0

    if not os.path.exists(expected_dir):
        print(f"Expected output directory '{expected_dir}' does not exist.")
        return
    if not os.path.exists(output_dir):
        print(f"Output directory '{output_dir}' does not exist.")
        return
    
    for expected_file in os.listdir(expected_dir):
        expected_file_path = os.path.join(expected_dir, expected_file)
        output_file_path = os.path.join(output_dir, expected_file)

        if not os.path.isfile(expected_file_path):
            continue

        total_count += 1
        
        if os.path.isfile(output_file_path):
            if compare_files(expected_file_path, output_file_path):
                print(f"Test case '{expected_file}' passed.")
                passed_count += 1
            else:
                print(f"Test case '{expected_file}' failed (content mismatch).")
        else:
            print(f"Test case '{expected_file}' failed (missing in output).")
            
    print(f"\nSummary: {passed_count}/{total_count} test cases passed.")

if __name__ == "__main__":

    TESTFILES_DIR = "testfiles"  
    OUTPUT_DIR = "output"  
    EXPECTED_DIR= "expected_output"
    MAIN_GO_PATH = "."
    
    convert_all_testfiles(TESTFILES_DIR, OUTPUT_DIR, MAIN_GO_PATH)
    compare_test_cases(EXPECTED_DIR, OUTPUT_DIR)