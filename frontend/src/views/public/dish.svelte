<script lang="ts">
  import { onMount } from "svelte";
  import { SetDishCount, GetDishCount } from "../../client/cart";
  import PreviewImage from "../components/preview_image.svelte";
  import { Dish } from "../../client/dish";
  import { FormatPriceDefault } from "../../utils/format_price";

  export let dish: Dish;
  let count = 0;

  function incrDishCount() {
    count = Math.min(count + 1, 40);
    SetDishCount(dish.id, count);
  }

  let showIncrDecrButtons = false;
  function decrDishCount() {
    count = Math.max(count - 1, 0);
    if (count == 0) {
      showIncrDecrButtons = false;
    }
    SetDishCount(dish.id, count);
  }

  onMount(() => {
    count = GetDishCount(dish.id);
    if (count > 0) {
      showIncrDecrButtons = true;
    }
  });
</script>

<div class="dish">
  <div class="dish-img-count">
    <div class="dish-image">
      <PreviewImage bind:url={dish.url} bind:alt={dish.name} />
    </div>
    {#if count > 0}
      <div class="dish-count">{count}</div>
    {/if}
  </div>
  <div class="dish-caption">{dish.name}</div>
  <div class="dish-caption">{FormatPriceDefault(dish.price)}</div>
  <div class="dish-count-buttons">
    {#if showIncrDecrButtons}
      <button
        class="btn dish-count-button"
        id="dish-count-incr-button"
        on:click={incrDishCount}>+</button
      >
      <button
        class="btn dish-count-button"
        id="dish-count-decr-button"
        on:click={decrDishCount}>-</button
      >
    {:else}
      <button
        class="dish-add btn"
        on:click={() => {
          count = 1;
          showIncrDecrButtons = true;
          SetDishCount(dish.id, count);
        }}>добавить</button
      >
    {/if}
  </div>
</div>

<style>
  .dish-image {
    width: 150px;
    height: 150px;
  }
  .dish-img-count {
    display: flex;
    position: relative;
    justify-content: center;
    width: 100%;
  }
  .dish-count {
    position: absolute;
    border-radius: 40%;
    background-color: orange;
    width: 36px;
    height: 36px;
    right: 0px;
    top: 0px;
    font-weight: bold;
    font-size: medium;
    text-align: center;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .dish-caption {
    margin-top: auto;
    font-size: medium;
  }

  .dish {
    background-color: var(--tg-theme-secondary-bg-color);
    border-radius: 5%;
    padding: 10px;
    margin: 5px;
    text-align: center;
    display: inline-block;
  }

  .dish-count-buttons {
    display: flex;
    flex-direction: row;
    justify-content: space-around;
    text-align: center;
    font-size: medium;
    font-weight: bold;
  }

  .dish-count-button {
    height: 100%;
  }

  .dish-add {
    width: 100%;
    height: 100%;
    text-align: center;
  }

  .dish-count-buttons #dish-count-incr-button {
    background-color: rgb(17, 151, 17);
    width: 45%;
  }

  .dish-count-buttons #dish-count-decr-button {
    background-color: rgb(182, 19, 19);
    width: 45%;
  }
</style>
