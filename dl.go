package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func DownloadFile(fileloc, url, mime string) (int64, error) {
	head, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	if head.StatusCode != 200 {
		return 0, errors.New(fmt.Sprintf("Bad response, expected 200 got %d (%s)", head.StatusCode, url))
	}
	if mime != "*" && head.Header.Get("Content-Type") != mime {
		return 0, errors.New(fmt.Sprintf("Wrong MIME, expected %s got %s (%s)", mime, head.Header.Get("Content-Type"), url))
	}
	if !fileExists(fileloc) {
		return writeUrlToFile(fileloc, url, mime)
	} else {
		contlength, err := strconv.ParseInt(head.Header.Get("Content-Length"), 10, 0)
		if err != nil {
			return 0, err
		}
		file, err := os.Open(fileloc)
		if err != nil {
			return 0, err
		}
		defer file.Close()

		fstat, err := file.Stat()
		if err != nil {
			return 0, err
		}

		if fstat.Size() == contlength {
			fmt.Printf("Skipping %s\n", fileloc)
			return 0, nil
		} else {
			return writeUrlToFile(fileloc, url, mime)
		}
	}
}

func writeUrlToFile(fileloc, url, mime string) (int64, error) {
	fmt.Printf("Writing %s\n", fileloc)
	var out *os.File

	get, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer get.Body.Close()

	if fileExists(fileloc) {
		out, err = os.OpenFile(fileloc, os.O_RDWR, os.ModePerm)
		if err != nil {
			return 0, err
		}
	} else {
		out, err = os.Create(fileloc)
		if err != nil {
			return 0, err
		}
	}
	defer out.Close()

	n, err := io.Copy(out, get.Body)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}
