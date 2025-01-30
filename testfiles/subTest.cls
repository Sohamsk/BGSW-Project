Function Mult(ByVal Factor1 As Integer, Factor2 As Integer)
    if Factor1 = 0 then
        Mult = 0
        Exit function
    End if
    Mult = Factor1 * Factor2
End function
