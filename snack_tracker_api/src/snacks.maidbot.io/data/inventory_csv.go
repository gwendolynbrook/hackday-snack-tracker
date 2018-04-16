package data

import (
    "os"
    "log"
    "encoding/csv"
    "path/filepath"
    "time"

    "app/src/snacks.maidbot.io/domain"
)

var FORMAT_STRING = "2006-01-02T15:04:05.999999-07:00"

func WriteChangeCsv(agg *domain.InventoryAggregate, writeHeaders bool, generateTime *time.Time) error {
  data_dir := os.Getenv(DATA_DIR_ENV)
  dump_file_name := "inventory_changes_" + generateTime.Format(FORMAT_STRING) + ".csv"
  dump_file, err :=  os.OpenFile(data_dir + "/" + dump_file_name,
    os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    log.Print("Cannot open change csv file.")
    return err
  }

  defer dump_file.Close()
  writer := csv.NewWriter(dump_file)
  defer writer.Flush()

  if writeHeaders {
    headers := agg.InventoryChanges[0].GetHeaders()
    if err = writer.Write(headers); err != nil {
      log.Print("Failed to write headers.")
      return err
    }
  }

  for _, change := range agg.InventoryChanges {
    values := change.ToSlice()
    if err = writer.Write(values); err != nil {
      log.Print("Failed to write values.")
      continue
    }
  }

  return nil
}

func WriteSummaryCsv(aggs []*domain.InventoryAggregate, generateTime *time.Time) error {
  data_dir := os.Getenv(DATA_DIR_ENV)
  dump_file_name := "inventory_summary_" + generateTime.Format(FORMAT_STRING) + ".csv"
  dump_file, err := os.OpenFile(data_dir + "/" + dump_file_name,
    os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    log.Print("Cannot open change csv file.")
    return err
  }

  defer dump_file.Close()
  writer := csv.NewWriter(dump_file)
  defer writer.Flush()


  headers := aggs[0].GetHeaders()
  writeHeaders := true
  if err = writer.Write(headers); err != nil {
    log.Print("Failed to write headers.")
    return err
  }

  for _, agg := range aggs {
    values := agg.ToSlice()
    if err = writer.Write(values); err != nil {
      log.Print("Failed to write values.")
      continue
    }

    WriteChangeCsv(agg, writeHeaders, generateTime)
    writeHeaders = false
  }

  return nil
}

func CleanupCsvs() {
  dataDir := os.Getenv(DATA_DIR_ENV)
  csvFiles, _ := filepath.Glob(dataDir + "/*.csv")
  for _, fileName := range csvFiles {
    os.Remove(fileName)
  }
}
