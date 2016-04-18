package actions
import (
  "os"
  "io/ioutil"
  "path/filepath"
  "strings"
  "fmt"
  "bytes"
  "bufio"
  "container/list"
  "text/template"
)
type Snippet struct{
  NONBLOCKING string
  BLOCKING string
  TAB_BLOCKING string
  TAB_NON_BLOCKING string
}

const en_blocking = "Blocking API"
const ja_blocking = "ブロッキング API"
const cn_blocking = "阻塞 API"
var tab_blocking string

const en_non_blocking = "Non-Blocking API"
const ja_non_blocking = "ノンブロッキング API"
const cn_non_blocking = "非阻塞 API"
var tab_non_blocking string
const tmpl_tab = `
**Swift:**

{% tabcontrol %}

{% tabpage {{.TAB_BLOCKING}} %}
{% highlight swift %}
{{.BLOCKING}}
{% endhighlight %}
{% endtabpage %}

{% tabpage {{.TAB_NON_BLOCKING}} %}
{% highlight swift %}
{{.NONBLOCKING}}
{% endhighlight %}
{% endtabpage %}

{% endtabcontrol %}`
const tmpl_single = `
**Swift:**

` + "```swift" +`
{{.BLOCKING}}
` + "```"

type SyncSnippetAction struct {
    Prefix            string
    Index             int
    SnippetSourceDir  string
    DocTargetDir      string
    IsTrial           bool
}

func NewSyncSnippetAction() *SyncSnippetAction {
  return &SyncSnippetAction{"",-1,"","",false}
}

func (as *SyncSnippetAction) ExecuteAction(path string, f os.FileInfo, err error) (e error) {
  if filepath.Ext(path) != ".swift" || !strings.HasPrefix(f.Name(), as.Prefix){
    return
  }
  dir := filepath.Dir(path)
  var base string
  if strings.HasPrefix(f.Name(),"guides_ab-"){
    base = strings.Replace(f.Name(),"_","/",2)
  }else{
    base = strings.Replace(f.Name(),"_","/",-1)
  }
  r_name := filepath.Join(dir, f.Name())
  filename := strings.Replace(filepath.Join(as.DocTargetDir, base),".swift",".mkd",1)
  fmt.Println(filename)
  f1, err := os.OpenFile(r_name, os.O_RDONLY, 0666)
  if err != nil {
    panic(err)
  }

  defer f1.Close()
  scanner := bufio.NewScanner(f1)
  var x list.List
  isParsingSingle := false
  str := ""
  blocking := ""
  non_blocking := ""
  var snip Snippet
  for scanner.Scan() {

    if strings.Contains(scanner.Text(),"//dummy") ||
    strings.Contains(scanner.Text(),"print(") ||
    strings.Contains(scanner.Text(),"//snippet "){
      if ! strings.Contains(scanner.Text(),"print(\""){
      continue
      }
    }
    if strings.HasPrefix(scanner.Text(),"private func snippet") {
      str = ""
      //log.Println("start")
      if strings.HasSuffix(scanner.Text(),"blocking(){") {
      if strings.HasSuffix(scanner.Text(),"non_blocking(){") {
        non_blocking = " "
        }else{
          snip = Snippet{}
          blocking = " "
        }

        } else{
          isParsingSingle = true
          snip = Snippet{}
          blocking = " "
        }
        continue
    }
    if strings.HasPrefix(scanner.Text(),"}"){
      if blocking != "" {
        snip.BLOCKING = strings.Trim(str,"\n")
        blocking = ""
      }
      if non_blocking != "" {
        snip.NONBLOCKING = strings.Trim(str,"\n")
        non_blocking = ""
        x.PushBack(snip)
      }
      if isParsingSingle {
        isParsingSingle = false
        x.PushBack(snip)
      }
      continue
    }

    str = str +"\n"+ strings.Replace(scanner.Text(),"  ","",1)
  }

  f2, err := os.OpenFile(filename, os.O_RDONLY, 0666)
  if err != nil {
    panic(err)
  }


  scanner = bufio.NewScanner(f2)
  shouldSkip := false
  var element = x.Front()
  str = ""
  scanned := ""
  for scanner.Scan() {
    scanned =scanner.Text()
    if strings.HasPrefix(scanned,"layout:") {
      if strings.HasSuffix(scanned,"en-doc") {
        tab_non_blocking = en_non_blocking
        tab_blocking = en_blocking
      } else if strings.HasSuffix(scanned,"ja-doc") {
        tab_non_blocking = ja_non_blocking
        tab_blocking = ja_blocking
      } else if strings.HasSuffix(scanned,"cn-doc") {
        tab_non_blocking = cn_non_blocking
        tab_blocking = cn_blocking
      }
    }
    //page-id:
    if strings.HasPrefix(scanned,"**Swift:**"){
    shouldSkip = true
    continue
    }

    if str=="" {
      str = str + scanned
      }else if !shouldSkip{
      str = str +"\n"+ scanned
    }
    if scanned == "```" || strings.HasPrefix(scanned,"{% endtabcontrol %}"){
      if element == nil {
        continue
      }
      if shouldSkip == true {
        shouldSkip = false
        snip = element.Value.(Snippet)
        text := writeToString(snip)
        str = str+text
        element = element.Next()
      }
    }
  }
  f2.Close()
  src_dir,_ := os.Getwd()
  if as.IsTrial {
    filename = src_dir+"/test_files/temp.mkd"
  }

  err = ioutil.WriteFile(filename, []byte(str), 0644)
  return
}
func writeToString(snippet Snippet) (result string){
  var tmpl string
  if snippet.NONBLOCKING != "" {
    tmpl = tmpl_tab
    snippet.TAB_BLOCKING = tab_blocking
    snippet.TAB_NON_BLOCKING = tab_non_blocking
    }else{
      tmpl = tmpl_single
    }
    t, err := template.New("person").Parse(tmpl)

    if err == nil {
      buff := bytes.NewBufferString("")
      t.Execute(buff, snippet)
      result = buff.String()
    }
  return
}
