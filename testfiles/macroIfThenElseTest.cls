
#If Win32 Then
    MsgBox "Running on a 32-bit Windows system"
#ElseIf Win64 Then
    MsgBox "Running on a 64-bit Windows system"
#Else
    MsgBox "Unknown system architecture"
#End If
