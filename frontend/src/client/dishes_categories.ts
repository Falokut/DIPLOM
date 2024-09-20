import { GetUserIdHeader } from './user'
import { GetBackendBasePath } from '../main'

export class DishCategory {
    id: number
    name: string
}

const dishesCategoriesEndpoint = '/dishes/categories'
export async function GetDishesCategories(): Promise<DishCategory[]> {
    return await fetch(GetBackendBasePath() + dishesCategoriesEndpoint).then(response => response.json())
}

export async function GetDishCategoriesById(categoryId: number): Promise<DishCategory[]> {
    return await fetch(GetBackendBasePath() + dishesCategoriesEndpoint + '/' + categoryId).then(response => response.json()).catch(reason => alert(reason))
}

export async function AddDishCategory(userId: string, name: string): Promise<DishCategory> {
    let req = {
        name: name,
    }
    let headers = GetUserIdHeader(userId);
    headers.set("content-type", "application/json; charset=utf8")
    return await fetch(GetBackendBasePath() + dishesCategoriesEndpoint, {
        method: "POST",
        headers: headers,
        body: JSON.stringify(req)
    }).then(response => response.json()).catch(reason => alert(reason))
}

export async function RenameDishesCategory(userId: string, newName: string, categoryId: number) {
    let req = {
        name: newName,
    }
    let headers = GetUserIdHeader(userId);
    headers.set("content-type", "application/json; charset=utf8")
    return await fetch(GetBackendBasePath() + dishesCategoriesEndpoint + '/' + categoryId, {
        method: "POST",
        headers: headers,
        body: JSON.stringify(req)
    }).catch(reason => alert(reason))
}

export async function DeleteDishesCategory(userId: string, categoryId: number) {
    return await fetch(GetBackendBasePath() + dishesCategoriesEndpoint + '/' + categoryId, {
        method: "DELETE",
        headers: GetUserIdHeader(userId),
    }).catch(reason => alert(reason))
}


