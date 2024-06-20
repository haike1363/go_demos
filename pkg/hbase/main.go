package main

import (
	"context"
	"fmt"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"os"
	"strings"
)

func main() {

	zkRoot := gohbase.ZookeeperRoot(os.Getenv("ZK_ROOT"))
	var client gohbase.Client
	if strings.HasPrefix(os.Getenv("AUTH"), "KERBEROS") {
		auth := gohbase.Auth(os.Getenv("AUTH"))
		user := gohbase.EffectiveUser(os.Getenv("KRB5_USER"))
		client = gohbase.NewClient(os.Getenv("ZK_QUORUM"), []gohbase.Option{auth, user, zkRoot}...)
	} else {
		client = gohbase.NewClient(os.Getenv("ZK_QUORUM"), []gohbase.Option{zkRoot}...)
	}

	// Values maps a ColumnFamily -> Qualifiers -> Values.
	values := map[string]map[string][]byte{
		os.Getenv("CF"): {"a": []byte{0}}}
	putRequest, err := hrpc.NewPutStr(context.Background(), os.Getenv("TABLE"), "row1", values)
	rsp, err := client.Put(putRequest)
	fmt.Println(rsp)
	fmt.Println(err)
}
