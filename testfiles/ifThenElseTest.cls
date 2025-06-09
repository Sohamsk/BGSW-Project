sub something()
Dim IsVisible As Boolean
Dim MyBool As Boolean
Dim myVar As Boolean

IsVisible = true
MyBool = true

If IsVisible Then
    Beep
    IsVisible = false
ElseIf MyBool Then
  myVar = true
Else
    MyBool = false
End If
end sub
