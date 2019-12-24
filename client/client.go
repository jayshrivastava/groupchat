package main

import (
	"context"
	"flag"
	"log"
	"time"
	"io"
	"bufio"
	"os"
	"fmt"
	"strings"
	"sync"

	"google.golang.org/grpc"
	pb "../chat"
)

var (
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
)

func printItem(client pb.ChatClient, key *pb.ItemKey) {
	log.Printf("Getting message at index %d", key.Index)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	message, err := client.GetItem(ctx, key)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	log.Println(message)
}

// printFeatures lists all the features within the given bounding Rectangle.
func printItems(client pb.ChatClient, rng *pb.Range) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListItems(ctx, rng)
	if err != nil {
		log.Fatalf("%v.ListItems(_) = _, %v", client, err)
	}
	for {
		itemValue, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListItems(_) = _, %v", client, err)
		}
		log.Println(itemValue)
	}
}

func main() {
	fmt.Println("starting")
	flag.Parse()
	var wg sync.WaitGroup

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewChatClient(conn)

	// Looking for a valid feature
	printItem(client, &pb.ItemKey{Index: 1})

	//get all items 1 by 1
	// printItems(client, &pb.Range{StartIndex: 2, EndIndex: 9})

	wg.Add(1)
	go process_commands(client, &wg)
	wg.Wait()
}

func SendChatMessage (client pb.ChatClient, chatmsg *pb.ChatMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.SendChatMessage(ctx, chatmsg)
	if err != nil {
		log.Fatalf("Errot Sending Message")
	}
}

func message_reader(wg *sync.WaitGroup) {

}

func process_commands(client pb.ChatClient, wg *sync.WaitGroup) {
	defer wg.Done()

	user := pb.User{}
	group := pb.Group{}
	
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		entry := strings.SplitN(text, " ", 2)
		command := entry[0]
		argument := ""
		
		if len(entry) > 1 {
			argument = entry[1]
		}

		switch command {
			case "register":
				user = pb.User{UserId: 1, Username: argument}
				log.Println("Registered", user.Username)
			case "connect":
				group = pb.Group{GroupId: 1, Groupname: argument}
				log.Println("Connected to", group.Groupname)				
			case "help":
				fmt.Println("Usage:\n")
				fmt.Println("register <username> \n")
				fmt.Println("connect <group chat name> \n")
				fmt.Println("<message> \n")
			default:
				if user.Username != "" && group.Groupname != "" {
					msg := pb.ChatMessage{Text: text, Date: "Some Data", Sender: &user, Group: &group}
					SendChatMessage(client, &msg)
				} else {
					log.Println("Please use 'register <username>' and 'connect <groupname' before sending a message. Type \"help\" for help")
				}
		}
	}
}