import { DefaultClient } from '../utils/client'

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

export async function AddDishCategory(userId: string, name: string): Promise<DishCategory> {
    return await DefaultClient.PostJSON(dishesCategoriesEndpoint,
        {
            name: name,
        },
        DefaultClient.UserAuthHeader(userId),
    ).then(response => response.json()).catch(reason => alert(reason))
}

export async function RenameDishesCategory(userId: string, newName: string, categoryId: number) {
    return await DefaultClient.PostJSON(dishesCategoriesEndpoint + '/' + categoryId,
        {
            name: newName,
        },
        DefaultClient.UserAuthHeader(userId),
    ).catch(reason => alert(reason))
}

export async function DeleteDishesCategory(userId: string, categoryId: number) {
    return await DefaultClient.Delete(
        dishesCategoriesEndpoint + '/' + categoryId,
        null,
        DefaultClient.UserAuthHeader(userId)
    ).catch(reason => alert(reason))
}


