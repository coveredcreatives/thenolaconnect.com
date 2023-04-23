import * as React from 'react';
import api from "../../http"
import { useQuery } from '@tanstack/react-query'

import { Layout, theme } from 'antd';
import { CreateQRForm } from './create';
import { ListQRMappings } from './list';

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
                <ListQRMappings isLoading={isLoading} isError={isError} data={data} refetch={refetch} />
            </Layout.Content>
        </Layout>
    )
}