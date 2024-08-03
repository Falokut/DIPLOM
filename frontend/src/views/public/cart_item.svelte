<script>
  import { onMount } from "svelte";
  import { SetDishCount, GetDishCount } from "../../client/cart";
  import { dishFromObject } from "../../client/dish";

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
    dish = dishFromObject(dishInst);
    count = GetDishCount(dish.id);
  });
</script>

<div>
  <img class="dish-preview" src={dish.url} alt={dish.name} />
  <div class="dish-caption">{dish.name} {dish.price}</div>
  <div class="dish-count">{count}</div>
</div>

<style>
  :root {
    --cart-item-height: 32px;
  }
  .dish-preview {
    width: auto;
    height: var(--cart-item-height);
    object-fit: scale-down;
  }
</style>
