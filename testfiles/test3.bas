sub test_3
Dim i As Integer, j As Integer
i = 1

Do While i <= 5
    j = 1
    
    Do While j <= 5
        Debug.Print i & " x " & j & " = " & (i * j)
        j = j + 1
    Loop
    
    Debug.Print "-----"
    i = i + 1
Loop
end sub