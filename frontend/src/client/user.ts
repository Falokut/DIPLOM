import { retrieveLaunchParams } from "@telegram-apps/sdk";
import { DefaultClient } from "../utils/client";

const accessTokenKey = "access";
const accessTokenExpiresAtKey = "access_expire_at";
const refreshTokenKey = "refresh";
const refreshTokenExpiresAtKey = "refresh_expire_at";

export async function GetAccessToken() {
    const accessToken = localStorage.getItem(accessTokenKey);
    const accessTokenExpiresAt = localStorage.getItem(accessTokenExpiresAtKey)
    if (accessToken && accessTokenExpiresAt && new Date(accessTokenExpiresAt) > new Date()) {
        return accessToken;
    }
    const refreshToken = localStorage.getItem(refreshTokenKey);
    const refreshTokenExpiresAt = localStorage.getItem(refreshTokenExpiresAtKey)
    if (refreshToken && refreshTokenExpiresAt && new Date(refreshTokenExpiresAt) > new Date()) {
        return refreshAccessToken(refreshToken);
    }

    return await auth();
}

const authByTelegramEndpoint = '/auth/login_by_telegram'
export async function auth(): Promise<string> {
    const { initDataRaw } = retrieveLaunchParams();
    if (initDataRaw == undefined) {
        throw "undefined init data!";
    }
    let resp = await DefaultClient.PostJSON(authByTelegramEndpoint,
        {
            initTelegramData: initDataRaw
        })
    if (!resp.ok) {
        return ""
    }
    const jsonResp = await resp.json();
    localStorage.setItem(accessTokenKey, jsonResp.accessToken.token);
    localStorage.setItem(accessTokenExpiresAtKey, jsonResp.accessToken.expiresAt);
    localStorage.setItem(refreshTokenKey, jsonResp.refreshToken.token);
    localStorage.setItem(refreshTokenExpiresAtKey, jsonResp.refreshToken.expiresAt);
    return jsonResp.accessToken.token;
}

async function refreshAccessToken(refreshToken: string): Promise<string> {
    let resp = await DefaultClient.Get(authByTelegramEndpoint, null, DefaultClient.UserBearerAuthHeader(refreshToken))
    if (!resp.ok) {
        return ""
    }
    const jsonResp = await resp.json();
    localStorage.setItem(accessTokenKey, jsonResp.accessToken.token);
    localStorage.setItem(accessTokenExpiresAtKey, jsonResp.accessToken.expiresAt);
    return jsonResp.accessToken.token;
}

const isUserAdminEndpoint = '/has_admin_privileges'
export async function IsUserAdmin(): Promise<boolean> {
    let accessToken = await GetAccessToken()
    if (!accessToken) {
        return false;
    }
    let resp = await DefaultClient.Get(isUserAdminEndpoint, null, DefaultClient.UserBearerAuthHeader(accessToken))
    if (!resp.ok) {
        return false
    }
    let isUserAdminResp = await resp.json()
    return isUserAdminResp.hasAdminPrivileges
}