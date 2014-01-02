EXE := sxed

exe:
	go build -o $(EXE) ./main

.PHONY: exe
