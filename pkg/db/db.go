package db

import (
	"example/pkg/db/schemas"
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

type DB struct {
	C *gocql.ClusterConfig
	S *gocql.Session
}

func NewDB() *DB {
	/* The example assumes the following CQL was used to setup the keyspace:
	create keyspace example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
	create table example.tweet(timeline text, id UUID, text text, PRIMARY KEY(id));
	create index on example.tweet(timeline);
	*/
	cluster := gocql.NewCluster("localhost:9050")
	cluster.Keyspace = "greetings"
	cluster.Consistency = gocql.Quorum
	// connect to the cluster
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	defer session.Close()

	// ctx := context.Background()

	var g schemas.Greeting
	/* Search for a specific set of records whose 'timeline' column matches
	 * the value 'me'. The secondary index that we created earlier will be
	 * used for optimizing the search */
	if err := session.Query(`SELECT id, name, message FROM greeting WHERE name=? LIMIT 1`,
		"fake").Consistency(gocql.One).Scan(&g.ID, &g.Name, &g.Message); err != nil {
		log.Fatal(err)
	}
	fmt.Println("msg:", g.Name, g.Message)
	fmt.Println()
	return &DB{
		C: cluster,
		S: session,
	}
}
