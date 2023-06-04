package utils

type Input struct {
	Action     string
	Ip         string
	Port       string
	Encryption string
	Username   string
}

func Format(args []string) Input {

	input := Input{}

	for index, element := range args {
		switch element {
		case "-l":
			input.Action = "listen"
		case "-c":
			input.Action = "connect"
		case "-h":
			input.Ip = args[index+1]
		case "-p":
			input.Port = args[index+1]
		case "-e":
			input.Encryption = args[index+1]
		case "-u":
			input.Username = args[index+1]
		}
	}

	return input
}
