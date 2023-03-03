import axios from "axios";

export const localURI = `localhost:3001`;

export const baseURI = `us-central1-the-new-orleans-connection.cloudfunctions.net`;

export const URL = (process.env.NODE_ENV === 'production') ? `https://${baseURI}` : `http://${localURI}`;

export const httpClient = axios.create({
    baseURL: URL,
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