import { GetUserIdHeader } from './user'
export class Dish {
    id: number
    name: string
    description: string
    categories: string
    url: string
    price: number
}

const dishesUrl = 'https://falokut.ru/api/dish_as_a_service/dishes'
export async function GetDishes(ids: string[] | null) {
    let url = dishesUrl
    if (ids && ids.length > 0) {
        url += "?" + new URLSearchParams({
            ids: ids.join(',')
        })
    }
    return await fetch(url).then(response => response.json())
}

export class AddDishObj {
    name: string
    description: string
    price: number
    categories: number[]
    image: any
}

export async function AddDish(dish: AddDishObj, userId: string): Promise<boolean | void> {
    let headers = GetUserIdHeader(userId);
    headers.set("content-type", "application/json; charset=utf8")
    return await fetch(dishesUrl, {
        method: "POST",
        headers: headers,
        body: JSON.stringify(dish)
    }).then(resp => resp.ok).catch(reason => alert(reason))
}