package downloadpage

import (
	"tasklist"
)

func getDownloadTask() tasklist.PageObject {
	return tasklist.PopPageTask()
}

