package main

import (
	"context"
	"log"
	"path"
	"time"

	"github.com/google/uuid"

	"google.golang.org/api/option"

	"cloud.google.com/go/datastore"
)

type Task struct {
	Category        string
	Done            bool
	Priority        float64
	Description     string `datastore:",noindex"`
	PercentComplete float64
	Created         time.Time
}

func main() {
	delete()
}

func delete() {
	ctx := context.Background()

	// usr, _ := user.Current()
	p := path.Join("/path/to/credential.json")

	dsClient, err := datastore.NewClient(ctx, "Your Project ID", option.WithCredentialsFile(p))
	if err != nil {
		log.Fatal(err)
	}

	tx, err := dsClient.NewTransaction(ctx)
	if err != nil {
		log.Fatalf("NewTransaction: %v", err)
	}

	const encodedKey = "THE ENCODED KEY"
	k, err := datastore.DecodeKey(encodedKey)
	if err != nil {
		log.Fatalf("DecodeKey: %v", err)
	}

	if err := tx.Delete(k); err != nil {
		log.Fatalf("Delete: %v", err)
	}

	if _, err := tx.Commit(); err != nil {
		log.Fatalf("Commit: %v", err)
	}
}

func update() {
	ctx := context.Background()

	// usr, _ := user.Current()
	p := path.Join("/path/to/credential.json")

	dsClient, err := datastore.NewClient(ctx, "Your Project ID", option.WithCredentialsFile(p))
	if err != nil {
		log.Fatal(err)
	}

	tx, err := dsClient.NewTransaction(ctx)
	if err != nil {
		log.Fatal(err)
	}
	var task Task
	taskKey := datastore.NameKey("Task", "String ID", nil)
	if err := tx.Get(taskKey, &task); err != nil {
		log.Fatal(err)
	}

	log.Println(task)

	task.Priority = 100
	if _, err := tx.Put(taskKey, &task); err != nil {
		log.Fatalf("tx.Put: %v", err)
	}
	if _, err := tx.Commit(); err != nil {
		log.Fatalf("tx.Commit: %v", err)
	}
}

func create() {
	ctx := context.Background()

	// usr, _ := user.Current()
	p := path.Join("/path/to/credential.json")

	dsClient, err := datastore.NewClient(ctx, "Your Project ID", option.WithCredentialsFile(p))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully created a client.")

	// NameKey(Kind, Identifier, ParentKey)
	uu, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("UUID: %v", err)
	}
	taskKey := datastore.NameKey("Task", uu.String(), nil)
	// taskKey := datastore.IncompleteKey("Task", nil)

	log.Println(taskKey)

	task := &Task{
		Category:        "Personal",
		Done:            false,
		Priority:        4,
		Description:     "Learn Cloud Datastore",
		PercentComplete: 10.0,
		Created:         time.Now(),
	}

	_, err = dsClient.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var empty Task
		if err := tx.Get(taskKey, &empty); err != datastore.ErrNoSuchEntity {
			return err
		}
		_, err := tx.Put(taskKey, task)
		/*
			if err != nil {
				log.Println("Something went wrong.")
			}
		*/
		return err
	})

	if err != nil {
		log.Fatalf("RunInTransaction: %v", err)
	}

	log.Println("Successfully created an entity.")
}
