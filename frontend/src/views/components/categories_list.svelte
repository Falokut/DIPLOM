<script lang="ts">
  import { onMount } from "svelte";
  import { GetDishesCategories } from "../../client/dishes_categories";

  export const CategoryChangedEventType = "category_changed";
  let categories = [];

  onMount(async () => {
    categories = await GetDishesCategories();
  });
  let selectedCategory = -1;
  function selectCategory(categoryId) {
    if (categoryId == selectedCategory) {
      selectedCategory = -1;
    } else {
      selectedCategory = categoryId;
    }
    let categoryChangedEvent = new CustomEvent(CategoryChangedEventType, {
      detail: { categoryId: selectedCategory },
    });
    dispatchEvent(categoryChangedEvent);
  }
</script>

<section class="dishes-categories">
  {#each categories as category}
    <button
      class={category.id == selectedCategory
        ? "selected"
        : ""}
      on:click={() => selectCategory(category.id)}>{category.name}</button
    >
  {/each}
</section>

<style>
  .dishes-categories {
    display: flex;
    gap: 1rem;
    overflow-x: scroll;
  }
</style>
