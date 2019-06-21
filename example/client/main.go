package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"example/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":5000", grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to dial server (%s)\n", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	client := proto.NewKeyValueServiceClient(conn)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		text, _ := reader.ReadString('\n')
		parts := strings.Split(strings.TrimSpace(text), " ")

		if parts[0] == "get" && len(parts) == 2 {
			resp, err := client.Get(context.Background(), &proto.GetRequest{
				Key: parts[1],
			})

			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get payload (%s)\n", err.Error())
			} else {
				fmt.Printf("%s\n", resp.GetValue())
			}

			continue
		}

		if parts[0] == "set" && len(parts) == 3 {
			_, err := client.Set(context.Background(), &proto.SetRequest{
				Key:   parts[1],
				Value: parts[2],
			})

			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to set payload (%s)\n", err.Error())
			}

			continue
		}

		fmt.Fprintf(os.Stderr, "malformed command, must have the form 'get key' or 'set key value'\n")
	}
}
