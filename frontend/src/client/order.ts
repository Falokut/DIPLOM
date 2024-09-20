import { GetUserIdHeader } from './user'
import { GetBackendBasePath } from '../main'

export class UserOrder {
    id: string
    items: UserOrderItem[]
    paymentMethod: string
    total: number
    wishes: string
    createdAt: string
    status: string
}

export class UserOrderItem {
    dishId: number
    name: string
    price: number
    count: number
    totalPrice: number
}

const userOrdersEndpoint = "/orders/my"
export async function GetUserOrders(userId: string, offset: number, limit: number): Promise<UserOrder[]> {
    let url = GetBackendBasePath() + userOrdersEndpoint + "?limit=" + limit + "&offset=" + offset
    return await fetch(url,
        { headers: GetUserIdHeader(userId) }
    ).then(response => response.json()).catch(reason => alert(reason))
}