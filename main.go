package main

import (
	"log"
	"time"
)

func main() {
	node1 := NewNode("node 1", ":8001")
	go node1.Serve()
	time.Sleep(time.Millisecond * 100)

	node2 := NewNode("node 2", ":8002")
	go node2.Serve()
	time.Sleep(time.Millisecond * 100)

	log.Println("Client interacting with node 1")
	client1, err := NewClient(":8001")

	if err != nil {
		log.Fatalf("Failed to create client for Node 1: %v", err)
	}
	defer client1.Close()

	err = client1.Put("name", "Akash")
	if err != nil {
		log.Printf("Error putting 'name': %v", err)
	} else {
		log.Println("Put 'name: Akash' successful on Node 1")
	}

	err = client1.Put("city", "Los Santos")
	if err != nil {
		log.Printf("Error putting 'city': %v", err)
	} else {
		log.Println("Put 'city: Los Santos' successful on Node 1")
	}

	val, err := client1.Get("name")
	if err != nil {
		log.Printf("Error getting 'name': %v", err)
	} else {
		log.Printf("Got 'name': %s from Node 1", val)
	}

	val, err = client1.Get("country")
	if err != nil {
		log.Printf("Error getting 'country': %v", err)
	} else {
		log.Printf("Got 'country': %s from Node 1", val)
	}

	log.Println("Client interacting with node 2")

	client2, err := NewClient(":8002")
	if err != nil {
		log.Fatalf("Failed to create client for Node 2: %v", err)
	}
	defer client2.Close()

	val, err = client2.Get("name") // This key was put on node 1
	if err != nil {
		log.Printf("Error getting 'name' from Node 2 (expected): %v", err)
	} else {
		log.Printf("Got 'name': %s from Node 2", val)
	}

	err = client2.Put("language", "Go")
	if err != nil {
		log.Printf("Error putting 'language': %v", err)
	} else {
		log.Println("Put 'language: Go' successful on Node 2")
	}

	val, err = client2.Get("language")
	if err != nil {
		log.Printf("Error getting 'language' from Node 2: %v", err)
	} else {
		log.Printf("Got 'language': %s from Node 2", val)
	}

	log.Println("Operations complete")
	select {}

}
