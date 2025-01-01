const cartKey = "cart"
import { initMainButton, MainButton } from '@telegram-apps/sdk';
import { GetUserIdByTelegramId } from './user'
import { DefaultClient } from '../utils/client';
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

const processOrderEndpoint = '/orders'
export async function ProcessOrder(telegramId: number, wishes: string): Promise<boolean> {
    let userId = await GetUserIdByTelegramId(telegramId)
    if (userId.length == 0) {
        return false
    }
    let items = objectFromCart(GetCart())

    let resp = await DefaultClient.PostJSON(processOrderEndpoint, {
        userId: userId,
        paymentMethod: "telegram",
        wishes: wishes,
        items: items,
    })
    if (resp.ok) {
        localStorage.removeItem(cartKey);
    }
    return resp.ok;
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