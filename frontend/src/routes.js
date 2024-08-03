import RootLayout from './views/public/root_layout.svelte'
import { UserIsAdmin } from './client/user';
// -----------------
import AdminLayout from './views/admin/layout.svelte'
import AdminIndex from './views/admin/index.svelte'
import AdminDishesIndex from './views/admin/dishes.svelte'
// ----------
import PublicLayout from './views/public/layout.svelte';
import DishesIndex from './views/public/dishes.svelte'
import CartIndex from './views/public/cart.svelte'

const routes = [
    {
        name: "/",
        layout: RootLayout
    },
    {
        name: "/dishes",
        component: DishesIndex,
        layout: PublicLayout
    },
    {
        name: '/cart',
        component: CartIndex,
        layout: PublicLayout
    },
    {
        name: '/admin',
        component: AdminIndex,
        layout: AdminLayout,
        onlyIf: { guard: UserIsAdmin, redirect: '/' },
        nestedRoutes: [
            {
                name: "dishes",
                component: AdminDishesIndex,
                layout: AdminLayout
            }
        ]
    }
]

export { routes }