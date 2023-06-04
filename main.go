package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tikv/client-go/v2/txnkv"
	"github.com/tikv/client-go/v2/txnkv/transaction"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("run failed %v\n", err)
	}
}

func createTable(txn *transaction.KVTxn, table string) error {
	err := txn.Set([]byte(tableKey(table)+"/a"), []byte("table beginning"))
	if err != nil {
		return fmt.Errorf("set failed, %w", err)
	}
	err = txn.Set([]byte(tableKey(table)+"/z"), []byte("table ending"))
	if err != nil {
		return fmt.Errorf("set failed, %w", err)
	}
	return nil
}
func tableKey(table string) string {
	return fmt.Sprintf("t/%s", table)
}

func createRecord(txn *transaction.KVTxn, table, id string, columns map[string]string) error {
	err := txn.Set([]byte(tableKey(table)+"/r/a"+id), []byte("row beginning"))
	if err != nil {
		return fmt.Errorf("set failed, %w", err)
	}
	err = txn.Set([]byte(tableKey(table)+"/r/z"+id), []byte("row ending"))
	if err != nil {
		return fmt.Errorf("set failed, %w", err)
	}
	for col, value := range columns {
		err = txn.Set([]byte(tableKey(table)+"/r/"+id+"/c/"+col), []byte(value))
		if err != nil {
			return fmt.Errorf("set failed, %w", err)
		}
	}
	return nil
}

func run(ctx context.Context) error {
	client, err := txnkv.NewClient([]string{"127.0.0.1:2379"})
	if err != nil {
		return fmt.Errorf("new client failed, %w", err)
	}

	txn, err := client.Begin()
	if err != nil {
		return fmt.Errorf("begin failed, %w", err)
	}

	err = createTable(txn, "people")
	if err != nil {
		return fmt.Errorf("set failed, %w", err)
	}
	log.Printf("commit ok")

	if err := printAllTableRecords(ctx, client, "people"); err != nil {
		return fmt.Errorf("print all table records failed, %w", err)
	}

	log.Printf("creating websites")
	txn, err = client.Begin()
	if err != nil {
		return fmt.Errorf("begin failed, %w", err)
	}
	if err := createTable(txn, "websites"); err != nil {
		return fmt.Errorf("create table failed, %w", err)
	}
	if err := createRecord(txn, "websites", "w1", map[string]string{
		"url": "google.com",
	}); err != nil {
		return fmt.Errorf("create record failed, %w", err)
	}
	if err := createRecord(txn, "websites", "ms", map[string]string{
		"url":   "microsoft.com",
		"title": "Microsoft",
	}); err != nil {
		return fmt.Errorf("create record failed, %w", err)
	}
	if err := txn.Commit(ctx); err != nil {
		return fmt.Errorf("commit failed, %w", err)
	}
	if err := printAllTableRecords(ctx, client, "websites"); err != nil {
		return fmt.Errorf("print all table records failed, %w", err)
	}
	return nil
}

func printAllTableRecords(ctx context.Context, client *txnkv.Client, table string) error {
	log.Printf("listing table %s records", table)
	txn, err := client.Begin()
	if err != nil {
		return fmt.Errorf("begin failed, %w", err)
	}
	iter, err := txn.Iter([]byte(tableKey(table)+"/a"), []byte(tableKey(table)+"/z"))
	if err != nil {
		return fmt.Errorf("iter failed, %w", err)
	}
	defer iter.Close()
	for iter.Valid() {
		k, v := iter.Key(), iter.Value()
		kk := append([]byte{}, k...)
		vv := append([]byte{}, v...)
		fmt.Printf("Entry: %s\n", fmt.Sprintf("k=%s, v=%s", kk, vv))
		if err := iter.Next(); err != nil {
			return fmt.Errorf("iter next failed, %w", err)
		}
	}
	return nil
}
