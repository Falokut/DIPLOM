
const getUserIdByTelegramIdEndpoint = 'https://falokut.ru/api/dish_as_a_service/users/get_by_telegram_id/'

export async function GetUserIdByTelegramId(telegramId): Promise<number> {
    console.log(telegramId)
    let resp = await fetch(getUserIdByTelegramIdEndpoint + telegramId)
    if (!resp.ok) {
        console.error(resp)
        return -1
    }
    let userIdResp = await resp.json()
    return userIdResp.userId
}