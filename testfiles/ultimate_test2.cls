' VBA Class Module: Operations.cls
Option Explicit

' Declare class-level variables
Private count As Integer
Private statusCode As Integer

' Function to subtract two numbers
Public Function SubtractNumbers(a As Integer, b As Integer) As Integer
    SubtractNumbers = a - b
End Function

' Function to multiply two numbers
Public Function MultiplyNumbers(a As Integer, b As Integer) As Integer
    MultiplyNumbers = a * b
End Function

' Subroutine to perform calculations
Public Sub PerformCalculations(ByVal x As Integer, ByVal y As Integer)
    Dim result As Integer
    result = MultiplyNumbers(x, y)
End Sub

' Subroutine demonstrating loops and conditionals (without flag)
Public Sub ProcessData()
    Dim i As Integer
    Dim j As Integer
    Dim total As Integer
    total = 10  

    ' Conditional checks without using a flag variable
    If total < 10 Then
        count = 1
    ElseIf total = 10 Then
        count = 2
    Else
        count = 3
    End If

    ' Nested loop example
    For i = 1 To 4
        For j = 2 To 8 Step 2
            total = total + j
        Next j
    Next i

    ' Using the subtraction function
    count = SubtractNumbers(100, 25)

    ' Do-While loop
    Do While count > 0
        count = count - 10
    Loop

    ' Assign status code
    If count <= 0 Then
        statusCode = 200
    Else
        statusCode = 500
    End If
End Sub
