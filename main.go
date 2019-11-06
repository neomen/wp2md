package main

import (
  "fmt"
  "log"
  "os"

  "github.com/urfave/cli"
  "github.com/frankbille/go-wxr-import"
)

func main() {
  app := cli.NewApp()
  app.Name = "wp2md"
  app.Version = "0.1.0"
  app.Usage = "WXR to MD content"
  app.Action = func(c *cli.Context) error {
    if fileExists(c.Args().Get(0)) {
      wxrXmlData, _ := ioutil.ReadFile(c.Args().Get(0))
      wxr := ParseWxr(wxrXmlData)
    }else{
      fmt.Printf("file not Exists %q\n", c.Args().Get(0))
    }


    log.Print(wxr)
    fmt.Printf("Hello %q\n", c.Args().Get(0))
    fmt.Println("Hello friend!\n")
    fmt.Printf("Hello %q\n", c.String("out"))

    return nil
  }
  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name:  "out, o",
      Value: "./content",
      Usage: "Output content `DIR`",
    },
  }
  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
func fileExists(filename string) bool {
  info, err := os.Stat(filename)
  if os.IsNotExist(err) {
      return false
  }
  return !info.IsDir()
}
func writeFile(fileName string, node *model.JoinedNodeDataBody, alias string, tags []string, menus []*model.JoinedMenu, emvideos []model.Emvideo) {
	file, err := os.Create(fileName)
	util.CheckErrFatal(err, "create", fileName)

	w := bufio.NewWriter(file)
	writeFrontMatter(w, node, alias, tags, menus)
	writeContent(w, node, emvideos)
	w.Flush()
	file.Close()
}

func writeFrontMatter(w io.Writer, node *model.JoinedNodeDataBody, alias string, tags []string, menus []*model.JoinedMenu) {
	created := time.Unix(node.Created, 0).Format("2006-01-02")
	changed := time.Unix(node.Changed, 0).Format("2006-01-02")
	fmt.Fprintln(w, "---")
	fmt.Fprintf(w, "title:       \"%s\"\n", node.Title)
	//fmt.Fprintf(w, "description: \"%s\"\n", node.BodySummary)
	fmt.Fprintf(w, "type:        %s\n", node.Type)
	fmt.Fprintf(w, "date:        %s\n", created)
	if changed != created {
		fmt.Fprintf(w, "changed:     %s\n", changed)
	}
	//fmt.Fprintf(w, "weight:      %d\n", node.Nid) // the node-id is normally ascending in Drupal and is always unique
	fmt.Fprintf(w, "draft:       %v\n", !node.Published)
	fmt.Fprintf(w, "promote:     %v\n", node.Promote)
	fmt.Fprintf(w, "sticky:      %v\n", node.Sticky)
	//fmt.Fprintf(w, "deleted:     %v\n", node.Deleted)
	fmt.Fprintf(w, "url:         /%s\n", alias)
	fmt.Fprintf(w, "aliases:     [ node/%d ]\n", node.Nid)
	if len(menus) > 0 {
		fmt.Fprintf(w, "menu:        [ \"%s\" ]\n", flattenMenuNames(menus))
	}
	for _, tag := range tags {
		fmt.Fprintf(w, "%s\n", tag)
	}
}

func writeContent(w io.Writer, node *model.JoinedNodeDataBody, emvideos []model.Emvideo) {
	if node.BodySummary != "" {
		fmt.Fprintf(w, "\n# Summary:\n")
		for _, line := range strings.Split(node.BodySummary, "\n") {
			fmt.Fprintf(w, "# %s\n", line)
		}
	}

	fmt.Fprintln(w, "\n---")
	body := node.BodyValue
	if strings.HasPrefix(body, node.BodySummary) {
		body = body[len(node.BodySummary):]
		fmt.Fprintln(w, node.BodySummary)
		fmt.Fprintln(w, "<!--more-->")
	}
	for _, emvideo := range emvideos {
		fmt.Fprintf(w, "{{< %s %s >}}", emvideo.Provider, emvideo.VideoId)
	}
	fmt.Fprintln(w, body)
}