import DishesLayout from './views/public/dishes_layout.svelte'
import CartLayout from './views/public/cart_layout.svelte'
// import  from './views/public/index.svelte'


function userIsAdmin() {
    //check if user is admin and returns true or false
}

const routes = [
    {
        name: '/',
        component: DishesLayout,
    },
    {
        name: '/cart',
        component: CartLayout,
    },
    // {
    //     name: 'admin',
    //     component: AdminLayout,
    //     onlyIf: { guard: userIsAdmin, redirect: '/login' },
    //     nestedRoutes: [
    //         { name: 'index', component: AdminIndex },
    //         {
    //             name: 'employees',
    //             component: '',
    //             nestedRoutes: [
    //                 { name: 'index', component: EmployeesIndex },
    //                 { name: 'show/:id', component: EmployeesShow },
    //             ],
    //         },
    //     ],
    // },
]

export { routes }