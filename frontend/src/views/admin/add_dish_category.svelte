<script lang="ts">
  import {
    AddDishCategory,
    DishCategory,
  } from "../../client/dishes_categories";

  export let categoryName = "";
  export let OnAdd = (category: DishCategory) => {};

  async function addDishCategory() {
    if (categoryName == "") {
      return;
    }

    let dishCategory = await AddDishCategory(categoryName);
    if (dishCategory.id != undefined) {
      dishCategory.name = categoryName;
      OnAdd(dishCategory);
      categoryName = "";
      return;
    }
  }
</script>

<section class="category">
  <input
    class="input-area category-input"
    type="text"
    bind:value={categoryName}
  />
  <button class="add-button" on:click={addDishCategory}>Добавить</button>
</section>

<style>
  .category {
    border-radius: 5px;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-around;
    background-color: var(--secondary-bg-color);
    align-items: center;
    margin: auto;
    padding: 5px;
    width: 100vw;
  }
  .category-input {
    border-radius: 3px;
  }

  .add-button {
    display: block;
    border: 1px solid var(--primary-bg-color);
    margin-right: 5px;
  }
</style>
