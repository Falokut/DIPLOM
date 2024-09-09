<script lang="ts">
  import { onMount } from "svelte";
  import { SetDishCount, GetDishCount } from "../../client/cart";
  import { DishFromObject } from "../../client/dish";
  import PreviewImage from "../components/preview_image.svelte";
  export let dishInst = {
    id: 0,
    url: "",
    name: "",
    description: "",
    categories: "",
    price: 0,
  };

  let dish = {
    id: 0,
    url: "",
    name: "",
    price: "",
  };
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
    dish = DishFromObject(dishInst);
    count = GetDishCount(dish.id);
    if (count > 0) {
      showIncrDecrButtons = true;
    }
  });
  var previewSize = 90;
</script>

<div class="dish">
  {#if count > 0}
    <div class="dish-count">{count}</div>
  {/if}
  <PreviewImage bind:url={dish.url} bind:size={previewSize} bind:alt={dish.name}
  ></PreviewImage>
  <div class="dish-caption">{dish.name}</div>
  <div class="dish-caption">{dish.price}</div>
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
  :root {
    --dish-width: 110px;
    --dish-height: 200px;
    --dish-preview-size: 90px;
    --dish-counter-size: 25px;
  }

  .dish-count {
    position: relative;
    width: var(--dish-counter-size);
    height: var(--dish-counter-size);
    border-radius: 40%;

    margin-left: calc(var(--dish-width) - var(--dish-counter-size));
    margin-bottom: calc(var(--dish-counter-size) * -1);

    background-color: orange;

    font-size: medium;
    font-weight: bold;
    text-align: center;
  }

  .dish-caption {
    font-size: small;
    margin-top: auto;
  }

  .dish {
    text-align: center;
    display: inline-block;

    width: var(--dish-width);
    height: var(--dish-height);
  }

  .dish-count-buttons {
    position: relative;
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    grid-auto-flow: column;
    gap: 0.5rem;
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
  }

  .dish-count-buttons #dish-count-decr-button {
    background-color: rgb(182, 19, 19);
  }
</style>
