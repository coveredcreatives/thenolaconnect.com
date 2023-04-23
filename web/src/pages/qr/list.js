import * as React from 'react';
import api, { getApiURL } from "../../http"

import { MinusOutlined } from '@ant-design/icons';
import { Image, Card, Typography, Row, Col, Badge, Button, } from 'antd';

const { Text } = Typography;

export function ListQRMappings({ isLoading, isError, data, refetch }) {
  if (isLoading === true) return "Loading..."

  if (isError === true) return "Error fetching data"

  const handleDelete = (qr_encoded_data) => {
    api.qr_mapping_hide(qr_encoded_data)
      .then(() => refetch())
      .catch((err) => console.err(err))
  }

  const URL = getApiURL();
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
                            <b>Redirects To:</b> <a target="_blank" rel="noreferrer" href={`${URL}/v1/qr_mapping/retrieve?qr_encoded_data=${qrmapping.qr_encoded_data}`}>Link</a>
                          </Text>
                        </Col>
                      </Row>
                      <Row>
                        <Col>
                          <Button icon={<MinusOutlined />} onClick={() => handleDelete(qrmapping.qr_encoded_data)}>Delete</Button>
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