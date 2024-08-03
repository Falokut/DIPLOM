const cartKey = "cart"
import { initMainButton, MainButton } from '@telegram-apps/sdk';
import { GetUserIdByTelegramId } from './user'
const mainButtonRes = initMainButton();
let mainButton = mainButtonRes[0]

export function GetCart(): Map<string, number> {
    let localCart = localStorage.getItem(cartKey);
    if (localCart === null) {
        return new Map<string, number>()
    }

    return new Map(Object.entries(JSON.parse(localCart)));
}

export function GetDishCount(dishId: number): number {
    let count = GetCart().get(dishId.toString())
    return count == undefined ? 0 : count;
}

export function LoadCart(): MainButton {
    mainButton.setParams({
        text: "Корзина",
    })
    mainButton.enable();
    if (GetCart().size != 0) {
        mainButton.show();
    } else {
        mainButton.hide();
    }
    return mainButton;
}

const processOrderUrl = 'https://falokut.ru/api/dish_as_a_service/orders'
export async function ProcessOrder(telegramId: number): Promise<boolean> {
    let userId = await GetUserIdByTelegramId(telegramId)
    if (userId === -1) {
        return false
    }
    let items = objectFromCart(GetCart())
    let req = {
        userId: userId,
        paymentMethod: "telegram",
        items: items,
    }
    console.log(userId, JSON.stringify(req));

    const headers = new Headers();
    headers.append("Content-Type", "application/json");
    let processOrderOptions = {
        method: "POST",
        headers: headers,
        body: JSON.stringify(req),
    }
    let resp = await fetch(processOrderUrl, processOrderOptions)
    if (resp.ok) {
        localStorage.removeItem(cartKey);
        return true
    }
    return false;
}

export function SetDishCount(dishId: number, count: number) {
    count = Math.max(0, count);
    let cart = GetCart();
    if (count == 0) {
        cart.delete(dishId.toString())
        if (cart.size != 0) {
            mainButton.show();
        } else {
            mainButton.hide();
        }
    } else {
        cart.set(dishId.toString(), count);
        mainButton.show()
    }
    saveCart(cart);
}

function saveCart(cart: Map<string, number>) {
    localStorage.setItem(cartKey, JSON.stringify(objectFromCart(cart)));
}

function objectFromCart(cart: Map<string, number>): any {
    let obj = new Object();
    cart.forEach((v, k) => {
        obj[k.toString()] = v;
    })
    return obj;
}