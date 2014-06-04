package tutorials

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type courseStep struct {
	Instruction string `json:"instruction"`
	Cmd         string `json:"cmd"` //command to advance to the next step
}

type courseInfo struct {
	Name  string       `json:"name"`
	Steps []courseStep `json:"steps"`
}

type progressInfo struct {
	courseName  string
	currentStep int
}

type TutorialsInfo struct {
	courses    map[string]courseInfo
	inProgress bool
	progress   progressInfo
}

func New(path string) *TutorialsInfo {
	t := &TutorialsInfo{make(map[string]courseInfo), false, progressInfo{"", 0}}
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		b, err := ioutil.ReadFile(path)

		// marshal json into stuct and append to course list
		var course courseInfo
		json.Unmarshal(b, &course)
		if course.Name != "" {
			t.courses[course.Name] = course
		}
		return nil
	})
	return t
}

func (t *TutorialsInfo) StartCourse(cName string) (instruction string, step string) {
	if t.inProgress == true {
		return "", ""
	}
	if _, ok := t.courses[cName]; ok {
		t.inProgress = true
		t.progress.courseName = cName
		t.progress.currentStep = 0
		return t.courses[cName].Steps[0].Instruction, strconv.Itoa(t.progress.currentStep+1) + "/" + strconv.Itoa(len(t.courses[cName].Steps))
	} else {
		return "", ""
	}

}

func (t *TutorialsInfo) ProgressCourse(commandDone string) (instruction string, courseName string, step string, err error) {
	if t.inProgress {
		if t.courses[t.progress.courseName].Steps[t.progress.currentStep].Cmd == commandDone {
			if t.progress.currentStep < len(t.courses[t.progress.courseName].Steps)-1 {
				t.progress.currentStep++

				if t.progress.currentStep == len(t.courses[t.progress.courseName].Steps)-1 {
					defer t.EndCourse()
				}

				return t.courses[t.progress.courseName].Steps[t.progress.currentStep].Instruction, t.progress.courseName, strconv.Itoa(t.progress.currentStep+1) + "/" + strconv.Itoa(len(t.courses[t.progress.courseName].Steps)), nil

			} else {
				return "", "", "", nil
			}
		} else {
			return "", "", "", errors.New("Please follow the instruction to advance to the next step.")
		}
	} else {
		return "", "", "", errors.New("No tutorial is in progress currently.")
	}
}

func (t *TutorialsInfo) EndCourse() {
	t.inProgress = false
	t.progress.currentStep = 0
	t.progress.courseName = ""
}

func (t *TutorialsInfo) InProgress() bool {
	return t.inProgress
}

func (t *TutorialsInfo) TotalCourses() int {
	total := 0
	for _, _ = range t.courses {
		total++
	}
	return total
}
