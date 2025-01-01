import { DefaultClient } from '../utils/client'

export class Dish {
    id: number
    name: string
    description: string
    categories: string
    url: string
    price: number
}

const dishesEndpoint = '/dishes'
export async function GetDishes(dishIds?: any[], limit?, offset?, categoriesIds?) {
    let queryParams = {}
    if (dishIds && dishIds.length > 0) {
        queryParams = { "ids": dishIds.join(',') }
    } else {
        queryParams = { "limit": limit, offset: offset }
        if (categoriesIds) queryParams["categoriesIds"] = categoriesIds
    }
    return await DefaultClient.Get(dishesEndpoint, queryParams).
        then(response => response.json()).
        catch(reason => console.log(reason))
}

export class AddDishObj {
    name: string
    description: string
    price: number
    categories: number[]
    image: any
}

export async function AddDish(dish: AddDishObj, userId: string): Promise<boolean | void> {
    return await DefaultClient.PostJSON(dishesEndpoint, dish, DefaultClient.UserAuthHeader(userId)).
        then(resp => resp.ok).
        catch(reason => alert(reason))
}