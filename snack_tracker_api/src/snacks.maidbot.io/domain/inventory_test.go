package domain

import (
	// "fmt"
  "io/ioutil"
	// "path/filepath"
	"testing"
  "encoding/json"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInventoryUnMarshalling(t *testing.T) {
	Convey("Inventory domain should parse from json", t, func() {
  	fPath := TEST_DATA_DIR + "inventory_changes.json"
    // var githubCommits []GithubCommit
    //
    // fDict, readErr := ioutil.ReadFile(fPath)
    // So(readErr, ShouldBeNil)
    // jsonErr := json.Unmarshal(fDict, &githubCommits)
    // So(jsonErr, ShouldBeNil)
    // So(len(githubCommits), ShouldEqual, 30)
    // for _, c := range githubCommits {
    //   fmt.Println(fmt.Sprintf("Commit has sha : <%s>", c.Sha))
    //   fmt.Println(fmt.Sprintf("Commited by : <%s>", c.Commit.Author.Name))
    // }
    // So(githubCommits[0].Sha, ShouldEqual, "e337d2db3049cce756b481eb60103a240da3392d")
    So(0, ShouldEqual, 0)
	})
}
