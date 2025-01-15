import { DefaultClient } from '../utils/client'
import { GetAccessToken } from './user'

export class DishCategory {
    id: number
    name: string
}

const allDishesCategoriesEndpoint = '/dishes/all_categories'
export async function GetAllDishesCategories(): Promise<DishCategory[]> {
    return await DefaultClient.Get(allDishesCategoriesEndpoint).then(response => response.json())
}

const dishesCategoriesEndpoint = '/dishes/categories'
export async function GetDishesCategories(): Promise<DishCategory[]> {
    return await DefaultClient.Get(dishesCategoriesEndpoint).then(response => response.json())
}


export async function GetDishCategoriesById(categoryId: number): Promise<DishCategory[]> {
    return await DefaultClient.Get(dishesCategoriesEndpoint + '/' + categoryId).
        then(response => response.json()).catch(reason => alert(reason))
}

export async function AddDishCategory(name: string): Promise<DishCategory> {
    let accessToken = await GetAccessToken();
    return await DefaultClient.PostJSON(dishesCategoriesEndpoint,
        {
            name: name,
        },
        DefaultClient.UserBearerAuthHeader(accessToken),
    ).then(response => response.json()).catch(reason => alert(reason))
}

export async function RenameDishesCategory(newName: string, categoryId: number) {
    let accessToken = await GetAccessToken();
    return await DefaultClient.PostJSON(dishesCategoriesEndpoint + '/' + categoryId,
        {
            name: newName,
        },
        DefaultClient.UserBearerAuthHeader(accessToken),
    ).catch(reason => alert(reason))
}

export async function DeleteDishesCategory(categoryId: number) {
    let accessToken = await GetAccessToken();
    return await DefaultClient.Delete(
        dishesCategoriesEndpoint + '/' + categoryId,
        null,
        DefaultClient.UserBearerAuthHeader(accessToken)
    ).catch(reason => alert(reason))
}


