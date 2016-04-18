package main

import (
  "os"
  "github.com/codegangsta/cli"
  "path/filepath"
  "strings"
  "github.com/rezacute/batchfiles/actions"
)
type Snippet struct{
  NONBLOCKING string
  BLOCKING string
  TAB_BLOCKING string
  TAB_NON_BLOCKING string
}
var (
  ext string
  src_dir string
  in_prefix string
  out_prefix string
  dst_dir string
  success_count int
  skip_count int
)
const en_blocking = "Blocking API"
const ja_blocking = "ブロッキング API"
const cn_blocking = "阻塞 API"
var tab_blocking string

const en_non_blocking = "Non-Blocking API"
const ja_non_blocking = "ノンブロッキング API"
const cn_non_blocking = "非阻塞 API"
var tab_non_blocking string

func main() {
  success_count = 0
  skip_count = 0
  app := cli.NewApp()
  app.Name = "batchfiles"
  app.Usage = "Batch File operation CLI"
  app.Commands = []cli.Command{
    {
    Name:        "rename",
    Usage:       "use it to batch rename",
    Description: "This will do batch operation for rename",
    Subcommands: []cli.Command{
      {
      Name:        "files",
      Usage:       "rename files in a folder",
      Description: "rename only files, fill skip directories",
      Flags: []cli.Flag{
        cli.StringFlag{
          Name:  "source",
          Value: "",
          Usage: "The directory source path ",
        },
        cli.StringFlag{
          Name:  "src-prefix",
          Value: "",
          Usage: "The prefix source filename ",
        },
        cli.StringFlag{
          Name:  "add-prefix",
          Value: "",
          Usage: "The prefix dest filename ",
        },
        cli.StringFlag{
          Name:  "src-extension",
          Value: "",
          Usage: "The filter extension source filename ",
        },
      },
      Action: func(c *cli.Context) {
        src_dir = c.String("source")

        if src_dir == ""{
          src_dir,_ = os.Getwd()
        }
        ap := actions.AddPrefixAction{
          c.String("src-prefix"),
          c.String("add-prefix"),
          c.String("src-extension"),
        }
        filepath.Walk(src_dir, ap.ExecuteAction)
      },
    },
  },
},
{
Name:        "sync",
Usage:       "use it to batch synchronize snippet with pattern",
Description: "use '_' as separator",
Subcommands: []cli.Command{
  {
  Name:        "snippet",
  Usage:       "Synch ",
  Description: "rename only files, will skip directories",
  Flags: []cli.Flag{
    cli.StringFlag{
      Name:  "source",
      Value: "",
      Usage: "The directory source path ",
    },
    cli.StringFlag{
      Name:  "prefix",
      Value: "",
      Usage: "The prefix source filename ",
    },
    cli.StringFlag{
      Name:  "base_destination",
      Value: "",
      Usage: "The base destination directory, the default is current ",
      EnvVar: "BASE_DST",
    },
    cli.BoolFlag{
      Name:  "trial",
      Usage: "Set true for trial mode. The result will be a new file in the current directory. ",
    },
  },
  Action: func(c *cli.Context) {
    in_prefix = c.String("prefix")
    src_dir = c.String("source")
    dst_dir = c.String("base_destination")
    if src_dir == ""{
      src_dir,_ = os.Getwd()
    }
    if dst_dir == ""{
      dst_dir,_ = os.Getwd()
    }
    if c.Bool("trial") {

    }
    ap := actions.SyncSnippetAction{
      in_prefix,
      -1,
      src_dir,
      dst_dir,
      c.Bool("trial"),
    }
    filepath.Walk(src_dir, ap.ExecuteAction)
  },
},
},
},
}

app.Run(os.Args)
}

func addPrefix(path string, f os.FileInfo, err error) (e error) {
  if filepath.Ext(path) != ext || !strings.HasPrefix(f.Name(), in_prefix){
    return
  }
  dir := filepath.Dir(path)
  base := filepath.Base(path)
  newname := filepath.Join(dir, out_prefix + base)
  os.Rename(path, newname)
  return
}
