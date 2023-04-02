package pgw

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"sort"

	"google.golang.org/api/forms/v1"
	"gorm.io/gorm"

	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/model"
	template_tools "github.com/coveredcreatives/thenolaconnect.com/pkg/template"
)

type element struct {
	Question model.Item
	Answer   string
}

func BuildOrderHTML(db *gorm.DB, pgw model.PDFGenerationWorker) (b []byte, err error) {
	defer alog.Trace("BuildOrderHTML").Stop(&err)
	pid := os.Getpid()
	err = db.
		Model(&model.PDFGenerationWorker{}).
		Where(&model.PDFGenerationWorker{PDFGenerationWorkerId: pgw.PDFGenerationWorkerId}).
		Update("process_id", pid).Error
	if err != nil {
		return
	}
	defer db.
		Model(&model.PDFGenerationWorker{}).
		Where(&model.PDFGenerationWorker{PDFGenerationWorkerId: pgw.PDFGenerationWorkerId}).
		Update("process_id", 0)
	items := []model.Item{}
	err = db.Model(&model.Item{FormId: pgw.FormId}).Find(&items).Error
	if err != nil {
		return
	}
	answers := []model.Answer{}
	err = db.Model(&model.Answer{FormId: pgw.FormId}).Find(&answers).Error
	if err != nil {
		return
	}
	responses_by_question_id := map[string]model.Answer{}
	for _, response := range answers {
		responses_by_question_id[response.QuestionId] = response
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].ItemId == items[j].ItemId {
			return items[i].QuestionId < items[j].QuestionId
		} else {
			return items[i].ItemId < items[j].QuestionId
		}
	})
	priority_item_ids := []string{
		"6afe4bff",
		"2f9abe9f",
		"5ea1bf77",
		"07f0edcd",
		"147e7d7d",
		"24e5d9f5",
		"5ed65d27",
		"73b1cac1",
		"2a280fad",
		"23bae6e0",
		"6bec128a",
	}
	menu_item_ids := []string{}
	elements_by_item_id := map[string][]element{}
	elements := []element{}
	for _, question := range items {
		response, ok := responses_by_question_id[question.QuestionId]
		if !ok {
			continue
		}
		item_id_located := false
		for _, id := range append(priority_item_ids, menu_item_ids...) {
			item_id_located = question.ItemId == id
			if item_id_located {
				break
			}
		}
		if !item_id_located {
			menu_item_ids = append(menu_item_ids, question.ItemId)
		}
		answer := forms.TextAnswers{}
		err = json.Unmarshal([]byte(response.TextAnswers), &answer)
		if err != nil {
			return
		}
		elements_by_item_id[question.ItemId] = append(elements_by_item_id[question.ItemId], element{
			Question: question,
			Answer:   answer.Answers[0].Value,
		})
	}
	for _, id := range priority_item_ids {
		elements = append(elements, elements_by_item_id[id]...)
	}
	for _, id := range menu_item_ids {
		elements = append(elements, elements_by_item_id[id]...)
	}
	ordersheet, err := template_tools.OrderSheetTemplate()
	if err != nil {
		return
	}
	outfile_title := fmt.Sprint("outfile_template_", pgw.PDFGenerationWorkerId)
	ordersheet_template, err := template.New("ordersheet").Parse(ordersheet)
	f, err := os.Create(fmt.Sprint(outfile_title, ".html"))
	if err != nil {
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = ordersheet_template.Execute(w, elements)
	w.Flush()

	_, err = exec.Command("wkhtmltopdf", fmt.Sprint(outfile_title, ".html"), fmt.Sprint(outfile_title, ".pdf")).Output()
	if err != nil {
		return
	}

	s, err := os.ReadFile(fmt.Sprint(outfile_title, ".pdf"))
	if err != nil {
		return
	}

	b = []byte(s)
	return
}
