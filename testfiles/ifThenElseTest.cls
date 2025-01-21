Dim IsVisible As Boolean
Dim MyBool As Boolean

IsVisible = True
MyBool = True

If IsVisible Then 'Block if then else
    Beep
    IsVisible = False
ElseIf MyBool Then
  myVar = True
Else
    MyBool = False
End If

' Inline if then else'
If MyBool Then MyBool = False ' without else'
If MyBool Then myVar = False Else MyBool = False 
'with else'

