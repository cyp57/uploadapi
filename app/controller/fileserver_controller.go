package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"time"

	"github.com/cyp57/uploadapi/pkg/mongodb"
	"github.com/cyp57/uploadapi/setting"
	"github.com/cyp57/uploadapi/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//upload file to fileserver or local

type FileServer struct {
	Id         string    `json:"id" bson:"id"`
	FileName   string    `json:"fileName" bson:"fileName"`
	FileType   string    `json:"fileType" bson:"fileType"`
	FileSize   int64     `json:"fileSize" bson:"fileSize"`
	Created_at time.Time `json:"created_at" bson:"created_at"`
}

func (f *FileServer) SaveFile(file *multipart.File) (interface{}, error) {
	var result = bson.M{}

	now, err := utils.GetCurrentTime()
	if err != nil {
		return nil, err
	}
	f.Created_at = now
	f.Id = utils.GenerateOid("FS")

	root := utils.GetYaml("RootFile")
	filepath := path.Join(root, f.Id+"."+f.FileType)

	dst, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}

	defer dst.Close()

	_, err = io.Copy(dst, *file)
	if err != nil {
		return nil, err
	}
	// save fileinfo to database
	_, err = f.saveFileInfo()
	if err != nil {
		return nil, err
	}

	url := path.Join(utils.GetYaml("BasePath")+":"+utils.GetYaml("HTTPPort"), utils.GetYaml("ServiceName"), "file/server", "url", f.Id+"."+f.FileType)

	result = bson.M{"id": f.Id, "url": url}

	return result, nil
}

func (f *FileServer) saveFileInfo() (interface{}, error) {
	collection := setting.CollectionSetting.FileserverCollection
	db := mongodb.Database

	_, err := db.Collection(collection).InsertOne(context.TODO(), f)
	if err != nil {
		return nil, err
	}

	return f.Id, nil
}

func (f *FileServer) GetFileDetail(id string) (interface{}, error) {

	collection := setting.CollectionSetting.FileserverCollection
	db := mongodb.Database

	result := make(primitive.M)
	opts := options.FindOne().SetProjection(bson.M{"_id": 0})
	filter := bson.M{"id": id}
	err := db.Collection(collection).FindOne(context.TODO(), filter, opts).Decode(&result)
	if err != nil {
		return nil, err
	}

	url := path.Join(utils.GetYaml("BasePath")+":"+utils.GetYaml("HTTPPort"), utils.GetYaml("ServiceName"), "file/server", "url", fmt.Sprint(result["id"])+"."+fmt.Sprint(result["fileType"]))
	result["url"] = url

	return result, nil
}

func (f *FileServer) DeleteFile(id string) (interface{}, error) {

	detail, err := f.GetFileDetail(id)
	if err != nil {
		return nil, err
	}
	filedata := detail.(primitive.M)
	filename := id + "." + fmt.Sprint(filedata["fileType"])

	root := utils.GetYaml("RootFile")

	pathDir := path.Join(root, filename)

	err = os.Remove(pathDir)
	if err != nil {
		return nil, err
	}

	// delete info on database
	collection := setting.CollectionSetting.FileserverCollection
	db := mongodb.Database
	filter := bson.M{"id": id}

	deleteResult, err := db.Collection(collection).DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if deleteResult.DeletedCount != 1 {
		return false, errors.New("DeletedCount : " + fmt.Sprint(deleteResult.DeletedCount))
	}

	result := bson.M{
		"id": id,
	}

	return result, nil
}
