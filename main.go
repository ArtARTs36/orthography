package main

import (
	"context"
	"fmt"
	"github.com/artarts36/orthography/firstname"
)

func main() {
	store, err := firstname.LoadCSVStore("./data/first_names_male.csv")
	if err != nil {
		panic(err)
	}

	finder := firstname.NewStorableFinder(store)

	fmt.Println(finder.Find(context.Background(), "Артем"))
}
