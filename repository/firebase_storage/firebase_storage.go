package repository_firebase_storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"palco-planner-api/api/repository"
	"palco-planner-api/entity"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func initFirebase() *firebase.App {
	ctx := context.Background()
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório atual:", err)
		panic(err)
	}
	path := filepath.Join(dir, "serviceAccountKey.json")
	sa := option.WithCredentialsFile(path)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}
	return app
}

func getConfigJsonFirebase() (sa entity.ServiceAccount) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório atual:", err)
	}
	path := filepath.Join(dir, "serviceAccountKey.json")

	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to read service account file: %v", err)
	}

	if err := json.Unmarshal(bytes, &sa); err != nil {
		log.Fatalf("Unable to parse service account JSON: %v", err)
	}

	return sa
}

type FirebaseStorage struct {
	clientFirebase *firebase.App
	bucketFirebase string
	sa             entity.ServiceAccount
}

func (repositorty FirebaseStorage) ListFiles(folderName string) (listFiles []string, err error) {
	app := initFirebase()

	ctx := context.TODO()

	client, err := app.Storage(ctx)
	if err != nil {
		return
	}

	bucket, err := client.Bucket(repositorty.bucketFirebase)
	if err != nil {
		return
	}

	it := bucket.Objects(ctx, &storage.Query{Prefix: ""})
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			break
		}
		nameFile, _ := repositorty.generateSignedURL(objAttrs.Name)
		listFiles = append(listFiles, nameFile)
	}

	return
}

func (repositorty FirebaseStorage) UploadFile(file *multipart.FileHeader) (err error) {
	ctx := context.Background()
	client, err := repositorty.clientFirebase.Storage(ctx)
	if err != nil {
		return
	}

	bucket, err := client.Bucket(repositorty.bucketFirebase)
	if err != nil {
		return
	}

	namefile := fmt.Sprintf("%s.pdf", uuid.NewString())

	sw := bucket.Object(namefile).NewWriter(ctx)
	defer sw.Close()

	f, err := file.Open()
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(sw, f)
	return
}

func (repositorty FirebaseStorage) generateSignedURL(objectPath string) (string, error) {
	client, err := repositorty.clientFirebase.Storage(context.TODO())
	if err != nil {
		return "", fmt.Errorf("error creating storage client: %v", err)
	}

	bucket, err := client.Bucket(repositorty.bucketFirebase)
	if err != nil {
		return "", fmt.Errorf("error getting default bucket: %v", err)
	}

	signedURL, err := bucket.SignedURL(objectPath, &storage.SignedURLOptions{
		GoogleAccessID: repositorty.sa.ClientEmail,
		PrivateKey:     []byte(repositorty.sa.PrivateKey),
		Method:         "GET",
		Expires:        time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create signed URL: %v", err)
	}

	return signedURL, nil
}

func NewFirebaseStorage() repository.Storage {
	fibaseCLient := initFirebase()
	sa := getConfigJsonFirebase()
	return FirebaseStorage{
		fibaseCLient,
		"palco-planner-test.appspot.com",
		sa,
	}
}
