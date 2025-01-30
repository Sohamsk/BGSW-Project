Public Property Get UserName() As String
    UserName = frmLogin.UserName
End Property
Private Sub Class_Initialize()
    '##ReplaceStatement frmLogin = New frmLoginWin

    Set frmLogin = New frmLoginWin
    If IsVisible Then
    Beep
    IsVisible = False
    ElseIf MyBool Then
    myVar = True
    Else
        MyBool = False
    End If
End Sub