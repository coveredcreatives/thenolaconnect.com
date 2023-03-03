import * as React from 'react';
import api, { URL } from "../http"
import { useQuery } from '@tanstack/react-query'

import { UploadOutlined } from '@ant-design/icons';
import { Image, Card, Layout, Typography, Form, Row, Col, Badge, Upload, Button, Input, theme } from 'antd';

const { Text } = Typography;

export function GenerateQRCode() {
    const {
        token: { colorBgContainer },
    } = theme.useToken();

    const { isLoading, isError, data, refetch } = useQuery(['qr_mapping_list'], async () => {
        const response = await api.qr_mapping_list();
        return response.data;
    });

    return (
        <Layout style={{ background: colorBgContainer }}>
            <Layout.Sider style={{ background: colorBgContainer, padding: "10px" }}>
                <CreateQRForm refetch={refetch} />
            </Layout.Sider>
            <Layout.Content style={{ background: colorBgContainer, padding: "10px" }}>
                <ListQRMappings isLoading={isLoading} isError={isError} data={data} />
            </Layout.Content>
        </Layout>
    )
}

function CreateQRForm(refetch) {
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

function ListQRMappings({ isLoading, isError, data }) {
    if (isLoading === true) return "Loading..."

    if (isError === true) return "Error fetching data"

    return (
        <Row>
            <Col>
                {
                    data.map((qrmapping, i) => (
                        <Row key={i}>
                            <Col>
                                <Card>
                                    <Row>
                                        <Col span={16}>
                                            <Image src={qrmapping.qr_file_storage_url} alt={qrmapping.qr_file_storage_url} />
                                        </Col>
                                        <Col span={8}>
                                            <Row>
                                                <Col>
                                                    <Text>
                                                        <b>Label:</b> {qrmapping.name}
                                                    </Text>
                                                </Col>
                                            </Row>

                                            <br />
                                            <Row>
                                                <Col>
                                                    <Text>
                                                        <b>Unique Impressions:</b> <Badge count={qrmapping.unique_impressions} showZero />
                                                    </Text>
                                                </Col>
                                            </Row>
                                            <Row>
                                                <Col>
                                                    <Text>
                                                        <b>Redirects To:</b> <a target="_blank" rel="noreferrer" href={`http://${URL}/qr_mapping/retrieve?qr_encoded_data=${qrmapping.qr_encoded_data}`}>Link</a>
                                                    </Text>
                                                </Col>
                                            </Row>
                                        </Col>
                                    </Row>
                                </Card>
                            </Col>
                        </Row>
                    ))
                }
            </Col>
        </Row >

    )
}