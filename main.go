package main

import (
	"github.com/ChmaraX/notidb/cmd"
	"github.com/ChmaraX/notidb/internal"
)

func main() {
	internal.CreateNotionClient()
	cmd.Execute()
}
