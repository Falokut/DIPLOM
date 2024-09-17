import { GetUserIdHeader } from './user'
import { toDates } from 'ts-transformer-dates';

export class UserOrder {
    id: string
    items: UserOrderItem[]
    paymentMethod: string
    total: number
    wishes: string
    createdAt: string
}

export class UserOrderItem {
    dishId: number
    name: string
    price: number
    count: number
    totalPrice: number
    status: string
}

const userOrdersUrl = "https://falokut.ru/api/dish_as_a_service/orders/my"
export async function GetUserOrders(userId: string, offset: number, limit: number): Promise<UserOrder[]> {
    let url = userOrdersUrl + "?limit=" + limit + "&offset=" + offset
    return await fetch(url,
        { headers: GetUserIdHeader(userId) }
    ).then(response => response.json()).catch(reason => alert(reason))
}