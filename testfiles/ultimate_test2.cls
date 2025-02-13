VERSION 1.0 CLASS
BEGIN
  MultiUse = -1  ' True
END
Attribute VB_Name = "Example"
Attribute VB_GlobalNameSpace = False
Attribute VB_Creatable = True
Attribute VB_PredeclaredId = False
Attribute VB_Exposed = True

' *********************
' Example VB6 Class Module
' Demonstrates various VB6 features
' *********************

Option Explicit

' ------------------------
' Enum Declaration
' ------------------------
Public Enum UserRole
    Admin = 1
    Editor = 2
    Viewer = 3
End Enum

' ------------------------
' User-Defined Type
' ------------------------
Public Type Person
    Name As String
    Age As Integer
    Role As UserRole
End Type

' ------------------------
' Private Variables (Encapsulation)
' ------------------------
Private pName As String
Private pAge As Integer
Private pRole As UserRole

' ------------------------
' Property Get/Let for Name
' ------------------------
Public Property Get Name() As String
    Name = pName
End Property

Public Property Let Name(ByVal value As String)
    pName = value
End Property

' ------------------------
' Property Get/Set for Role (Object)
' ------------------------
Public Property Get Role() As UserRole
    Role = pRole
End Property

Public Property Set Role(value As UserRole)
    pRole = value
End Property

' ------------------------
' Function: Returns Full Details
' ------------------------
Public Function GetDetails() As String
    GetDetails = "Name: " & pName & ", Age: " & pAge & ", Role: " & Role
End Function

' ------------------------
' Sub: Sets User Details
' ------------------------
Public Sub SetUserDetails(ByVal name As String, ByVal age As Integer, ByVal role As UserRole)
    pName = name
    pAge = age
    pRole = role
End Sub

' ------------------------
' If-Then-Else Example
' ------------------------
Public Function IsAdult() As Boolean
    If pAge >= 18 Then
        IsAdult = True
    Else
        IsAdult = False
    End If
End Function

' ------------------------
' Loop Example: For Next
' ------------------------
Public Sub PrintNumbers()
    Dim i As Integer
    For i = 1 To 5
        Debug.Print "Number: " & i
        If i = 3 Then Exit For  ' Exit statement example
    Next i
End Sub

' ------------------------
' Loop Example: Do While
' ------------------------
Public Sub Countdown()
    Dim counter As Integer
    counter = 5
    Do While counter > 0
        Debug.Print "Countdown: " & counter
        counter = counter - 1
        If counter = 2 Then Exit Do ' Exit statement example
    Loop
End Sub

' ------------------------
' With Statement Example
' ------------------------
Public Sub DisplayPersonDetails(ByVal person As Person)
    With person
        Debug.Print "Name: " & .Name
        Debug.Print "Age: " & .Age
        Debug.Print "Role: " & .Role
    End With
End Sub

' ------------------------
' Class Constructor
' ------------------------
Private Sub Class_Initialize()
    pName = "Unknown"
    pAge = 0
    pRole = Viewer
    Debug.Print "Class Initialized!"
End Sub

' ------------------------
' Class Destructor
' ------------------------
Private Sub Class_Terminate()
    Debug.Print "Class Terminated!"
End Sub
