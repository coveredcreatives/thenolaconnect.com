import * as React from 'react';
import api from "../../http"

import { UploadOutlined, } from '@ant-design/icons';
import { Card, Form, Row, Col, Upload, Button, Input, Typography } from 'antd';

const { Text } = Typography;

export function CreateQRForm({ refetch }) {
  const [form] = Form.useForm();
  const [label, setLabel] = React.useState("");
  const [file, setFile] = React.useState();

  const [validationResult, setValidationResult] = React.useState()
  const doesValueContainSpaces = inputValue => /\s/g.test(inputValue)
  const handleLabelChange = e => {
    setLabel(e.target.value)
  }
  const handleFileChange = e => {
    console.log(e);
    setFile(e.file);
    e.onSuccess({ status: 'ok' })
  }
  const handleSubmit = event => {
    var formData = new FormData();
    formData.append("file", file, file.name);
    formData.append("label", label)

    return api.qr_mapping_generate(formData)
      .then((response) => refetch({ options: { throwOnError: false, cancelRefresh: true } }))
      .then((result) => {
        console.log('Success:', result);
      })
      .catch((error) => {
        console.error('Error:', error);
      });
  }

  const getUploadValueFromEvent = (e) => {
    if (Array.isArray(e)) {
      return e;
    }
    return e?.fileList;
  };

  React.useEffect(() => {
    if (doesValueContainSpaces(label)) {
      setValidationResult('error')
    } else if (label) {
      setValidationResult('success')
    }
  }, [label])

  return (
    <Row>
      <Col span={24}>
        <Card display="grid" gridGap={3}>
          <Row>
            <Col>
              <Text strong>Generate QR Code</Text>
            </Col>
          </Row>
          <Row>
            <Col>
              <Text>Please enter the file you want to QRCode.</Text>
            </Col>
          </Row>
          <Row style={{ paddingTop: "10px" }}>
            <Col>
              <Form form={form} onFinish={handleSubmit} layout="vertical">
                <Form.Item
                  name="label"
                  label="QRCode Label"
                  validateStatus={validationResult}
                  help={validationResult === "error" ? "should not contain spaces" : null}
                >
                  <Input
                    type="text"
                    value={label}
                    onChange={handleLabelChange} />
                </Form.Item>
                <Form.Item
                  name="upload"
                  label="Upload"
                  valuePropName="fileList"
                  getValueFromEvent={getUploadValueFromEvent}
                >
                  <Upload maxCount={1} customRequest={handleFileChange}>
                    <Button icon={<UploadOutlined />}>Choose File</Button>
                  </Upload>
                </Form.Item>
                <Form.Item>
                  <Button htmlType="submit">
                    Submit
                  </Button>
                </Form.Item>
              </Form>
            </Col>
          </Row>
        </Card>
      </Col>
    </Row >

  )
}