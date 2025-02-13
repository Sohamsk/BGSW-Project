Public Event Completed(timeStamp As String)

Public Sub DoWork()
    'Do work
    RaiseEvent Completed(Now)
End Sub
