# book-pipe-cmd 

1."命令":
cmd:=exec.Command("xxxx")
stdout,err:=cmd.StdoutPipe()
cmd.Start()


2."使用buffer":
var xx bytes.Buffer
for{
	n,err:=stdout.Read([]byte)
	xx.Write(n[:n])
}
fmt.Println(xx.String())

buffer-for-read-write-最后再println

3."不使用buffer":
o:=bufio.NewReader(stdout)
o.ReadLine()

4."io.Pipe()/os.Pipe()":
func inMemorySyncPipe() {
	reader, writer := io.Pipe()
	go func() {
		output := make([]byte, 100)
		n, err := reader.Read(output)
		if err != nil {
			fmt.Printf("Error: Couldn't read data from the named pipe: %s\n", err)
		}
		fmt.Printf("Read %d byte(s). [in-memory pipe]\n", n)
	}()
	input := make([]byte, 26)
	for i := 65; i <= 90; i++ {
		input[i-65] = byte(i)
	}
	n, err := writer.Write(input)
	if err != nil {
		fmt.Printf("Error: Couldn't write data to the named pipe: %s\n", err)
	}
	fmt.Printf("Written %d byte(s). [in-memory pipe]\n", n)
	time.Sleep(200 * time.Millisecond)
}

write的从read里去读


