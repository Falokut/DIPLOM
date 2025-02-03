import RootIndex from './views/root.svelte'
// -----------------
import AdminIndex from './views/admin/index.svelte'
import AdminAddDishesIndex from './views/admin/add_dishes.svelte'
import AdminDishesCategoriesIndex from './views/admin/dish_сategories.svelte'
import AdminDeleteDishesIndex from './views/admin/delete_dishes.svelte'
import AdminRestaurantsIndex from './views/admin/restaurants.svelte'
// ----------
import PublicIndex from './views/public/index.svelte'
import DishesIndex from './views/public/dishes.svelte'
import CartIndex from './views/public/cart.svelte'
import OrdersIndex from './views/public/orders.svelte'

const routes = [
    {
        name: "/",
        component: RootIndex
    },
    {
        name:"/public-index",
        component: PublicIndex
    },
    {
        name: "/dishes",
        component: DishesIndex,
    },
    {
        name: '/cart',
        component: CartIndex,
    },
    {
        name: '/orders',
        component: OrdersIndex,
    },
    {
        name: '/admin',
        component: AdminIndex,
    },
    {
        name: "/admin/dishes/add",
        component: AdminAddDishesIndex,
    },
    {
        name: "/admin/dishes/delete",
        component: AdminDeleteDishesIndex,
    },
    {
        name: "/admin/dishes/categories",
        component: AdminDishesCategoriesIndex,
    },
    {
        name: "/admin/restaurants",
        component: AdminRestaurantsIndex,
    }
]

export { routes }