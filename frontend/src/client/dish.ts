import { GetUserIdHeader } from './user'
import { GetBackendBasePath } from '../main'

export class Dish {
    id: number
    name: string
    description: string
    categories: string
    url: string
    price: number
}

const dishesEndpoint = '/dishes'
export async function GetDishes(ids: undefined | string[], limit, offset, categoryId) {
    let url = GetBackendBasePath() + dishesEndpoint
    if (ids && ids.length > 0) {
        url += "?ids=" + ids.join(',')
    } else {
        url += "?limit=" + limit + "&offset=" + offset
        if (categoryId) {
            url += "&categoriesIds=" + categoryId
        }
    }
    return await fetch(url).then(response => response.json()).catch(reason => console.log(reason))
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
    return await fetch(GetBackendBasePath() + dishesEndpoint, {
        method: "POST",
        headers: headers,
        body: JSON.stringify(dish)
    }).then(resp => resp.ok).catch(reason => alert(reason))
}