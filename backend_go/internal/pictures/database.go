package pictures

import (
	"errors"
	"fmt"
	"github.com/evanespen/vanespen.art_2025/configs"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/writer"
	"log"
	"os"
)

func Append(picture Picture) {
	if _, err := os.Stat(configs.PicturesDatabaseFile); err == nil {
		pictures, err := Read()
		if err != nil {
			log.Println(err, err)
		}
		pictures = append(pictures, picture)
		Write(pictures)
	} else {
		Write([]Picture{picture})
	}
}

func Write(pictures []Picture) {
	fw, err := local.NewLocalFileWriter(configs.PicturesDatabaseFile)
	if err != nil {
		log.Println("Can't create local file", err)
		return
	}

	pw, err := writer.NewParquetWriter(fw, new(Picture), 4)
	if err != nil {
		log.Println("Can't create parquet writer", err)
		return
	}

	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	for _, picture := range pictures {
		if err = pw.Write(picture); err != nil {
			log.Fatalf("error while writing parquet file: %v", err)
		}
	}

	if err = pw.WriteStop(); err != nil {
		log.Fatalf("error while closing parquet file: %v", err)
	}

	fmt.Println("write completed")
}

func Read() ([]Picture, error) {
	fr, err := local.NewLocalFileReader(configs.PicturesDatabaseFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open database file: %s", configs.PicturesDatabaseFile)
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, new(Picture), 4)
	if err != nil {
		return nil, errors.New("cannot create parquet reader")
	}
	defer pr.ReadStop()

	var pictures []Picture

	batchSize := 10
	num := int(pr.GetNumRows())

	if num < batchSize {
		batchSize = num
	}

	for {
		rows, err := pr.ReadByNumber(batchSize)
		if err != nil {
			fmt.Println(err)
			break // End of the file
		}

		for _, row := range rows {
			picture, ok := row.(Picture)
			if !ok {
				log.Printf("unable to convert: wanted Picture, got %T\n", row)
				continue
			}
			pictures = append(pictures, picture)
		}

		if len(rows) < batchSize {
			break // End of the file
		}
	}

	return pictures, nil
}
