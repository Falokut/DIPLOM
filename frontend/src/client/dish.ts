
export class Dish {
    id: number
    name: string
    description: string
    categories: string
    url: string
    price: string
}

export function DishFromObject(obj:
    { id: number; name: string; description: string; categories: string; url: string; price: number }): Dish {
    let dish = new Dish();
    dish.id = obj.id;
    dish.name = obj.name;
    dish.description = obj.description;
    dish.categories = obj.categories;
    dish.url = obj.url;
    dish.price = (obj.price / 100).toFixed(2) + '₽';
    return dish;
}

const getDishesUrl = 'https://falokut.ru/api/dish_as_a_service/dishes'
export async function GetDishes(ids: string[] | null) {
    let url = getDishesUrl
    if (ids && ids.length > 0) {
        url += "?" + new URLSearchParams({
            ids: ids.join(',')
        })
    }
    return await fetch(url).then(response => response.json())
}