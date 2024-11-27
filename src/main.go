package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	pgxDriver   string = "pgx"
	mysqlDriver string = "mysql"
)

// func main() {
// 	// Opening a driver typically will not attempt to connect to the database.
// 	uri := "postgres://admin:admin@db/sample?sslmode=disable"
// 	db, err := sql.Open("pgx", uri)
// 	if err != nil {
// 		// This will not be a connection error, but a DSN parse error or
// 		// another initialization error.
// 		log.Fatal(err)
// 	}
// 	db.SetConnMaxLifetime(0)
// 	db.SetMaxIdleConns(50)
// 	db.SetMaxOpenConns(50)

// 	s := &Service{db: map[string]*sql.DB{"postgresql": db}}

// 	http.ListenAndServe(":8080", s)
// }

func main() {
	dbs := make(map[string]*sql.DB)

	// pgx dirver(postgresql)
	{
		uri := "postgres://admin:admin@pg/sample?sslmode=disable"

		db, err := sql.Open(pgxDriver, uri)
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}

		db.SetConnMaxLifetime(0)
		db.SetMaxIdleConns(50)
		db.SetMaxOpenConns(50)

		dbs[pgxDriver] = db
	}

	// mysql driver for mariadb
	{
		uri := "admin:admin@tcp(maria:3306)/sample?tls=false"

		db, err := sql.Open(mysqlDriver, uri)
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}

		db.SetConnMaxLifetime(0)
		db.SetMaxIdleConns(50)
		db.SetMaxOpenConns(50)

		dbs[mysqlDriver] = db
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /variables/ErrConnDone", func(w http.ResponseWriter, r *http.Request) {
		variables(w, r, dbs, sql.ErrConnDone)
	})

	mux.HandleFunc("GET /variables/ErrNoRows", func(w http.ResponseWriter, r *http.Request) {
		variables(w, r, dbs, sql.ErrNoRows)
	})

	mux.HandleFunc("GET /variables/ErrTxDone", func(w http.ResponseWriter, r *http.Request) {
		variables(w, r, dbs, sql.ErrTxDone)
	})

	mux.HandleFunc("GET /drivers", func(w http.ResponseWriter, r *http.Request) {
		// 登録されてるdriverを出力するだけ
		io.WriteString(w, fmt.Sprintf("%v", sql.Drivers()))
	})

	mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		// driverを名前つけてregisterするだけ、用途はよくわからない
		// 名前衝突するとpanic
		defer func() {
			if r := recover(); r != nil {
				io.WriteString(w, fmt.Sprintf("%v, %s", sql.Drivers(), "Don't call it twice."))
			}
		}()

		// 関係ないけどゴキDBはpg向けのdriverでイケるらしい
		sql.Register("cockroach", dbs[pgxDriver].Driver())

		io.WriteString(w, fmt.Sprintf("%v", sql.Drivers()))
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		health(w, r, dbs)
	})
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		fetchUser(w, r, dbs)
	})

	s := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.ListenAndServe()
}

func variables(w http.ResponseWriter, r *http.Request, dbs map[string]*sql.DB, targetErr error) {
	db, ok := dbs[pgxDriver]
	if !ok {
		http.Error(w, "sql.Openまわりでしくじってるかも", http.StatusInternalServerError)
	}

	conn, err := db.Conn(r.Context())
	if err != nil {
		http.Error(w, "connが張れない", http.StatusInternalServerError)
	}
	defer conn.Close()

	if errors.Is(targetErr, sql.ErrConnDone) {
		// Queryを発行した際に利用したconnがすでにCloseされているケースなどに発生する
		// 利用しようとしたconnがcloseされている原因は色々考えられるため、このケースほど単純でなさそうなのでご注意を。
		conn.Close()
		if err := conn.PingContext(r.Context()); err != nil && errors.Is(err, sql.ErrConnDone) {
			fmt.Fprintf(w, "ErrConnDoneが発生した。")
		}
		return
	} else if errors.Is(targetErr, sql.ErrNoRows) {
		// database/sqlを利用した場合、queryしたデータが見つからない場合このErrが発生する
		// gormでは`ErrRecordNotFound`になっているのでwrapしているのかと考えて、調べてみたが不明。使っていない？
		// gormのソースコードで`ErrNoRows`をgrep => ない
		// mysql driverで`ErrNoRows`をgrep => testcodeだった
		//     => https://github.com/search?q=repo%3Ago-sql-driver%2Fmysql%20ErrNoRows&type=code
		//
		// 一方、entは`ErrNoRows`でgrepするとヒットした。
		//     => https://github.com/ent/ent/blob/365b49817603b056407b2c992fd92ce1790c8723/dialect/sql/scan.go#L30
		// database/sqlを利用せず、低レイヤ操作を自前実装しているからっぽい？
		//
		// 総じて、database/sqlをそのまんま実装するケースでしか見かけないかも？
		var username string
		err2 := conn.QueryRowContext(r.Context(), "SELECT * FROM users WHERE username = 'buriburi';").Scan(&username)
		if err2 != nil && errors.Is(err2, sql.ErrNoRows) {
			fmt.Fprintf(w, "ErrConnDoneが発生した。")
		}
		return
	} else if errors.Is(targetErr, sql.ErrTxDone) {
		fmt.Println("match! TxDone")

	} else {
		http.Error(w, "なんかミスってる", http.StatusInternalServerError)
		return
	}
}

