Public Sub DisplayPersonDetails(ByVal person As Person)
    With person
        Debug.Print "Name: " & .Name
        Debug.Print "Age: " & .Age
        Debug.Print "Role: " & .Role
    End With
End Sub

Sub Main()
    Dim myCar As Object
    Set myCar = New Car

    ' Using With statements
    With myCar
        .Brand = "Toyota"
        .Model = "Corolla"

        ' Nested With for Engine properties
        With .Engine
            .HorsePower = 132
            .FuelType = "Gasoline"
            .Start
        End With
    End With

    ' Display results
    MsgBox "Car: " & myCar.Brand & " " & myCar.Model & vbCrLf & "Engine: " & myCar.Engine.HorsePower & " HP, " & myCar.Engine.FuelType
End Sub

