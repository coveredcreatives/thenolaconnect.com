import * as React from 'react';
import { OrderSheet, OrderSheetSider } from '../components/ordersheet';
import { useQuery } from '@tanstack/react-query'
import api from '../http';
import { Row, Col, Layout, theme } from 'antd';

export function OrderPlacement() {
    return (
        <Layout.Content>
            <Row justify="center">
                <Col offset={4} span={16}>
                    <iframe title="orders" src="https://docs.google.com/forms/d/e/1FAIpQLSeKe8iSippG-8wLxdPaGrL2Bpbqw6O8lofNN6gti98MX-YzOw/viewform?embedded=true" width="640" height="640" overflow="scroll" frameborder="0" marginheight="0" marginwidth="0">Loading…</iframe>
                </Col>
            </Row>
        </Layout.Content>
    )
}

export function OrderReview() {
    const {
        token: { colorBgContainer },
    } = theme.useToken();

    const { isLoading, isError, data, refetch } = useQuery(['order_communication_list'], async () => {
        const response = await api.order_communication_list();
        return response.data;
    });
    let [filterBy, setFilterBy] = React.useState("none");
    let [asyncInProgress, setAsyncInProgress] = React.useState(false);
    const orderOptions = {
        refetch: refetch,
        data: data,
        isLoading: isLoading,
        isError: isError,
        filterBy: filterBy,
        setFilterBy: setFilterBy,
        asyncInProgress: asyncInProgress,
        setAsyncInProgress: setAsyncInProgress,
    };


    return (
        <Layout style={{ background: colorBgContainer }}>
            <Layout.Sider style={{ background: colorBgContainer, padding: "10px" }}>
                <OrderSheetSider {...orderOptions} />
            </Layout.Sider>
            <Layout.Content style={{ background: colorBgContainer, padding: "10px" }}>
                <OrderSheet {...orderOptions} />
            </Layout.Content>
        </Layout >

    )
}