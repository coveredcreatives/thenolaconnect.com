import { GenerateQRCode } from "./qr"
import { Welcome } from "./welcome"
import { OrderReview, OrderPlacement } from "./orders"

const pages = {
    Welcome: Welcome,
    GenerateQRCode: GenerateQRCode,
    OrderPlacement: OrderPlacement,
    OrderReview: OrderReview,
}

export default pages