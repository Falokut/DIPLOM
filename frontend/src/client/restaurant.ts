import { DefaultClient } from '../utils/client'
import { GetAccessToken } from './user'

export class Restaurant {
    id: number
    name: string
}

const restaurantsEndpoint = '/restaurants'
export async function GetAllRestaurants(): Promise<Restaurant[]> {
    return await DefaultClient.Get(restaurantsEndpoint).then(response => response.json())
}



export async function GetRestaurantsById(categoryId: number): Promise<Restaurant[]> {
    return await DefaultClient.Get(restaurantsEndpoint + '/' + categoryId).
        then(response => response.json()).catch(reason => alert(reason))
}

export async function AddRestaurant(name: string): Promise<Restaurant> {
    let accessToken = await GetAccessToken();
    return await DefaultClient.PostJSON(restaurantsEndpoint,
        {
            name: name,
        },
        DefaultClient.UserBearerAuthHeader(accessToken),
    ).then(response => response.json()).catch(reason => alert(reason))
}

export async function RenameRestaurant(newName: string, categoryId: number) {
    let accessToken = await GetAccessToken();
    return await DefaultClient.PostJSON(restaurantsEndpoint + '/' + categoryId,
        {
            name: newName,
        },
        DefaultClient.UserBearerAuthHeader(accessToken),
    ).catch(reason => alert(reason))
}

export async function DeleteRestaurant(categoryId: number) {
    let accessToken = await GetAccessToken();
    return await DefaultClient.Delete(
        restaurantsEndpoint + '/' + categoryId,
        null,
        DefaultClient.UserBearerAuthHeader(accessToken)
    ).catch(reason => alert(reason))
}


