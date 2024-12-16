package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/kwanpham2195/go-gcp-boilerplate/internal/derrors"
)

const testDBName = "discovery_postgres_test"

var testDB *DB

func TestMain(m *testing.M) {
	if err := CreateDBIfNotExists(testDBName); err != nil {
		if errors.Is(err, derrors.ErrNotFound) && os.Getenv("GO_DISCOVERY_TESTDB") != "true" {
			log.Printf("SKIPPING: could not connect to DB (see doc/postgres.md to set up): %v", err)
			return
		}
		log.Fatal(err)
	}

	var err error
	for _, driver := range []string{"postgres", "pgx"} {
		log.Printf("with driver %q", driver)
		testDB, err = Open(driver, DBConnURI(testDBName), "test")
		if err != nil {
			log.Fatalf("Open: %v %[1]T", err)
		}
		code := m.Run()
		if err := testDB.Close(); err != nil {
			log.Fatal(err)
		}
		if code != 0 {
			os.Exit(code)
		}
	}
}

func TestBulkInsert(t *testing.T) {
	table := "test_bulk_insert"

	for _, test := range []struct {
		name           string
		columns        []string
		values         []any
		conflictAction string
		wantErr        bool
		wantCount      int
		wantReturned   []string
	}{
		{
			name:      "test-one-row",
			columns:   []string{"colA"},
			values:    []any{"valueA"},
			wantCount: 1,
		},
		{
			name:      "test-multiple-rows",
			columns:   []string{"colA"},
			values:    []any{"valueA1", "valueA2", "valueA3"},
			wantCount: 3,
		},
		{
			name:    "test-invalid-column-name",
			columns: []string{"invalid_col"},
			values:  []any{"valueA"},
			wantErr: true,
		},
		{
			name:    "test-mismatch-num-cols-and-vals",
			columns: []string{"colA", "colB"},
			values:  []any{"valueA1", "valueB1", "valueA2"},
			wantErr: true,
		},
		{
			name:         "insert-returning",
			columns:      []string{"colA", "colB"},
			values:       []any{"valueA1", "valueB1", "valueA2", "valueB2"},
			wantCount:    2,
			wantReturned: []string{"valueA1", "valueA2"},
		},
		{
			name:    "test-conflict",
			columns: []string{"colA"},
			values:  []any{"valueA", "valueA"},
			wantErr: true,
		},
		{
			name:           "test-conflict-do-nothing",
			columns:        []string{"colA"},
			values:         []any{"valueA", "valueA"},
			conflictAction: OnConflictDoNothing,
			wantCount:      1,
		},
		{
			// This should execute the statement
			// INSERT INTO series (path) VALUES ('''); TRUNCATE series CASCADE;)');
			// which will insert a row with path value:
			// '); TRUNCATE series CASCADE;)
			// Rather than the statement
			// INSERT INTO series (path) VALUES (''); TRUNCATE series CASCADE;));
			// which would truncate most tables in the database.
			name:           "test-sql-injection",
			columns:        []string{"colA"},
			values:         []any{fmt.Sprintf("''); TRUNCATE %s CASCADE;))", table)},
			conflictAction: OnConflictDoNothing,
			wantCount:      1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			createQuery := fmt.Sprintf(`CREATE TABLE %s (
					colA TEXT NOT NULL,
					colB TEXT,
					PRIMARY KEY (colA)
				);`, table)
			if _, err := testDB.Exec(ctx, createQuery); err != nil {
				t.Fatal(err)
			}
			defer func() {
				dropTableQuery := fmt.Sprintf("DROP TABLE %s;", table)
				if _, err := testDB.Exec(ctx, dropTableQuery); err != nil {
					t.Fatal(err)
				}
			}()

			var err error
			var returned []string
			if test.wantReturned == nil {
				err = testDB.BulkInsert(ctx, table, test.columns, test.values, test.conflictAction)
			} else {
				err = testDB.BulkInsertReturning(ctx, table, test.columns, test.values, test.conflictAction,
					[]string{"colA"}, func(rows *sql.Rows) error {
						var r string
						if err := rows.Scan(&r); err != nil {
							return err
						}
						returned = append(returned, r)
						return nil
					})
			}
			if test.wantErr && err == nil || !test.wantErr && err != nil {
				t.Errorf("got error %v, wantErr %t", err, test.wantErr)
			}
			if err != nil {
				return
			}
			if test.wantCount != 0 {
				var count int
				query := "SELECT COUNT(*) FROM " + table
				row := testDB.QueryRow(ctx, query)
				err := row.Scan(&count)
				if err != nil {
					t.Fatalf("testDB.queryRow(%q): %v", query, err)
				}
				if count != test.wantCount {
					t.Errorf("testDB.queryRow(%q) = %d; want = %d", query, count, test.wantCount)
				}
			}
			if test.wantReturned != nil {
				sort.Strings(returned)
				if !cmp.Equal(returned, test.wantReturned) {
					t.Errorf("returned: got %v, want %v", returned, test.wantReturned)
				}
			}
		})
	}
}

