<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { navigate } from "svelte-routing";
  import { initBackButton } from "@telegram-apps/sdk";

  import Dish from "./dish.svelte";
  import { GetDishes } from "../../client/dish";
  import { LoadCart } from "../../client/cart";
  import { GetDishesCategories } from "../../client/dishes_categories";

  let dishes = [];
  let categories = [];
  let mainButton = LoadCart();
  let removeListFn;
  let removeBackButtonListFn;
  const backButtonRes = initBackButton();
  var backButton = backButtonRes[0];
  onMount(async () => {
    dishes = await GetDishes(null, pageLimit, currentOffset, null);
    currentOffset += dishes.length;
    categories = await GetDishesCategories();
    removeListFn = mainButton.on("click", () => {
      navigate("/cart");
      mainButton.hide();
    });
    removeBackButtonListFn = backButton.on("click", () => {
      navigate("/", { replace: true });
    });
  });

  onDestroy(() => {
    removeListFn();
    removeBackButtonListFn();
  });

  const pageLimit = 30;
  let currentOffset = 0;
  let selectedCategory = -1;
  async function selectCategory(categoryId) {
    if (categoryId == selectedCategory) {
      currentOffset = 0;
      dishes = await GetDishes(null, pageLimit, currentOffset, null);
      currentOffset += dishes.length;
      selectedCategory = -1;
    } else {
      currentOffset = 0;
      dishes = await GetDishes(null, pageLimit, currentOffset, categoryId);
      currentOffset += dishes.length;
      selectedCategory = categoryId;
    }
  }
</script>

<div class="dish">
  <div class="dishes_categories">
    {#each categories as category}
      <button
        class={category.id == selectedCategory
          ? "selected_category_button"
          : "category_button"}
        on:click={() => selectCategory(category.id)}>{category.name}</button
      >
    {/each}
  </div>
  <div class="spacer"></div>
  <div class="grid">
    {#each dishes as dish}
      <Dish bind:dish />
    {/each}
  </div>
</div>

<style>
  .grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 1rem;
  }
  .dishes_categories {
    display: inline-flex;
    overflow-x: scroll;
    max-width: 90vw;
  }
  .category_button {
    font-size: medium;
    margin: 5px;
    white-space: nowrap;
  }
  .selected_category_button {
    background-color: var(--tg-theme-hint-color)!important;
    margin: 5px;
    white-space: nowrap;
  }
  .spacer {
    height: 5vh;
  }
</style>
