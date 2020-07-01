# ssp
String Slice Parser

SSP parses string slices in a regex-y way.
.
It allows to parse slices containing structured or patterned data and extract that data into variables for further use.

For example, if you have a slice of strings

```go
inputSlice := []string{"er", "er", "er", "er", "yolo", "1,2,3,4"}
```

you can make an expression for it and parse your slice:

```go
rx := SSP.Repeat(SSP.Flag("er"), 2, 200).CaptureString(&theYoloString).CaptureStringSlice(&the1234StringSlice)
rx.Parse(inputSlice)
```

thus making sure it fits the pattern, and extracting the "yolo" string and a list of "1", "2", "3" & "4" strings into variables.