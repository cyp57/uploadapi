package controller

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path"
	"time"

	"github.com/cyp57/uploadapi/pkg/mongodb"
	"github.com/cyp57/uploadapi/setting"
	"github.com/cyp57/uploadapi/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//upload file with gridfs

type GridFs struct {
	Id       string
	FileName string
	FileType string
	FileSize int64
	Data     []byte
}


func (g *GridFs) UploadFile() (interface{}, error) {

	db := mongodb.Database

	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName("myfile"))
	if err != nil {
		return nil, err
	}

	id := utils.GenerateOid("FS")
	now, err := utils.GetCurrentTime()
	if err != nil {
		return nil, err
	}

	uploadSream, err := bucket.OpenUploadStream(id + "." + g.FileType, options.GridFSUpload().SetMetadata(
		bson.M{"id": id, "fileName": g.FileName, "fileType": g.FileType, "fileSize": g.FileSize, "created_at": now}))

	if err != nil {
		return nil, err
	}

	defer uploadSream.Close()

	fileSize, err := uploadSream.Write(g.Data)
	if err != nil {
		return nil, err
	}

	defer bucket.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer bucket.SetWriteDeadline(time.Now().Add(10 * time.Second))
	log.Printf("Write file to DB was successful. File size: %d\n", fileSize)

	url := path.Join(utils.GetYaml("BasePath")+":"+utils.GetYaml("HTTPPort"), utils.GetYaml("ServiceName"), "gridfs/file", id +"." +g.FileType)
	
	result := bson.M{"id": id , "url" : url}


	return result, nil
}

func (g *GridFs) GetFileDetail(id string) (interface{}, error) {

	db := mongodb.Database
	collectionName := setting.CollectionSetting.GridFsCollection

	var resultDb = primitive.M{}

	err := db.Collection(collectionName).FindOne(context.TODO(), bson.M{"metadata.id": id}).Decode(&resultDb)
	if err != nil {
	
		return nil, err
	}

	result := resultDb["metadata"].(primitive.M)


	url := path.Join(utils.GetYaml("BasePath")+":"+utils.GetYaml("HTTPPort"), utils.GetYaml("ServiceName"), "gridfs/file", fmt.Sprint(resultDb["filename"]))
	result["url"] = url

	return result, nil
}

func (g *GridFs) DownloadFile(fileName string) (int64, []byte, error) {

	db := mongodb.Database
	collectionName := setting.CollectionSetting.GridFsCollection
	fileCollection := db.Collection(collectionName)

	var resultDb = primitive.M{}
	err := fileCollection.FindOne(context.TODO(), bson.M{"filename": fileName}).Decode(&resultDb)
	if err != nil {
		// not found
		return 0, nil, err
	}

	
	bucket, err := gridfs.NewBucket(
		db,options.GridFSBucket().SetName("myfile"),
	)
	if err != nil {
		return 0, nil, err
	}

	var buf bytes.Buffer

	dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
	if err != nil {
		log.Println(err.Error())
		return 0, nil, err
	}
	fmt.Printf("File size to download: %v\n", dStream)
	content := buf.Bytes()

	defer bucket.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer bucket.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return dStream, content, nil
}

func (g *GridFs)DeleteFile(id string) (interface{} , error) {
  
	db := mongodb.Database
	collectionName := setting.CollectionSetting.GridFsCollection
	fileCollection := db.Collection(collectionName)
   
	var resultDb = primitive.M{}
    err := fileCollection.FindOne(context.TODO(), bson.M{"metadata.id": id}).Decode(&resultDb)
    if err != nil {
        return nil ,err
    }
    bucket, _ := gridfs.NewBucket(
        db,options.GridFSBucket().SetName("myfile"),
    )

    if err := bucket.Delete(resultDb["_id"]); err != nil && err != gridfs.ErrFileNotFound {
        return nil ,err
    }
    defer bucket.SetReadDeadline(time.Now().Add(10 * time.Second))
    defer bucket.SetWriteDeadline(time.Now().Add(10 * time.Second))
    return id , err
}
