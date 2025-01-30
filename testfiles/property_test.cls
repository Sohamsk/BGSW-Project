Private mName As String
Private mOwner As Person

' Name property (String) - uses Property Get and Property Let
Public Property Get Name() As String
    Name = mName
End Property

Public Property Let Name(value As String)
    mName = value
End Property

' Owner property (Object) - uses Property Get and Property Set
Public Property Get Owner() As Person
    Set Owner = mOwner
End Property

Public Property Set Owner(value As Person)
    Set mOwner = value
End Property
