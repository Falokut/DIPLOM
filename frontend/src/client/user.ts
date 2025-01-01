import { retrieveLaunchParams } from "@telegram-apps/sdk";
import { DefaultClient } from "../utils/client";

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

const getUserIdByTelegramIdEndpoint = '/users/get_by_telegram_id'
export async function GetUserIdByTelegramId(telegramId: number): Promise<string> {
    let resp = await DefaultClient.Get(getUserIdByTelegramIdEndpoint + "/" + telegramId)
    if (!resp.ok) {
        console.error(resp)
        return ''
    }
    let userIdResp = await resp.json()
    return userIdResp.userId
}

const isUserAdminEndpoint = '/users/is_admin'
export async function IsUserAdmin(userId: string): Promise<boolean> {
    let resp = await DefaultClient.Get(isUserAdminEndpoint, {"userId":userId})
    if (!resp.ok) {
        console.error(resp)
        return false
    }
    let isUserAdminResp = await resp.json()
    return isUserAdminResp.isAdmin
}