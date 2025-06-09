Public Property Get UserName() As String
    UserName = frmLogin.UserName
End Property
Private Sub Class_Initialize()
    '##ReplaceStatement frmLogin = New frmLoginWin
    Dim myVar as Boolean
    Dim IsVisible as Boolean
    Dim MyBool as Boolean

    myVar=true
    IsVisible=true
    MyBool=true
    Set frmLogin = New frmLoginWin
    If IsVisible Then
    Beep
    IsVisible = false
    ElseIf MyBool Then
    myVar = true
    Else
        MyBool = false
    End If
End Sub