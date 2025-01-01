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

<div class="dish-container">
  <div class="dish-name">{dish.name} x {count}</div>
  <div class="dish-price">{FormatPriceDefault(dish.price)}</div>
</div>

<style>
  .dish-container {
    display: flex;
    padding: 10px;
    justify-content: space-between;
    align-items: center;
    flex-direction: row;
  }
  .dish-price {
    text-align: right;
  }
</style>
