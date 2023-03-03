import { GenerateQRCode } from "./qr"
import { Welcome } from "./welcome"
import { OrderReview, OrderPlacement } from "./orders"
import { Footer } from "./footer";
import { Navigation } from "./navigation";

const pages = {
    Welcome: Welcome,
    GenerateQRCode: GenerateQRCode,
    OrderPlacement: OrderPlacement,
    OrderReview: OrderReview,
    Footer: Footer,
    Navigation: Navigation,
}

export default pages