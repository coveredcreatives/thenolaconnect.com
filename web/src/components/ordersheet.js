import * as React from 'react'
import Grid from './grid'
import { Row, Col, Radio, Space, Typography, Card } from 'antd';

export function OrderSheetSider({ data, filterBy, setFilterBy }) {
    return (
        <div>
            <Row>
                <Col>
                    <Radio.Group onChange={(e) => setFilterBy(e.target.value)} value={filterBy}>
                        <Space direction="vertical">
                            <Radio value={"unread"}>Unread ({data?.filter(applyFilterBy("unread")).length || 0})</Radio>
                            <Radio value={"rejected"}>Rejected ({data?.filter(applyFilterBy("rejected")).length || 0})</Radio>
                            <Radio value={"deliver_to_kitchen"}>Deliver To Kitchen ({data?.filter(applyFilterBy("deliver_to_kitchen")).length || 0})</Radio>
                            <Radio value={"none"}>No Filter ({data?.length || 0})</Radio>
                        </Space>
                    </Radio.Group>
                </Col>
            </Row>
        </div>
    )
}

function applyFilterBy(selectFilterBy) {
    return (value, index, array) => {
        if (selectFilterBy === "none") {
            return true;
        }
        if (selectFilterBy === "unread") {
            return !value.is_viewed_by_manager;
        }
        if (selectFilterBy === "rejected") {
            return value.is_viewed_by_manager && !value.is_accepted_by_manager;
        }
        if (selectFilterBy === "deliver_to_kitchen") {
            return value.is_accepted_by_manager && !value.is_delivered_to_kitchen;
        }
        return false;
    }
}

export function OrderSheet({ data, isLoading, isError, refetch, filterBy, setFilterBy, asyncInProgress, setAsyncInProgress }) {



    if (isLoading === true) { return "Loading recorded orders" }

    if (isError === true) { return "Failed to load recorded orders" }

    if (data) {
        let orders = filterBy === "none" ? data :
            filterBy === "unread" ? data.filter(applyFilterBy("unread")) :
                filterBy === "rejected" ? data.filter(applyFilterBy("rejected")) :
                    filterBy === "deliver_to_kitchen" ? data.filter(applyFilterBy("deliver_to_kitchen")) : null;

        return (
            <Card>
                {
                    orders.length > 0 ? (
                        <Row>
                            <Col span={24}>
                                <Grid data={orders} isLoading={isLoading} refetch={refetch} asyncInProgress={asyncInProgress} setAsyncInProgress={setAsyncInProgress} />
                            </Col>
                        </Row>
                    ) : <Typography.Text>No Orders to display</Typography.Text>
                }
            </Card>

        )
    } else {
        return (
            <Card>
                <Row>
                    <Col span={24}>
                        <Typography.Text>No Orders to display</Typography.Text>
                    </Col>
                </Row>
            </Card>
        )
    }


}