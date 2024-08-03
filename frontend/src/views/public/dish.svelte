<script lang="ts">
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
    dish = dishFromObject(dishInst);
    count = GetDishCount(dish.id);
    if (count > 0) {
      showIncrDecrButtons = true;
    }
  });
</script>

<div class="dish">
  <div>
    {#if count > 0}
      <div class="dish-count">{count}</div>
    {/if}
    <img class="dish-preview" src={dish.url} alt={dish.name} />
    <div class="dish-name">{dish.name}</div>
    <div class="dish-price">{dish.price}</div>
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
</div>

<style>
  :root {
    --dish-size: 120px;
    --dish-button-width: 120px;
    --dish-counter-size: 25px;
  }

  .dish-count {
    position: relative;
    width: var(--dish-counter-size);
    height: var(--dish-counter-size);
    border-radius: 50%;

    margin-left: calc(var(--dish-size) * 0.7);
    margin-bottom: calc(var(--dish-counter-size) * -1);

    background-color: orange;

    font-size: medium;
    font-weight: bold;
    text-align: center;
  }

  .dish-name {
    height: 10%;
    font-size: small;
  }

  .dish-price {
    height: 10%;
    font-size: small;
  }

  .dish-preview {
    width: auto;
    height: calc(var(--dish-size) * 0.7);
    object-fit: scale-down;
  }

  .dish {
    text-align: center;
    width: var(--dish-size);
    height: var(--dish-size);
  }

  .dish-count-buttons {
    position: relative;
    grid-template-columns: calc(var(--dish-button-size) / 2 * 0.9);
    gap: calc(var(--dish-button-size) / 2 * 0.1);
    grid-auto-flow: column;

    color: white;
    text-align: center;
    font-size: medium;
    font-weight: bold;
    width: 100%;
  }

  .dish-count-button {
    display: inline-block;
    justify-content: center;
    align-items: center;
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
