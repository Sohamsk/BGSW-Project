Private Sub cmdClear_Click()
    txtHello.Text = ""
End Sub

Private Sub test7()
    Dim obj As Object
    Dim anotherObject As Object

    Set anotherObject = New AnotherClass
    Set obj = New SomeClass
    Set obj.SomeProperty=AnotherObject
End Sub
