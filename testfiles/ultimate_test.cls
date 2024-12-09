
Dim count As Integer
Dim result As String
Dim demo as Integer
Function AddNumbers(a As Integer, b As Integer) As Integer
    AddNumbers = a + b
End Function

Sub Mult(ByVal Factor1 As Integer, Factor2 As Integer)
	Dim Product
  Product = Factor1 * Factor2
End Sub


Sub Demo()
    Dim i As String
    Dim sum As Integer
    sum = 0
    
    If sum = 0 Then
      demo= 1
    ElseIf sum > 0 Then
      demo = 2
    Else
      demo = 3
    End If
    For i = 1 To 3
      For j = 1 To 6 Step 2
        sum += 1
      Next j
    Next i
    count = AddNumbers(10, 20)

    Do While count > 0
        count = count - 1
    Loop

End Sub
