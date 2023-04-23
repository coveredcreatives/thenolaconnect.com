import { GenerateQRCode } from "./qr"
import { Welcome } from "./welcome"
import { OrderReview } from "./orders"

const pages = {
    Welcome: Welcome,
    GenerateQRCode: GenerateQRCode,
    OrderReview: OrderReview,
}

export default pages