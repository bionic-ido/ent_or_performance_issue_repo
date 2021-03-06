package bug

import (
	"context"
	"entgo.io/bug/ent/admin"
	"entgo.io/bug/ent/user"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"entgo.io/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"entgo.io/bug/ent"
	"entgo.io/bug/ent/enttest"
)

func TestBugSQLite(t *testing.T) {
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	test(t, client)
}

func TestBugMySQL(t *testing.T) {
	for version, port := range map[string]int{"56": 3306, "57": 3307, "8": 3308} {
		addr := net.JoinHostPort("localhost", strconv.Itoa(port))
		t.Run(version, func(t *testing.T) {
			client := enttest.Open(t, dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
			defer client.Close()
			test(t, client)
		})
	}
}

func TestBugPostgres(t *testing.T) {
	for version, port := range map[string]int{"10": 5430, "11": 5431, "12": 5432, "13": 5433, "14": 5434} {
		t.Run(version, func(t *testing.T) {
			client := enttest.Open(t, dialect.Postgres, fmt.Sprintf("host=localhost port=%d user=postgres dbname=test password=pass sslmode=disable", port))
			defer client.Close()
			test(t, client)
		})
	}
}

func TestBugMaria(t *testing.T) {
	for version, port := range map[string]int{"10.5": 4306, "10.2": 4307, "10.3": 4308} {
		t.Run(version, func(t *testing.T) {
			addr := net.JoinHostPort("localhost", strconv.Itoa(port))
			client := enttest.Open(t, dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
			defer client.Close()
			test(t, client)
		})
	}
}

func fillDatabase(client *ent.Client, ctx context.Context) []*ent.Admin {
	// Just fill the database with a lot of data - 100 admins, each with 1 team lead & 1000 members.
	// This will help demonstrate the difference in query times, as one will do so with index scans, while the other will be a sequential scan on this huge database
	i := 0
	bulkAdmins := make([]*ent.AdminCreate, 100)
	for i < 100 {
		bulkAdmins[i] = client.Admin.Create().SetName("a").SetAge(35)
		i++
	}

	dbAdmins := client.Admin.CreateBulk(bulkAdmins...).SaveX(ctx)

	var bulkUsers []*ent.UserCreate
	for _, dbAdmin := range dbAdmins {
		bulkUsers = append(bulkUsers, client.User.Create().SetName("a").SetAge(20).SetLeadAdminID(dbAdmin.ID))

		j := 0
		for j < 1000 {
			bulkUsers = append(bulkUsers, client.User.Create().SetName("b").SetAge(18).SetMemberAdminID(dbAdmin.ID))
			j++
		}
	}

	// Create in bulk since there are too many to do at once
	for index := 0; index < len(bulkUsers); index += 5000 {
		end := index + 5000
		if end > len(bulkUsers) {
			end = len(bulkUsers)
		}

		client.User.CreateBulk(bulkUsers[index:end]...).SaveX(ctx)
	}

	return dbAdmins
}

// Note I only checked this with TestPostgres here since that is what interests me, idk if it is an issue on the others or not.
// This tests query time of the Or predicate vs query time of doing the same thing as 2 separate queries (applicative union. Union would also have worked to improve the time from manual testing in raw sql in front of postgres). This shows the problem with the Or predicate performance, as we see that 2 separate queries take less than a second, where the Or predicate query takes a long time that only increases the more we insert into the database (as it does a sequential scan)
func test(t *testing.T, client *ent.Client) {
	ctx := context.Background()

	fmt.Println("Clean database")
	client.User.Delete().ExecX(ctx)
	client.Admin.Delete().ExecX(ctx)

	// Clean database at the end as well
	defer client.User.Delete().ExecX(ctx)
	defer client.Admin.Delete().ExecX(ctx)

	// Do fillDatabase a couple of times to really fill it up, and then only get 100 admins to query on from one of the inserts
	fmt.Println("Fill up database")
	i := 0
	for i < 9 {
		i++
		fillDatabase(client, ctx)
	}

	dbAdmins := fillDatabase(client, ctx)

	// Test query that doesn't have performance issue
	fmt.Println("Check timing of non-performance issue query (2 separate queries for member admin & lead admin)")
	start := time.Now().UnixMilli()

	for _, dbAdmin := range dbAdmins {
		client.User.Query().Where(user.HasLeadAdminWith(admin.ID(dbAdmin.ID))).AllX(ctx)
		client.User.Query().Where(user.HasMemberAdminWith(admin.ID(dbAdmin.ID))).AllX(ctx)
	}

	end := time.Now().UnixMilli()
	timeSeconds := float64(end-start) / 1000
	//require.True(t, timeSeconds < 1, "Expected query time to be less then 1 seconds, got %f", timeSeconds)
	fmt.Printf("Took %f\n", timeSeconds)

	// Test query that has performance issue
	fmt.Println("Check timing of performance issue query (1 query using Or predicate for member admin & lead admin)")
	start = time.Now().UnixMilli()

	for _, dbAdmin := range dbAdmins {
		client.User.Query().Where(user.Or(
			user.LeadAdminID(dbAdmin.ID),
			user.MemberAdminID(dbAdmin.ID),
		)).AllX(ctx)
	}

	end = time.Now().UnixMilli()
	timeSeconds = float64(end-start) / 1000
	//require.True(t, timeSeconds < 1, "Expected query time to be less then 1 seconds, got %f", timeSeconds)
	fmt.Printf("Took %f\n", timeSeconds)
}