func health(w http.ResponseWriter, r *http.Request, dbs map[string]*sql.DB) {
	// 疎通してるか確認用
	for k, v := range maps.All(dbs) {
		if err := v.PingContext(r.Context()); err != nil {
			fmt.Fprintf(w, "driver: %s is not working", k)
			r.Response.StatusCode = http.StatusInternalServerError
			return
		}
	}
	fmt.Fprintf(w, "All connections are healthy!")
}

func fetchUser(w http.ResponseWriter, r *http.Request, dbs map[string]*sql.DB) {
	driver := "pgx"
	// db, ok := dbs[r.PathValue("driver")]
	db, ok := dbs[driver]
	fmt.Println(r.PathValue("driver"))
	if !ok {
		http.Error(w, "no `dirver` provided", http.StatusBadRequest)

		return
	}
	ID := r.PathValue("id")
	fmt.Printf("ID: %s\n", ID)
	if !ok {
		http.Error(w, "no `id` provided", http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// TODO: sanitize
	query := fmt.Sprintf("SELECT username FROM users WHERE id = %s;", ID)
	fmt.Println(query)
	var firstName string
	var unko string
	err := db.QueryRowContext(ctx, query).Scan(&firstName)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, fmt.Sprintf("username: %s, Unko: %s", firstName, unko))

}

type Service struct {
	db map[string]*sql.DB
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	action := paths[1]
	dbname := paths[2]

	db := s.db[dbname]
	switch action {
	default:
		http.Error(w, "not found", http.StatusNotFound)
		return
	case "healthz":
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		err := db.PingContext(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("db down: %v", err), http.StatusFailedDependency)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case "quick-action":
		// This is a short SELECT. Use the request context as the base of
		// the context timeout.
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		query := "SELECT username FROM users WHERE id = 1;"
		var name string
		err := db.QueryRowContext(ctx, query).Scan(&name)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.WriteString(w, fmt.Sprintf("username: %s", name))
		return
	case "long-action":
		// This is a long SELECT. Use the request context as the base of
		// the context timeout, but give it some time to finish. If
		// the client cancels before the query is done the query will also
		// be canceled.
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		var names []string
		rows, err := db.QueryContext(ctx, "select p.name from people as p where p.active = true;")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var name string
			err = rows.Scan(&name)
			if err != nil {
				break
			}
			names = append(names, name)
		}
		// Check for errors during rows "Close".
		// This may be more important if multiple statements are executed
		// in a single batch and rows were written as well as read.
		if closeErr := rows.Close(); closeErr != nil {
			http.Error(w, closeErr.Error(), http.StatusInternalServerError)
			return
		}

		// Check for row scan error.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check for errors during row iteration.
		if err = rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(names)
		return
	case "async-action":
		// This action has side effects that we want to preserve
		// even if the client cancels the HTTP request part way through.
		// For this we do not use the http request context as a base for
		// the timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var orderRef = "ABC123"
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = tx.ExecContext(ctx, "stored_proc_name", orderRef)

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tx.Commit()
		if err != nil {
			http.Error(w, "action in unknown state, check state before attempting again", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}
