const fs = require('fs');
const express = require('express');
const { logger, httpLogger } = require('./logger');
const React = require('react');
const ReactPDF = require('@react-pdf/renderer');
const db = require("./db");

const { Document, Page, View, Text, StyleSheet } = ReactPDF;

const styles = StyleSheet.create({
  page: {
    flexDirection: 'column',
    backgroundColor: '#E4E4E4'
  },
  section: {
    margin: 10,
    padding: 10,
    display: "flex",
    flexDirection: "column",
    flexWrap: "wrap",
    justifyContent: "space-between",
    alignItems: "flex-start"
  },
  row: {
    flex: "0 1 auto",
    marginVertical: "5px",
    border: "10px",
    borderColor: "black",
  },
  title: {
    width: "40%",
  },
  answers: {
    width: "60%",
    alignSelf: "flex-end",
  }
});

function questionSortFn(a, b) {
  let item_id_compare = a.item_id.localeCompare(b.item_id);
  if (item_id_compare !== 0) { return item_id_compare };
  return a.question_id.localeCompare(b.question_id);
}

function buildPDF(responses_by_question_id, questions) {
  let collect_elements_by_item_id = {};
  let priority_elements = [];
  let priority_item_ids = [
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
  ];
  let menu_elements = [];
  let menu_item_ids = [];
  questions.sort(questionSortFn);
  questions.forEach((question) => {
    let response = responses_by_question_id[question.question_id];
    if (response === undefined) { return };
    if (!priority_item_ids.includes(question.item_id) && !menu_item_ids.includes(question.item_id)) { menu_item_ids.push(question.item_id); }
    let answers = JSON.parse(response.text_answers);
    let answer = React.createElement(Text, styles.answers, answers.answers[0].value);
    let title = React.createElement(Text, styles.title, question.title);
    let element = React.createElement(View, styles.row, title, answer);
    if (collect_elements_by_item_id[question.item_id] === undefined) { collect_elements_by_item_id[question.item_id] = []; }
    collect_elements_by_item_id[question.item_id].push(element);
  });
  priority_item_ids.forEach((item_id) => {
    if (collect_elements_by_item_id[item_id] === undefined) { return };
    priority_elements.push(...collect_elements_by_item_id[item_id]);
  });
  let view1 = React.createElement(View, { style: styles.section }, ...priority_elements);

  menu_item_ids.forEach((item_id) => {
    if (collect_elements_by_item_id[item_id] === undefined) { return };
    menu_elements.push(...collect_elements_by_item_id[item_id]);
  });
  let view2 = React.createElement(View, { style: styles.section }, ...menu_elements);
  let page = React.createElement(Page, { size: "A4", style: styles.page }, view1, view2);
  return React.createElement(Document, {
    title: "Order Request",
    author: "theneworleansconnection",
    subject: "catering",

  }, page);
};

for (required_var of [process.env.PGHOST, process.env.PGPORT, process.env.PGDATABASE, process.env.PGUSER, process.env.PGPASSWORD]) {
  if (required_var === undefined) {
    logger.log({
      level: 'error',
      message: 'missing required connection variables see: https://www.postgresql.org/docs/9.1/libpq-envars.html'
    });
    process.exit()
  }
}

const server = express();
const port = 8080;
server.use(httpLogger);
server.get('/pdf_generation_worker/:pdf_generation_worker_id', (req, res) => {
  let questions = [];
  let responses_by_question_id = {};
  let form_id = "";
  let set_process_id_promise = db.query('UPDATE order_communication.pdf_generation_worker SET process_id=$1 WHERE pdf_generation_worker_id = $2', [process.pid, req.params.pdf_generation_worker_id]);
  let unset_process_id_promise = db.query('UPDATE order_communication.pdf_generation_worker SET process_id=$1 WHERE pdf_generation_worker_id = $2', [0, req.params.pdf_generation_worker_id]);

  set_process_id_promise.then(_ => {
    return db.query('SELECT form_id FROM order_communication.pdf_generation_worker WHERE pdf_generation_worker_id = $1', [req.params.pdf_generation_worker_id]);
  })
    .then(({ rows }) => {
      form_id = rows[0].form_id;
      return Promise.all([
        db.query('SELECT description, item_id, title, is_question_item, question_id, is_required, form_id FROM google_workspace_forms.item WHERE form_id = $1', [form_id]),
        db.query('SELECT question_id, form_response_id, text_answers FROM google_workspace_forms.answer WHERE form_response_id = (SELECT form_response_id FROM order_communication.pdf_generation_worker WHERE pdf_generation_worker_id = $1)', [req.params.pdf_generation_worker_id])
      ]);
    })
    .then(([fetched_form_items, fetched_response_answers]) => {
      questions = fetched_form_items.rows;
      for (let row of fetched_response_answers.rows) {
        if (row.question_id) {
          responses_by_question_id[row.question_id] = row;
        }
      }
      const pdf = buildPDF(responses_by_question_id, questions);
      return ReactPDF.renderToStream(pdf);
    })
    .then(stream => {
      res.setHeader('Content-Type', 'application/pdf');
      return stream.pipe(res);
    })
    .then(_ => {
      return unset_process_id_promise;
    })
    .catch(e => {
      logger.log({
        level: 'error',
        message: 'failed to generate pdf',
        pdf_generation_worker_id: req.params.pdf_generation_worker_id,
        err: e,
      }); res.status(400).send(e.stack)
    });
});

server.listen(port, () => {
  logger.log({
    level: 'info',
    message: 'pdf generation server running',
    port: port,
  })
});