func TestLargeBulkInsert(t *testing.T) {
	ctx := context.Background()
	if _, err := testDB.Exec(ctx, `CREATE TEMPORARY TABLE test_large_bulk (i BIGINT);`); err != nil {
		t.Fatal(err)
	}
	const size = 150001
	vals := make([]any, size)
	for i := 0; i < size; i++ {
		vals[i] = i + 1
	}
	start := time.Now()
	if err := testDB.Transact(ctx, sql.LevelDefault, func(db *DB) error {
		return db.BulkInsert(ctx, "test_large_bulk", []string{"i"}, vals, "")
	}); err != nil {
		t.Fatal(err)
	}
	t.Logf("large bulk insert took %s", time.Since(start))
	rows, err := testDB.Query(ctx, `SELECT i FROM test_large_bulk;`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	sum := int64(0)
	for rows.Next() {
		var i int64
		if err := rows.Scan(&i); err != nil {
			t.Fatal(err)
		}
		sum += i
	}
	var want int64 = size * (size + 1) / 2
	if sum != want {
		t.Errorf("sum = %d, want %d", sum, want)
	}
}

func TestBulkUpsert(t *testing.T) {
	ctx := context.Background()
	if _, err := testDB.Exec(ctx, `CREATE TEMPORARY TABLE test_replace (C1 int PRIMARY KEY, C2 int);`); err != nil {
		t.Fatal(err)
	}
	for _, values := range [][]any{
		{2, 4, 4, 8},                 // First, insert some rows.
		{1, -1, 2, -2, 3, -3, 4, -4}, // Then replace those rows while inserting others.
	} {
		err := testDB.Transact(ctx, sql.LevelDefault, func(tx *DB) error {
			return tx.BulkUpsert(ctx, "test_replace", []string{"C1", "C2"}, values, []string{"C1"})
		})
		if err != nil {
			t.Fatal(err)
		}
		var got []any
		err = testDB.RunQuery(ctx, `SELECT C1, C2 FROM test_replace ORDER BY C1`, func(rows *sql.Rows) error {
			var a, b int
			if err := rows.Scan(&a, &b); err != nil {
				return err
			}
			got = append(got, a, b)
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(got, values) {
			t.Errorf("%v: got %v, want %v", values, got, values)
		}
	}
}

func TestBuildUpsertConflictAction(t *testing.T) {
	got := buildUpsertConflictAction([]string{"a", "b"}, []string{"c", "d"})
	want := "ON CONFLICT (c, d) DO UPDATE SET a=excluded.a, b=excluded.b"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDBAfterTransactFails(t *testing.T) {
	ctx := context.Background()
	var tx *DB
	err := testDB.Transact(ctx, sql.LevelDefault, func(d *DB) error {
		tx = d
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	var i int
	err = tx.QueryRow(ctx, `SELECT 1`).Scan(&i)
	if err == nil {
		t.Fatal("got nil, want error")
	}
}

func TestBuildBulkUpdateQuery(t *testing.T) {
	q := buildBulkUpdateQuery("tab", []string{"K", "C1", "C2"}, []string{"TEXT", "INT", "BOOL"})
	got := strings.Join(strings.Fields(q), " ")
	w := `
		UPDATE tab
		SET C1 = data.C1, C2 = data.C2
		FROM (SELECT UNNEST($1::TEXT[]) AS K, UNNEST($2::INT[]) AS C1, UNNEST($3::BOOL[]) AS C2) AS data
		WHERE tab.K = data.K`
	want := strings.Join(strings.Fields(w), " ")
	if got != want {
		t.Errorf("\ngot\n%s\nwant\n%s", got, want)
	}
}

func TestBulkUpdate(t *testing.T) {
	ctx := context.Background()

	defer func(old int) { maxBulkUpdateArrayLen = old }(maxBulkUpdateArrayLen)
	maxBulkUpdateArrayLen = 5

	if _, err := testDB.Exec(ctx, `CREATE TABLE bulk_update (a INT, b INT)`); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := testDB.Exec(ctx, `DROP TABLE bulk_update`); err != nil {
			t.Fatal(err)
		}
	}()

	cols := []string{"a", "b"}
	var values []any
	for i := 0; i < 50; i++ {
		values = append(values, i, i)
	}
	err := testDB.Transact(ctx, sql.LevelDefault, func(tx *DB) error {
		return tx.BulkInsert(ctx, "bulk_update", cols, values, "")
	})
	if err != nil {
		t.Fatal(err)
	}

	// Update all even values of column a.
	updateVals := make([][]any, 2)
	for i := 0; i < len(values)/2; i += 2 {
		updateVals[0] = append(updateVals[0], i)
		updateVals[1] = append(updateVals[1], -i)
	}

	err = testDB.Transact(ctx, sql.LevelDefault, func(tx *DB) error {
		return tx.BulkUpdate(ctx, "bulk_update", cols, []string{"INT", "INT"}, updateVals)
	})
	if err != nil {
		t.Fatal(err)
	}

	err = testDB.RunQuery(ctx, `SELECT a, b FROM bulk_update`, func(rows *sql.Rows) error {
		var a, b int
		if err := rows.Scan(&a, &b); err != nil {
			return err
		}
		want := a
		if a%2 == 0 {
			want = -a
		}
		if b != want {
			t.Fatalf("a=%d: got %d, want %d", a, b, want)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestTransactSerializable(t *testing.T) {
	// Test that serializable transactions retry until success.
	// This test was taken from the example at https://www.postgresql.org/docs/11/transaction-iso.html,
	// section 13.2.3.
	ctx := context.Background()

	// Once in while, the test doesn't work. Repeat to de-flake.
	var msg string
	for i := 0; i < 20; i++ {
		msg = testTransactSerializable(ctx, t)
		if msg == "" {
			return
		}
	}
	t.Fatal(msg)
}

func testTransactSerializable(ctx context.Context, t *testing.T) string {
	const numTransactions = 4
	// A transaction that sums values in class 1 and inserts that sum into class 2,
	// or vice versa.
	insertSum := func(tx *DB, queryClass int) error {
		var sum int
		err := tx.QueryRow(ctx, `SELECT SUM(value) FROM ser WHERE class = $1`, queryClass).Scan(&sum)
		if err != nil {
			return err
		}
		insertClass := 3 - queryClass
		_, err = tx.Exec(ctx, `INSERT INTO ser (class, value) VALUES ($1, $2)`, insertClass, sum)
		return err
	}

	sawRetries := false
	for i := 0; i < 10; i++ {
		for _, stmt := range []string{
			`DROP TABLE IF EXISTS ser`,
			`CREATE TABLE ser (id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY, class INTEGER, value INTEGER)`,
			`INSERT INTO ser (class, value) VALUES (1, 10), (1, 20), (2, 100), (2, 200)`,
		} {
			if _, err := testDB.Exec(ctx, stmt); err != nil {
				t.Fatal(err)
			}
		}

		// Run the following two transactions multiple times concurrently:
		//   sum rows with class = 1 and insert as a row with class 2
		//   sum rows with class = 2 and insert as a row with class 1
		errc := make(chan error, numTransactions)
		for i := 0; i < numTransactions; i++ {
			i := i
			go func() {
				errc <- testDB.Transact(ctx, sql.LevelSerializable,
					func(tx *DB) error { return insertSum(tx, 1+i%2) })
			}()
		}
		// None of the transactions should fail.
		for i := 0; i < numTransactions; i++ {
			if err := <-errc; err != nil {
				return err.Error()
			}
		}
		t.Logf("max retries: %d", testDB.MaxRetries())
		// If nothing got retried, this test isn't exercising some important behavior.
		// Try again.
		if testDB.MaxRetries() > 0 {
			sawRetries = true
			break
		}
	}
	if !sawRetries {
		return "did not see any retries"
	}

	// Demonstrate serializability: there should be numTransactions new rows in
	// addition to the 4 we started with, and viewing the rows in insertion
	// order, each of the new rows should have the sum of the other class's rows
	// so far.
	type row struct {
		Class, Value int
	}
	rows, err := CollectStructs[row](ctx, testDB, `SELECT class, value FROM ser ORDER BY id`)
	if err != nil {
		return err.Error()
	}
	const initialRows = 4
	if got, want := len(rows), initialRows+numTransactions; got != want {
		return fmt.Sprintf("got %d rows, want %d", got, want)
	}
	sum := make([]int, 2)
	for i, r := range rows {
		if got, want := r.Value, sum[2-r.Class]; got != want && i >= initialRows {
			return fmt.Sprintf("row #%d: got %d, want %d", i, got, want)
		}
		sum[r.Class-1] += r.Value
	}
	return ""
}

func TestCopyDoesNotUpsert(t *testing.T) {
	// This test verifies that copying rows into a table will not overwrite existing rows.
	ctx := context.Background()
	conn, err := testDB.db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, stmt := range []string{
		`DROP TABLE IF EXISTS test_copy`,
		`CREATE TABLE test_copy (i  INTEGER PRIMARY KEY)`,
		`INSERT INTO test_copy (i) VALUES (1)`,
	} {
		if _, err := testDB.Exec(ctx, stmt); err != nil {
			t.Fatal(err)
		}
	}

	err = conn.Raw(func(c any) error {
		stdConn, ok := c.(*stdlib.Conn)
		if !ok {
			t.Skip("DB driver is not pgx")
		}
		rows := [][]any{{1}, {2}}
		_, err = stdConn.Conn().CopyFrom(ctx, []string{"test_copy"}, []string{"i"}, pgx.CopyFromRows(rows))
		return err
	})

	const constraintViolationCode = "23505"
	var gerr *pgconn.PgError
	if !errors.As(err, &gerr) || gerr.Code != constraintViolationCode {
		t.Errorf("got %v, wanted code %s", gerr, constraintViolationCode)
	}
}

func TestRunQueryIncrementally(t *testing.T) {
	ctx := context.Background()
	for _, stmt := range []string{
		`DROP TABLE IF EXISTS test_rqi`,
		`CREATE TABLE test_rqi (i  INTEGER PRIMARY KEY)`,
		`INSERT INTO test_rqi (i) VALUES (1), (2), (3), (4), (5)`,
	} {
		if _, err := testDB.Exec(ctx, stmt); err != nil {
			t.Fatal(err)
		}
	}
	query := `SELECT i FROM test_rqi ORDER BY i LIMIT $1`
	var got []int

	// Run until all rows consumed.
	err := testDB.RunQueryIncrementally(ctx, query, 2, func(rows *sql.Rows) error {
		var i int
		if err := rows.Scan(&i); err != nil {
			return err
		}
		got = append(got, i)
		return nil
	}, 4)
	if err != nil {
		t.Fatal(err)
	}
	want := []int{1, 2, 3, 4}
	if !cmp.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// Stop early.
	got = nil
	err = testDB.RunQueryIncrementally(ctx, query, 2, func(rows *sql.Rows) error {
		var i int
		if err := rows.Scan(&i); err != nil {
			return err
		}
		got = append(got, i)
		if len(got) == 3 {
			return io.EOF
		}
		return nil
	}, 10)
	if err != nil {
		t.Fatal(err)
	}
	want = []int{1, 2, 3}
	if !cmp.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCollectStrings(t *testing.T) {
	ctx := context.Background()
	for _, stmt := range []string{
		`DROP TABLE IF EXISTS test_cs`,
		`CREATE TABLE test_cs (s TEXT)`,
		`INSERT INTO test_cs (s) VALUES ('a'), ('b'), ('c')`,
	} {
		if _, err := testDB.Exec(ctx, stmt); err != nil {
			t.Fatal(err)
		}
	}
	got, err := Collect1[string](ctx, testDB, `SELECT s FROM test_cs`)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	want := []string{"a", "b", "c"}
	if !cmp.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
