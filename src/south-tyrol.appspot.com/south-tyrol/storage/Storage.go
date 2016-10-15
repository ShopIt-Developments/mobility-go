package storage

import (
    "google.golang.org/appengine"
    "cloud.google.com/go/storage"
    "google.golang.org/appengine/file"
    "io/ioutil"
    "net/http"
    "strings"
)

func CreateFile(r *http.Request, fileName string) {
    ctx := appengine.NewContext(r)
    client, _ := storage.NewClient(ctx)
    bucket, _ := file.DefaultBucketName(ctx)

    wc := client.Bucket(bucket).Object(fileName).NewWriter(ctx)
    wc.ContentType = "text/plain"

    wc.Metadata = map[string]string{
        "x-goog-meta-foo": "foo",
        "x-goog-meta-bar": "bar",
    }

    if _, err := wc.Write([]byte("abcde\n")); err != nil {
        //d.errorf("createFile: unable to write data to bucket %q, file %q: %v", bucket, fileName, err)
        return
    }

    if _, err := wc.Write([]byte(strings.Repeat("f", 1024 * 4) + "\n")); err != nil {
        //d.errorf("createFile: unable to write data to bucket %q, file %q: %v", bucket, fileName, err)
        return
    }

    if err := wc.Close(); err != nil {
        //d.errorf("createFile: unable to close bucket %q, file %q: %v", bucket, fileName, err)
        return
    }
}

func WriteFile(r *http.Request, fileName string) ([]byte, error) {
    ctx := appengine.NewContext(r)
    client, _ := storage.NewClient(ctx)
    bucket, _ := file.DefaultBucketName(ctx)

    rc, err := client.Bucket(bucket).Object(fileName).NewReader(ctx)

    if err != nil {
        return nil, err
    }

    defer rc.Close()

    contents, err := ioutil.ReadAll(rc)

    if err != nil {
        return nil, err
    }

    return contents, nil
}

func ReadFile(r *http.Request, fileName string) ([]byte, error) {
    ctx := appengine.NewContext(r)
    client, _ := storage.NewClient(ctx)
    bucket, _ := file.DefaultBucketName(ctx)

    rc, err := client.Bucket(bucket).Object(fileName).NewReader(ctx)

    if err != nil {
        return nil, err
    }

    defer rc.Close()

    contents, err := ioutil.ReadAll(rc)

    if err != nil {
        return nil, err
    }

    return contents, nil
}