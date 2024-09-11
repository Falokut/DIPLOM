import { retrieveLaunchParams } from "@telegram-apps/sdk";

export async function UserIsAdmin(): Promise<boolean> {
    const { initData } = retrieveLaunchParams();
    if (initData === undefined || initData.user === undefined) {
        return false;
    }
    const userId = await GetUserIdByTelegramId(initData.user.id);
    if (userId.length == 0) {
        return false;
    }

    let userAdmin = await IsUserAdmin(userId);
    return userAdmin;
}

const getUserIdByTelegramIdEndpoint = 'https://falokut.ru/api/dish_as_a_service/users/get_by_telegram_id/'

export async function GetUserIdByTelegramId(telegramId: number): Promise<string> {
    let resp = await fetch(getUserIdByTelegramIdEndpoint + telegramId)
    if (!resp.ok) {
        console.error(resp)
        return ''
    }
    let userIdResp = await resp.json()
    return userIdResp.userId
}

const isUserAdminUrl = 'https://falokut.ru/api/dish_as_a_service/users/'
export async function IsUserAdmin(userId: string): Promise<boolean> {
    let resp = await fetch(isUserAdminUrl + '/is_admin?userId=' + userId)
    if (!resp.ok) {
        console.error(resp)
        return false
    }
    let isUserAdminResp = await resp.json()
    return isUserAdminResp.isAdmin
}

export function GetUserIdHeader(userId: string): Headers {
    let headers = new Headers();
    headers.set("X-USER-ID", userId)
    return headers
}