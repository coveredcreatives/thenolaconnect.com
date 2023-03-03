import * as React from 'react';
import { DateTime } from 'luxon';
import { Pagination, Card, Row, Col, Button, Typography } from 'antd';
import api from '../http';

export default function Grid({ data, isLoading, isError, refetch, asyncInProgress, setAsyncInProgress }) {
    let [orderNumber, setOrderNumber] = React.useState(0);

    if (!data?.length) {
        return null
    }
    let orders = data;

    const onPageChange = (page, pageSize) => {
        setOrderNumber(page % pageSize);
    }

    const deliverToKitchen = (evt) => {
        setAsyncInProgress(true);
        api.order_communication_deliver_to_kitchen(orders[orderNumber].order_id)
            .catch((e) => console.error(e))
            .finally(() => { setAsyncInProgress(false); refetch() });
    }

    return (
        orders.length === 0 || orderNumber >= orders.length ? null : (
            <div>
                <Card>
                    <Row>
                        <Col span={8}>
                            <Row>
                                <Col span={24}>
                                    <Row>
                                        <Col>
                                            <Typography.Text>Order #{orderNumber + 1} / {orders.length}</Typography.Text>
                                        </Col>
                                    </Row>
                                    <Row>
                                        <Col>
                                            <Typography.Text>Created At: {DateTime.fromISO(orders[orderNumber].created_at).toLocaleString(DateTime.DATETIME_MED)}</Typography.Text>
                                        </Col>
                                    </Row>
                                    <Row>
                                        <Col>
                                            <Typography.Text>{orders[orderNumber].is_delivered_to_kitchen ? "Sent To Kitchen" : !orders[orderNumber].is_viewed_by_manager ? "Awaiting Manager Response" : orders[orderNumber].is_rejected_by_manager ? "Order Rejected" : orders[orderNumber].is_accepted_by_manager ? "Order Accepted" : null}</Typography.Text>
                                        </Col>
                                    </Row>
                                    <Row>
                                        <Col>
                                            <br />
                                        </Col>
                                    </Row>
                                    <Row>
                                        <Col>
                                            <Button shape={"round"} onClick={deliverToKitchen} loading={asyncInProgress} disabled={asyncInProgress || orders[orderNumber].is_rejected_by_manager}>
                                                <Typography.Text>{orders[orderNumber].is_delivered_to_kitchen ? "Resend" : "Deliver to Kitchen"}</Typography.Text>
                                            </Button>
                                        </Col>
                                    </Row>
                                </Col>
                            </Row>
                        </Col>
                        <Col span={16}>
                            <Row>
                                <Col span={24}>
                                    <iframe title="Order Preview" src={orders[orderNumber].order_file_storage_url} height="100%" width="100%"></iframe>
                                </Col>
                            </Row>
                            <Row>
                                <Col span={8}>
                                    <a href={orders[orderNumber].order_file_storage_url} target="_blank" rel="noopener noreferrer">Open</a>
                                </Col>
                            </Row>
                        </Col>
                    </Row>
                </Card >
                <br />
                <Row justify={"center"}>
                    <Col span={16}>
                        <Pagination simple total={orders.length} defaultCurrent={orderNumber + 1} onChange={onPageChange} hideOnSinglePage />
                    </Col>
                </Row>
            </div >
        ))
}