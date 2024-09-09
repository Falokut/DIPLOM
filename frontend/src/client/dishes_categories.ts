import { GetUserIdHeader } from './user'
export class DishCategory {
    id: number
    name: string
}

const dishesCategoriesUrl = 'https://falokut.ru/api/dish_as_a_service/dishes/categories'
export async function GetDishesCategories(): Promise<DishCategory[]> {
    return await fetch(dishesCategoriesUrl).then(response => response.json())
}

export async function GetDishCategoriesById(categoryId: number): Promise<DishCategory> {
    return await fetch(dishesCategoriesUrl + '/' + categoryId).then(response => response.json()).catch(reason => alert(reason))
}

export async function AddDishesCategory(userId: string, name: string): Promise<DishCategory[]> {
    let req = {
        name: name,
    }
    return await fetch(dishesCategoriesUrl, {
        method: "POST",
        headers: GetUserIdHeader(userId),
        body: JSON.stringify(req)
    }).then(response => response.json()).catch(reason => alert(reason))
}

export async function RenameDishesCategory(userId: string, newName: string, categoryId: number) {
    let req = {
        name: newName,
    }
    return await fetch(dishesCategoriesUrl + '/' + categoryId, {
        method: "POST",
        headers: GetUserIdHeader(userId),
        body: JSON.stringify(req)
    }).catch(reason => alert(reason))
}

export async function DeleteDishesCategory(userId: string, categoryId: number) {
    return await fetch(dishesCategoriesUrl + '/' + categoryId, {
        method: "DELETE",
        headers: GetUserIdHeader(userId),
    }).catch(reason => alert(reason))
}


