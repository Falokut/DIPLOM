import { GetBackendBasePath } from '../main'
import { DefaultClient } from '../utils/client'

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
    return await DefaultClient.Get(
        userOrdersEndpoint,
        { "limit": limit, "offset": offset },
        DefaultClient.UserAuthHeader(userId)
    ).
        then(response => response.json()).
        catch(reason => alert(reason))
}