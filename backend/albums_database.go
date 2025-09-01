package main

import (
	"errors"
	"fmt"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/writer"
	"log"
	"os"
)

func AppendAlbum(album Album) {
	if _, err := os.Stat(AlbumsDatabaseFile); err == nil {
		albums, err := ReadAlbums()
		if err != nil {
			log.Println(err, err)
		}
		albums = append(albums, album)
		WriteAlbums(albums)
	} else {
		WriteAlbums([]Album{album})
	}
}

func WriteAlbums(albums []Album) {
	fw, err := local.NewLocalFileWriter(AlbumsDatabaseFile)
	if err != nil {
		log.Println("Can't create local file", err)
		return
	}

	pw, err := writer.NewParquetWriter(fw, new(Album), 4)
	if err != nil {
		log.Println("Can't create parquet writer", err)
		return
	}

	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	for _, album := range albums {
		if err = pw.Write(album); err != nil {
			log.Fatalf("error while writing parquet file: %v", err)
		}
	}

	if err = pw.WriteStop(); err != nil {
		log.Fatalf("error while closing parquet file: %v", err)
	}

	fmt.Println("write completed")
}

func ReadAlbums() ([]Album, error) {
	fr, err := local.NewLocalFileReader(AlbumsDatabaseFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open database file: %s", AlbumsDatabaseFile)
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, new(Album), 4)
	if err != nil {
		return nil, errors.New("cannot create parquet reader")
	}
	defer pr.ReadStop()

	var albums []Album

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
			album, ok := row.(Album)
			if !ok {
				log.Printf("unable to convert: wanted Album, got %T\n", row)
				continue
			}
			albums = append(albums, album)
		}

		if len(rows) < batchSize {
			break // End of the file
		}
	}

	return albums, nil
}
