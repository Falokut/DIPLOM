<script lang="ts">
  import { onMount } from "svelte";
  import { GetDishCount } from "../../client/cart";
  import { Dish } from "../../client/dish";
  import { FormatPriceDefault } from "../../utils/format_price";

  export let dish: Dish;
  let count = 0;
  onMount(() => {
    count = GetDishCount(dish.id);
    dish.price = dish.price * count;
  });
</script>

<div class="dish_container">
  <div class="dish_name">{dish.name} x {count}</div>
  <div class="spacer"></div>
  <div class="dish_price">{FormatPriceDefault(dish.price)}</div>
</div>

<style>
  :root {
    --cart-item-height: 32px;
  }

  .dish_container {
    width: 100vw;
    height: var(--cart-item-height);
    display: inline-flex;
  }
  .dish_name {
    padding-left: var(--container-padding-lr);
  }

  .dish_price {
    text-align: right;
    padding-right: var(--container-padding-lr);
  }
</style>
