import { DefaultClient } from '../utils/client'
import { GetAccessToken } from './user'

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

const processOrderEndpoint = '/orders'
export async function ProcessOrder(items: any, wishes: string): Promise<boolean> {
    const accessToken = await GetAccessToken();

    let resp = await DefaultClient.PostJSON(processOrderEndpoint,
        {
            paymentMethod: "telegram",
            wishes: wishes,
            items: items,
        },
        DefaultClient.UserBearerAuthHeader(accessToken))
    return resp.ok;
}

const userOrdersEndpoint = "/orders/my"
export async function GetUserOrders(offset: number, limit: number): Promise<UserOrder[]> {
    const accessToken = await GetAccessToken();
    return await DefaultClient.Get(
        userOrdersEndpoint,
        { "limit": limit, "offset": offset },
        DefaultClient.UserBearerAuthHeader(accessToken)
    ).
        then(response => response.json()).
        catch(reason => alert(reason))
}