Sub Main()
    Dim myCar As Car
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

