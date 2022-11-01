package main

import (
	"context"
	"log"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	log.Println("git clone ...")
	repo := client.Git("https://github.com/haton14/dagger-go-example.git")
	src, err := repo.Branch("main").Tree().ID(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("set up go")
	golang := client.Container().From("golang:1.19")
	golang = golang.WithMountedDirectory("/app", src).WithWorkdir("/app")
	golang = golang.
		Exec(dagger.ContainerExecOpts{
			Args: []string{"go", "build"},
		}).
		Exec(dagger.ContainerExecOpts{
			Args: []string{"go", "test", "./..."},
		})

	log.Println("execute ci")
	if _, err := golang.ExitCode(ctx); err != nil {
		panic(err)
	}
}
