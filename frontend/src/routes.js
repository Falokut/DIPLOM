import RootIndex from './views/public/root.svelte'
// -----------------
import AdminIndex from './views/admin/index.svelte'
import AdminAddDishesIndex from './views/admin/add_dishes.svelte'
import AdminDishesCategoriesIndex from './views/admin/dish_сategories.svelte'
// ----------
import DishesIndex from './views/public/dishes.svelte'
import CartIndex from './views/public/cart.svelte'

const routes = [
    {
        name: "/",
        component: RootIndex
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
        name: '/admin',
        component: AdminIndex,
    },
    {
        name: "/admin/dishes/add",
        component: AdminAddDishesIndex,
    },
    {
        name: "/admin/dishes/categories",
        component: AdminDishesCategoriesIndex,
    }
]

export { routes }