Dim IsVisible As Boolean
Dim MyBool As Boolean

IsVisible = True
MyBool = True

If IsVisible Then
    Beep
    IsVisible = false
ElseIf MyBool Then
    Beep
    Beep
Else
    MyBool = False
End If
