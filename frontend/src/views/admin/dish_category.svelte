<script lang="ts">
  import {
    RenameDishesCategory,
    DeleteDishesCategory,
  } from "../../client/dishes_categories";
  import TextInput from "../components/text_input.svelte";

  export let category = {
    name: "",
    id: 0,
  };
  export let remove = function () {};

  async function updateDishName() {
    RenameDishesCategory(category.name, category.id);
  }

  async function deleteCategory() {
    await DeleteDishesCategory(category.id);
    remove();
  }
</script>

<div class="category-container">
  <TextInput bind:value={category.name} />
  <div class="buttons-container">
    <button class="apply-button" on:click={() => updateDishName()}>✓</button>
    <button class="remove-button" on:click={() => deleteCategory()}>✕</button>
  </div>
</div>

<style>
  .category-container {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-around;
    padding: 10px;
  }

  .buttons-container {
    text-align: center;
    vertical-align: middle;
    display: flex;
    flex-direction: row;
    font-weight: 100;
    align-items: center;
    justify-content: space-around;
  }

  .remove-button {
    background-color: rgb(182, 19, 19);
  }

  .remove-button:hover {
    border: 1px solid var(--tg-theme-bg-color);
  }

  .apply-button {
    background-color: rgb(17, 151, 17);
  }

  .apply-button:hover {
    border: 1px solid var(--tg-theme-bg-color);
  }
</style>
