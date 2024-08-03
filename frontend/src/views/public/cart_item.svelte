<script>
  import { onMount } from "svelte";
  import { GetDishCount } from "../../client/cart";
  import { DishFromObject } from "../../client/dish";

  export let dishInst = {
    id: 0,
    url: "",
    name: "",
    description: "",
    categories: "",
    price: 0,
  };

  let dish = {};
  let count = 0;
  onMount(() => {
    count = GetDishCount(dishInst.id);
    dishInst.price = dishInst.price * count;

    dish = DishFromObject(dishInst);
  });
</script>

<div class="dish">
  <div class="dish-name">{dish.name} x {count}</div>
  <div class="dish-price">{dish.price}</div>
</div>

<style>
  :root {
    --cart-item-height: 32px;
  }

  .dish {
    width: 100vw;
    height: var(--cart-item-height);
    display: grid;
    grid-auto-flow: column;
  }
  .dish-name {
    text-align: left;
    padding-left: 10px;
  }

  .dish-price {
    text-align: right;
    padding-right: 10px;
  }
</style>
