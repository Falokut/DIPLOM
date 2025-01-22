<script lang="ts">
  import { onMount } from "svelte";
  import { GetDishes } from "../../client/dish";

  export let pageLimit = 10;
  export let DishComponent;
  let selectedCategory = -1;

  export function selectedCategoryUpdated(newSelectedCategory) {
    if (newSelectedCategory == selectedCategory) return;

    selectedCategory = newSelectedCategory;
    allDishesFetched = false;
    dishes = [];
    fetchDishes();
  }

  onMount(() => {
    fetchDishes();
    window.addEventListener("scroll", handleScroll);
  });

  let dishes = [];
  let loading = false;
  let allDishesFetched = false;

  async function fetchDishes() {
    if (loading || allDishesFetched) return;

    loading = true;
    let fetchedDishes = await GetDishes(
      null,
      pageLimit,
      dishes.length,
      selectedCategory == -1 ? null : selectedCategory
    );

    dishes = [...dishes, ...fetchedDishes];
    dishes = dishes;
    loading = false;
    if (fetchedDishes.length < pageLimit) {
      allDishesFetched = true;
    }
  }

  let container;
  function handleScroll() {
    if (!container) return;
    const { scrollTop, scrollHeight, clientHeight } = container;

    if (scrollTop + clientHeight >= scrollHeight - 10) {
      fetchDishes();
    }
  }
</script>

<div bind:this={container} class="dishes-container">
  {#each dishes as dish}
    <svelte:component
      this={DishComponent}
      bind:dish
      onRemove={(dishId) => {
        dishes = dishes.filter((v) => v.id != dishId);
      }}
    />
  {/each}
</div>

<style>
  .dishes-container {
    overflow-y: auto;
    display: flex;
    flex-wrap: wrap;
    flex-direction: row;
    justify-content: space-around;
  }
</style>
