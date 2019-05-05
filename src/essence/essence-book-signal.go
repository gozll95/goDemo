# book-signal

1."pipeline执行命令":
func runCmds(cmds []*exec.Cmd) ([]string, error) {
	if cmds == nil || len(cmds) == 0 {
		return nil, errors.New("The cmd slice is invalid!")
	}
	first := true
	var output []byte
	var err error
	for _, cmd := range cmds {
		fmt.Printf("Run command: %v\n", getCmdPlaintext(cmd))
		if !first {
			var stdinBuf bytes.Buffer
			stdinBuf.Write(output)
			cmd.Stdin = &stdinBuf
		}
		var stdoutBuf bytes.Buffer
		cmd.Stdout = &stdoutBuf
		if err = cmd.Start(); err != nil {
			return nil, getError(err, cmd)
		}
		if err = cmd.Wait(); err != nil {
			return nil, getError(err, cmd)
		}
		output = stdoutBuf.Bytes()
		//fmt.Printf("Output:\n%s\n", string(output))
		if first {
			first = false
		}
	}
	var lines []string
	var outputBuf bytes.Buffer
	outputBuf.Write(output)
	for {
		line, err := outputBuf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, getError(err, nil)
			}
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}

2."bytes.Buffer"
var xx bytes.Buffer
xx.Write(xx)
for{
	line,err:=xx.ReadBytes('\n')
	lines=append(lines,string(line))
}

3."pids":
proc, err := os.FindProcess(pid)
err = proc.Signal(sig)
