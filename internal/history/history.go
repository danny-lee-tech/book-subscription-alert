package history

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

type History struct {
	key      string
	maxItems int
	file     string
}

func Init(key string, maxItems int) *History {
	newFile := key + ".txt"
	return &History{
		key:      key,
		maxItems: maxItems,
		file:     newFile,
	}
}

func (history *History) CheckIfExists(item string) (bool, error) {
	file, err := os.OpenFile(history.file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Sprintf("Error while closing resource %q: %+v", file.Name(), err))
		}
	}()

	lastPostIdsBytes, err := io.ReadAll(file)
	if err != nil {
		return false, err
	}
	lastPostsString := string(lastPostIdsBytes)
	lastPostItems := strings.Split(lastPostsString, "\n")

	if slices.Contains(lastPostItems, item) {
		return true, nil
	}

	return false, nil
}

func (history *History) RecordItemIfNotExist(item string) (bool, error) {
	file, err := os.OpenFile(history.file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Sprintf("Error while closing resource %q: %+v", file.Name(), err))
		}
	}()

	lastPostIdsBytes, err := io.ReadAll(file)
	if err != nil {
		return false, err
	}
	lastPostsString := string(lastPostIdsBytes)
	lastPostItems := strings.Split(lastPostsString, "\n")

	if slices.Contains(lastPostItems, item) {
		return false, nil
	}

	lastPostItems = append(lastPostItems, item)

	if _, err := file.Seek(0, 0); err != nil {
		panic(err)
	}
	if err := file.Truncate(0); err != nil {
		panic(err)
	}

	if len(lastPostItems) > history.maxItems {
		elementsToDelete := len(lastPostItems) - history.maxItems
		lastPostItems = lastPostItems[elementsToDelete:]
	}
	_, err = file.Write([]byte(strings.Join(lastPostItems, "\n")))
	if err != nil {
		return false, err
	}

	return true, nil
}
