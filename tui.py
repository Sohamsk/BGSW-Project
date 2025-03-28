import os
import sys
import subprocess
import glob
import platform
import argparse
from textual.app import App, ComposeResult
from textual.widgets import Header, Footer, Button, Static, Select
from textual.screen import Screen
from textual.containers import Container, Vertical

class VB6ConversionApp(App):
    """TUI for VB6 to C# Conversion Tool with Enhanced Output"""
    
    CSS = """
    Screen {
        align: center middle;
        background: $surface;
    }
    
    #main-container {
        width: 90%;
        height: 80%;
        border: tall $background;
        padding: 2 4;
    }
    
    .button-container {
        width: 100%;
        height: auto;
        align: center middle;
        margin-top: 2;
    }
    
    Button {
        margin: 1;
        width: 30;
    }
    
    Select {
        margin: 1;
    }
    
    Static.title {
        text-align: center;
        text-style: bold;
        color: $text;
    }
    
    Static.result {
        margin: 1;
        border: tall $background;
        padding: 1;
        width: 100%;
        height: 100%;
        background: $panel;
        color: $text;
    }
    """
    
    def convert_vb6_to_csharp(self, vb6_file):
        """Convert a single VB6 file to C# and capture full terminal output"""
        try:
            result = subprocess.run(
                ["go", "run", ".", vb6_file],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                check=True
            )
            return result.stdout, True
        except subprocess.CalledProcessError as e:
            return f"Error converting {vb6_file}: {e.stderr}", False
        except Exception as e:
            return f"Unexpected error: {str(e)}", False

    def compose(self) -> ComposeResult:
        yield Header()
        with Container(id="main-container"):
            yield Static("VB6 to C# Conversion Tool", classes="title")
            yield Select(self.get_vb6_sources(), id="vb6-source")
            with Container(classes="button-container"):
                yield Button("Convert", id="convert-btn")
                yield Button("Exit", id="exit-btn")
            yield Static("Results:", classes="result", id="result-output")
        yield Footer()
    
    def get_vb6_sources(self):
        sources = ["> Select Source"]
        sources += [f for f in os.listdir('.') if f.endswith(('.vb', '.bas', '.cls', '.frm'))]
        if os.path.exists('testfiles'):
            sources += [os.path.join('testfiles', f) for f in os.listdir('testfiles') if f.endswith(('.vb', '.bas', '.cls', '.frm'))]
        return [(f, f) for f in sources]
    
    def on_button_pressed(self, event: Button.Pressed) -> None:
        if event.button.id == "exit-btn":
            self.exit()
        elif event.button.id == "convert-btn":
            vb6_source = self.query_one("#vb6-source").value
            result_output = self.query_one("#result-output")
            
            if vb6_source and vb6_source != "> Select Source":
                output, success = self.convert_vb6_to_csharp(vb6_source)
                result_output.update(output)
            else:
                result_output.update("Please select a valid VB6 source.")

def main():
    app = VB6ConversionApp()
    app.run()

if __name__ == "__main__":
    main()