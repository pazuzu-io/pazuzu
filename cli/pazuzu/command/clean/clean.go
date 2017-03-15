package clean

import (
  "fmt"
  "os"
  "github.com/urfave/cli"
  "github.com/zalando-incubator/pazuzu"
  "github.com/zalando-incubator/pazuzu/shared"
)

var Command = cli.Command{
  Name: "clean",
  Usage: "Remove Pazuzufile and Dockerfile",
  Action: cleanAction,
}

func cleanAction(c *cli.Context) error {
  err := os.Remove(pazuzu.PazuzufileName)
  if err != nil {
    fmt.Println(err)
  }
  err = os.Remove(pazuzu.DockerfileName)
  if err != nil {
    fmt.Println(err)
  }
  err = os.Remove(shared.TestSpecFilename)
  if err != nil {
    fmt.Println(err)
  }
  return nil
}
