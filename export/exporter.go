package export

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"os"
	"os/exec"
	"time"
)

const (
	ExporterVersion = "0.1"
)

type Exporter struct {
	conn driver.Conn
	conf Config
	f    *os.File
}

func New(conf Config) (*Exporter, error) {
	exporter := &Exporter{
		conf: conf,
	}

	conn, err := exporter.connect(context.Background())
	if err != nil {
		return nil, err
	}
	exporter.conn = conn
	return exporter, nil
}

func (e *Exporter) connect(ctx context.Context) (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", e.conf.Host, e.conf.Port)},
		Auth: clickhouse.Auth{
			Database: e.conf.Database,
			Username: e.conf.Username,
			Password: e.conf.Password,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "clickhouse-export", Version: ExporterVersion},
			},
		},

		Debugf: func(format string, v ...interface{}) {
			fmt.Printf(format, v)
		},
		//TLS: &tls.Config{
		//	InsecureSkipVerify: true,
		//},
	})
	if err != nil {
		fmt.Println("1", err)
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		fmt.Println("2", err)
		return nil, err
	}
	return conn, nil
}

func (e *Exporter) BatchExport(ctx context.Context) error {
	f, err := os.OpenFile(e.conf.OutputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	e.f = f
	defer f.Close()

	if e.conf.BatchSize == 0 {
		return e.export(ctx)
	}

	limit := e.conf.BatchSize
	offset := 0
	format := e.conf.Format
	for {
		fmt.Printf("%v :: Exporting ... batch size: %v, offset: %v\n", time.Now(), limit, offset)
		query := setLimitOnQuery(e.conf.Query, limit, offset)
		data, err := e.runQuery(ctx, query, format)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			break
		}
		if err := e.write(data); err != nil {
			return err
		}
		format = "CSV"
		offset += limit
	}

	return nil
}

func (e *Exporter) export(ctx context.Context) error {
	data, err := e.runQuery(ctx, e.conf.Query, e.conf.Format)
	if err != nil {
		return err
	}
	if err := e.write(data); err != nil {
		return err
	}
	return nil
}

func (e *Exporter) runQuery(ctx context.Context, query, format string) ([]byte, error) {
	c := newCommand()
	c.appendParam("host", e.conf.Host)
	c.appendParam("port", e.conf.Port)
	if e.conf.Username != "" {
		c.appendParam("user", e.conf.Username)
	}
	if e.conf.Password != "" {
		c.appendParam("password", e.conf.Password)
	}
	c.appendParam("format", format)
	c.appendParam("query", query)

	//fmt.Println(c)

	out, err := exec.CommandContext(ctx, c.getBase(), c.getParams()...).Output()
	//fmt.Println(string(out))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (e *Exporter) write(data []byte) error {
	if _, err := e.f.Write(data); err != nil {
		return err
	}
	return nil
}

func setLimitOnQuery(query string, limit, offset int) string {
	return fmt.Sprintf("%s limit %v offset %v", query, limit, offset)
}
