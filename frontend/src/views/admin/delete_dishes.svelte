<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { navigate } from "svelte-routing";
  import { initBackButton } from "@telegram-apps/sdk";

  import { GetDishes } from "../../client/dish";
  import { LoadCart } from "../../client/cart";
  import { GetDishesCategories } from "../../client/dishes_categories";
  import DeleteDish from "./delete_dish.svelte";

  let dishes = [];
  let categories = [];
  let mainButton = LoadCart();
  let removeListFn;
  let removeBackButtonListFn;
  const backButtonRes = initBackButton();
  var backButton = backButtonRes[0];

  function dishDeleted(id) {
    currentOffset -= 1;
    dishes = dishes.filter((v) => v.id != id);
    console.log(dishes);
  }

  onMount(async () => {
    dishes = await GetDishes(null, pageLimit, currentOffset, null);
    currentOffset += dishes.length;
    categories = await GetDishesCategories();
    removeListFn = mainButton.on("click", () => {
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
  <div class="dishes-categories">
    {#each categories as category}
      <button
        class={category.id == selectedCategory
          ? "selected-category-button"
          : "category-button"}
        on:click={() => selectCategory(category.id)}>{category.name}</button
      >
    {/each}
  </div>
  <div class="spacer"></div>
  <div class="dishes-container">
    {#each dishes as dish}
      <DeleteDish bind:dish deletedFunc={dishDeleted} />
    {/each}
  </div>
</div>

<style>
  .dishes-container {
    display: flex;
    flex-wrap: wrap;
    flex-direction: row;
    justify-content: flex-start;
  }

  .dishes-categories {
    display: inline-flex;
    justify-content: center;
    overflow-x: scroll;
    scroll-snap-type: proximity;
    scroll-behavior: smooth;
    width: 95vw;
  }
  .category-button {
    font-size: medium;
    margin: 5px;
    white-space: nowrap;
  }
  .selected-category-button {
    background-color: var(--tg-theme-hint-color) !important;
    margin: 5px;
    white-space: nowrap;
  }
  .spacer {
    height: 5vh;
  }
</style>
