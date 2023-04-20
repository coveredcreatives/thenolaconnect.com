import axios from "axios";

const isLocalHost = window.location.origin.includes('localhost');

export const getApiURL = () => {
    if (isLocalHost) {
        return 'http://127.0.0.1:8080';
    }

    return 'https://api.thenolaconnect.com';
};


export const httpClient = axios.create({
    baseURL: getApiURL(),
});

const api = {
    order_communication_generate: () => httpClient.get(`/v1/order_communication/generate`),
    order_communication_list: () => httpClient.get(`/v1/order_communication/list`),
    order_communication_deliver_to_kitchen: (order_id) => httpClient.get(`/v1/order_communication/deliver_to_kitchen?order_id=${order_id}`),
    qr_mapping_generate: (form) => httpClient.post(`/v1/qr_mapping/generate`, form, {
        headers: {
            'content-type': 'multipart/form-data',
        }
    }),
    qr_mapping_list: () => httpClient.get(`/v1/qr_mapping/list`),
    qr_mapping_retrieve: (id) => httpClient.get(`/v1/qr_mapping/retrieve?id=${id}`)
};

export default api;