import axios from "axios";

const isLocalHost = window.location.origin.includes('localhost');

export const getApiURL = () => {
    if (isLocalHost) {
        return 'http://127.0.0.1:3001';
    }

    return 'http://127.0.0.1:3001';
};


export const httpClient = axios.create({
    baseURL: getApiURL(),
});

const api = {
    order_communication_generate: () => httpClient.get(`/order_communication/generate`),
    order_communication_list: () => httpClient.get(`/order_communication/list`),
    order_communication_deliver_to_kitchen: (order_id) => httpClient.get(`/order_communication/deliver_to_kitchen?order_id=${order_id}`),
    qr_mapping_generate: (form) => httpClient.post(`/qr_mapping/generate`, form, {
        headers: {
            'content-type': 'multipart/form-data',
        }
    }),
    qr_mapping_list: () => httpClient.get(`/qr_mapping/list`),
    qr_mapping_retrieve: (id) => httpClient.get(`/qr_mapping/retrieve?id=${id}`)
};

export default api;