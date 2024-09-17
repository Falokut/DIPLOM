export const DefaultCurrency = "₽"
export function FormatPrice(price: number, currency: string): string {
    return (price / 100).toFixed(2) + currency
}
export function FormatPriceDefault(price: number): string {
    return (price / 100).toFixed(2) + DefaultCurrency
}